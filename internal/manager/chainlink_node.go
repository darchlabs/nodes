package manager

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	chainlinkDefaultEnvVars = map[string]string{
		"LOG_LEVEL":                  "debug",
		"MIN_OUTGOING_CONFIRMATIONS": "2",
		"CHAINLINK_TLS_PORT":         "0",
		"SECURE_COOKIES":             "false",
		"GAS_UPDATER_ENABLED":        "true",
		"FEATURE_FLUX_MONITOR":       "true",
		"ALLOW_ORIGINS":              "*",
		"DATABASE_TIMEOUT":           "0",
	}

	chainlinkNetworkEnvVars = map[string]map[string]string{
		"sepolia": {
			"ETH_CHAIN_ID":          "11155111",
			"LINK_CONTRACT_ADDRESS": "0x779877A7B0D9E8603169DdbD7836e478b4624789",
		},
	}

	chainlinkContainerEnvList = []string{
		//"DATABASE_URL",          // HANDLED
		//"ROOT",                  // HANDLED

		"ENVIRONMENT",    // FROM REQUEST
		"ETH_URL",        // FROM REQUEST
		"PASSWORD",       // FROM REQUEST
		"NODE_EMAIL",     // FROM REQUEST
		"NODE_EMAIL_PWD", // FROM REQUEST
	}
)

var chainlinkNodeDir = "chainlink"

func (m *Manager) ChainlinkNode(network string, env map[string]string) (*NodeInstance, error) {
	nodeID := m.idGenerator()
	nodeName := m.nameGenerator.Generate()
	networkEnv, ok := env["ENVIRONMENT"]
	if !ok {
		return nil, ErrKeyNotFound
	}
	managerName := fmt.Sprintf("%s-%s-%s", network, networkEnv, nodeName)

	// ## Check chainlink basics. Create folders and files
	networkDir := fmt.Sprintf("%s%s/%s/%s", m.basePathDB, network, networkEnv, nodeName)
	err := m.chainlinkBasics(networkDir, env["PASSWORD"], env["NODE_EMAIL"], env["NODE_EMAIL_PWD"])
	if err != nil {
		return nil, errors.Wrap(err, "manager: m.ChainlinkNode m.createChainlinkBasics error")
	}

	envVars := make([]corev1.EnvVar, 0)
	// user provided env vars
	envVars, err = getEnvVars(chainlinkContainerEnvList, env, envVars)
	if err != nil {
		return nil, errors.Wrap(err, "manager: getEnvVars error")
	}

	// default env vars
	defaultEnvKeys := make([]string, 0, len(chainlinkDefaultEnvVars))
	for k := range chainlinkDefaultEnvVars {
		defaultEnvKeys = append(defaultEnvKeys, k)
	}
	envVars, _ = getEnvVars(defaultEnvKeys, chainlinkDefaultEnvVars, envVars)

	// ## Handle env vars
	// TODO: use dynamic password
	dbPass := "ThisPasswordIsSecure"
	dbURL := fmt.Sprintf("postgres://postgres:%s@postgres:5432/%s?sslmode=disable", dbPass, nodeName)
	envVars = append(envVars, []corev1.EnvVar{
		{
			Name:  "DATABASE_URL",
			Value: dbURL,
		},
		{
			Name:  "ROOT",
			Value: fmt.Sprintf("/%s", nodeName),
		},
	}...)

	arts := &Artifacts{
		Deployments: []string{},
		Pods:        []string{},
		Services:    []string{},
	}

	// ## Create artifacts needed
	ctx := context.Background()

	// ### POSTGRES RELATED
	psqlNameRef := fmt.Sprintf("postgres-%s", managerName)

	// 1. Create Chainlink-postgres-db service
	arts.Services = append(arts.Services, psqlNameRef)
	psqlSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: psqlNameRef,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":  psqlNameRef,
				"role": psqlNameRef,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       psqlNameRef,
					Protocol:   corev1.ProtocolTCP,
					Port:       5432,
					TargetPort: intstr.FromInt(5432),
				},
			},
			ClusterIP: "None",
		},
	}
	_, err = m.clusterClient.CoreV1().Services("default").Create(ctx, psqlSvc, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.ChainlinkNode m.clusterClient.CoreV1().Service.Create psql error")
	}

	log.Printf("[MANAGER] Chainlink-postgres service <%s> created succesfully", psqlNameRef)

	// 2. Create Chainlink-postgres-db deployment
	dbReplicas := int32(1)
	postgresDeploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: psqlNameRef,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &dbReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": psqlNameRef,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": psqlNameRef,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  psqlNameRef,
							Image: "postgres:13",
							Env: []corev1.EnvVar{
								{
									Name:  "POSTGRES_DB",
									Value: managerName,
								},
								{
									Name:  "POSTGRES_USER",
									Value: "postgres",
								},
								{
									Name:  "POSTGRES_PASSWORD",
									Value: "ThisPasswordIsSecure",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "postgres",
									ContainerPort: 5432,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "postgres-db",
									MountPath: "/var/lib/postgresql/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: psqlNameRef,
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: fmt.Sprintf("/mnt/data/nodes-volume/databases/%s", nodeName),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = m.clusterClient.AppsV1().Deployments("default").Create(
		ctx,
		postgresDeploy,
		metav1.CreateOptions{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.ChainlinkNode m.clusterClient.AppsV1.Deployment.Create psql-deployment error")
	}

	log.Printf("[MANAGER] Chainlink-postgres deployment <%s> created succesfully", psqlNameRef)

	// ### CHAINLINK RELATED
	chainlinkNameRef := fmt.Sprintf("chainlink-%s", managerName)

	// 3. Create chainlink-node-svc
	arts.Services = append(arts.Services, chainlinkNameRef)
	chainlinkSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: chainlinkNameRef,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":  chainlinkNameRef,
				"role": chainlinkNameRef,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       psqlNameRef,
					Protocol:   corev1.ProtocolTCP,
					Port:       5432,
					TargetPort: intstr.FromInt(5432),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	_, err = m.clusterClient.CoreV1().Services("default").Create(ctx, chainlinkSvc, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.ChainlinkNode m.clusterClient.CoreV1().Service.Create psql error")
	}

	log.Printf("[MANAGER] Chainlink-postgres node service <%s> created succesfully", chainlinkNameRef)

	// 4. Create Chainlink-node-deployment
	chainlinkDeploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: psqlNameRef,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &dbReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": chainlinkNameRef,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": chainlinkNameRef,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  chainlinkNameRef,
							Image: "smartcontract/chainlink:1.13.1-root",
							Env:   envVars,
							Ports: []corev1.ContainerPort{
								{
									Name:          "postgres",
									ContainerPort: 5432,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "nodes-pvc",
									MountPath: fmt.Sprintf("/%s", nodeName),
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: chainlinkNameRef,
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: fmt.Sprintf("/mnt/data/nodes-volume/%s/%s/%s", network, networkEnv, nodeName),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = m.clusterClient.AppsV1().Deployments("default").Create(
		ctx,
		chainlinkDeploy,
		metav1.CreateOptions{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.ChainlinkNode m.clusterClient.AppsV1.Deployment.Create chainlink-deployment error")
	}

	log.Printf("[MANAGER] Chainlink-postgres deployment <%s> created succesfully", psqlNameRef)

	return &NodeInstance{
		ID:        nodeID,
		Name:      nodeName,
		Artifacts: arts,
		Config: &NodeConfig{
			Host:              "0.0.0.0",
			Network:           network,
			Environment:       networkEnv,
			BaseChainDataPath: networkDir,
			CreatedAt:         time.Now(),
		},
	}, nil
}

func (m *Manager) chainlinkBasics(networkDir, nodePwd, email, emailPwd string) error {
	// 1. check folder exist
	//		create folder for chainlink if not exist
	log.Printf("[MANAGER] directory created %s\n", networkDir)
	err := m.createDirIfNotExist(networkDir)
	if err != nil {
		return errors.Wrap(err, "manager: Manager.createChainlinkBasics exist error")
	}

	time.Sleep(1 * time.Second)
	// 2. create files

	err = ioutil.WriteFile(fmt.Sprintf("%s/password.txt", networkDir), []byte(nodePwd), 0644)
	if err != nil {
		return errors.Wrap(err, "manager: Manager.createChainlinkBasics ioutil.WriteFile password.txt error")
	}
	log.Printf("[MANAGER] file created %s/password.txt\n", networkDir)

	creds := fmt.Sprintf("%s\n%s", email, emailPwd)
	err = ioutil.WriteFile(fmt.Sprintf("%s/creds.txt", networkDir), []byte(creds), 0644)
	if err != nil {
		return errors.Wrap(err, "manager: Manager.createChainlinkBasics ioutil.WriteFile creds.txt error")
	}
	log.Printf("[MANAGER] file created %s/creds.txt\n", networkDir)

	return nil
}
