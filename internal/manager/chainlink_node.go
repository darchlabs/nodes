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

	// ## Check chainlink basics. Create folders and files
	networkDir := fmt.Sprintf("%s%s/%s/%s", m.basePathDB, network, networkEnv, nodeName)
	err := m.chainlinkBasics(networkDir, env["PASSWORD"], env["NODE_EMAIL"], env["NODE_EMAIL_PWD"])
	if err != nil {
		return nil, errors.Wrap(err, "manager: m.ChainlinkNode m.createChainlinkBasics error")
	}

	envVars := make([]corev1.EnvVar, 0)
	// ## Handle env vars
	// TODO: use dynamic password
	dbPass := "ThisPasswordIsSecure"
	psqlNameRef := fmt.Sprintf("postgres-%s", nodeName)
	dbURL := fmt.Sprintf("postgres://postgres:%s@%s:5432/postgres?sslmode=disable", dbPass, psqlNameRef)
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

	// user provided env vars
	// TODO: validate provided env vars
	providedEnvVars := getFromMap(env)

	// default env vars
	defaultEnvVars := getFromMap(chainlinkDefaultEnvVars)

	// network related env vars
	networkEnvVars := getFromMap(chainlinkNetworkEnvVars[networkEnv])

	// merge env vars
	containerEnvVars := mergeEnvVars(
		envVars,
		providedEnvVars,
		defaultEnvVars,
		networkEnvVars,
	)

	for _, v := range containerEnvVars {
		log.Printf("[MANAGER] ENV VAR - %+v\n", v)
	}

	arts := &Artifacts{
		Deployments: []string{},
		Pods:        []string{},
		Services:    []string{},
	}

	// ## Create artifacts needed
	ctx := context.Background()

	// ### POSTGRES RELATED

	// 1. Create Chainlink-postgres-db service
	log.Printf("[MANAGER] start Chainlink-postgres service <%s> creation", psqlNameRef)

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
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       5432,
					TargetPort: intstr.FromInt(5432),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	_, err = m.clusterClient.CoreV1().Services("default").Create(ctx, psqlSvc, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.ChainlinkNode m.clusterClient.CoreV1().Service.Create psql error")
	}

	log.Printf("[MANAGER] Chainlink-postgres service <%s> created [DONE ✔︎]", psqlNameRef)

	// 2. Create Chainlink-postgres-db deployment
	log.Printf("[MANAGER] start Chainlink-postgres deployment <%s> creation", psqlNameRef)
	dbReplicas := int32(1)
	postgresDeploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: psqlNameRef,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &dbReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  psqlNameRef,
					"role": psqlNameRef,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  psqlNameRef,
						"role": psqlNameRef,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  psqlNameRef,
							Image: "postgres:13",
							Env: []corev1.EnvVar{
								{
									Name:  "PGSSLMODE",
									Value: "disable",
								},
								{
									Name:  "POSTGRES_DB",
									Value: "postgres",
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
									Name:          "http",
									ContainerPort: 5432,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "nodes-pvc",
									MountPath: "/var/lib/postgresql/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "nodes-pvc",
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
		return nil, errors.Wrap(err, "manager: Manager.ChainlinkNode m.clusterClient.AppsV1().Deployment.Create psql-deployment error")
	}

	log.Printf("[MANAGER] Chainlink-postgres deployment <%s> created [DONE ✔︎]", psqlNameRef)

	// ### CHAINLINK RELATED
	chainlinkNameRef := fmt.Sprintf("chainlink-%s", nodeName)
	log.Printf("[MANAGER] start Chainlink-node service <%s> creation", chainlinkNameRef)

	// 3. Create chainlink-node-svc
	arts.Services = append(arts.Services, chainlinkNameRef)
	chainlinkSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: chainlinkNameRef,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": chainlinkNameRef,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       chainlinkNameRef,
					Protocol:   corev1.ProtocolTCP,
					Port:       6688,
					TargetPort: intstr.FromInt(6688),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	_, err = m.clusterClient.CoreV1().Services("default").Create(ctx, chainlinkSvc, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.ChainlinkNode m.clusterClient.CoreV1().Service.Create service error")
	}

	log.Printf("[MANAGER] Chainlink-node service <%s> created [DONE ✔︎]", chainlinkNameRef)

	// 4. Create Chainlink-node-deployment
	log.Printf("[MANAGER] start Chainlink-node deployment <%s> creation", chainlinkNameRef)

	chainlinkDeploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: chainlinkNameRef,
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
							//Command: []string{"/bin/sh", "-c"},
							Args: []string{
								"node",
								"start",
								"--password",
								fmt.Sprintf("/%s/password.txt", nodeName),
								"--api",
								fmt.Sprintf("/%s/creds.txt", nodeName),
							},
							Env: containerEnvVars,
							Ports: []corev1.ContainerPort{
								{
									Name:          "chainlink",
									ContainerPort: 6688,
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
							Name: "nodes-pvc",
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
		return nil, errors.Wrap(err, "manager: Manager.ChainlinkNode m.clusterClient.AppsV1().Deployment.Create chainlink-deployment error")
	}

	log.Printf("[MANAGER] Chainlink-node deployment <%s> created [DONE ✔︎]", chainlinkNameRef)

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

	creds := fmt.Sprintf("%s\n%s\n", email, emailPwd)
	err = ioutil.WriteFile(fmt.Sprintf("%s/creds.txt", networkDir), []byte(creds), 0644)
	if err != nil {
		return errors.Wrap(err, "manager: Manager.createChainlinkBasics ioutil.WriteFile creds.txt error")
	}
	log.Printf("[MANAGER] file created %s/creds.txt\n", networkDir)

	return nil
}
