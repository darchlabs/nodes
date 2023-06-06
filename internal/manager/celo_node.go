package manager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (m *Manager) CeloNode(network string, env map[string]string) (*NodeInstance, error) {
	nodeID := m.idGenerator()
	nodeName := m.nameGenerator.Generate()
	networkEnv, ok := env["ENVIRONMENT"]
	if !ok {
		return nil, ErrKeyNotFound
	}

	networkDir := fmt.Sprintf("%s%s/%s/%s", m.basePathDB, network, networkEnv, nodeName)
	err := m.createDirIfNotExist(networkDir)
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.CeloNode m.createDirIfNotExist error")
	}

	// get env vars
	//envVars := getFromMap(env)

	arts := &Artifacts{
		Deployments: []string{},
		Pods:        []string{},
		Services:    []string{},
	}

	// ## Create artifacts needed
	ctx := context.Background()

	// 1. Create Celo service
	celoNameRef := fmt.Sprintf("celo-%s", nodeName)

	log.Printf("[MANAGER] start Celo service <%s> creation", celoNameRef)

	celoSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: celoNameRef,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":  celoNameRef,
				"role": celoNameRef,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp",
					Protocol:   corev1.ProtocolTCP,
					Port:       8545,
					TargetPort: intstr.FromInt(8545),
				},
				{
					Name:       "ws",
					Protocol:   corev1.ProtocolTCP,
					Port:       8546,
					TargetPort: intstr.FromInt(8546),
				},
				{
					Name:       "p2p",
					Protocol:   corev1.ProtocolTCP,
					Port:       30303,
					TargetPort: intstr.FromInt(30303),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	_, err = m.clusterClient.CoreV1().Services("default").Create(ctx, celoSvc, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.CeloNode m.clusterClient.CoreV1().Services().Create celo error")
	}
	arts.Services = append(arts.Services, celoNameRef)

	log.Printf("[MANAGER]  Celo service <%s> created [DONE ✔︎]", celoNameRef)

	// 2. Create Celo deployment
	log.Printf("[MANAGER] start Celo service <%s> creation", celoNameRef)

	nodePwd, ok := env["PASSWORD"]
	if !ok {
		return nil, ErrKeyNotFound
	}

	podReplicas := int32(1)
	celoDeploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: celoNameRef,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &podReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  celoNameRef,
					"role": celoNameRef,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  celoNameRef,
						"role": celoNameRef,
					},
				},

				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "celo-node",
							Image: "us.gcr.io/celo-org/geth:alfajores",
							Args: []string{
								"--verbosity", "3",
								"--syncmode", "full",
								"--http", "--http.addr", "0.0.0.0",
								"--http.api", "eth,net,web3,debug,admin,personal",
								"--light.serve", "90",
								"--light.maxpeers", "1000",
								"--maxpeers", "1100",
								"--alfajores",
								"--datadir", "/root/.celo",
							},
							Env: []corev1.EnvVar{
								{
									Name:  "CELO_ACCOUNT_PASSWORD",
									Value: nodePwd,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8545,
								},
								{
									ContainerPort: 8546,
								},
								{
									ContainerPort: 30303,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "nodes-pvc",
									MountPath: "/root/.celo",
								},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:  "celo-account-creator",
							Image: "us.gcr.io/celo-org/geth:alfajores",
							Command: []string{
								"/bin/sh", "-c",
							},
							Args: []string{
								"echo $CELO_ACCOUNT_PASSWORD > /root/.celo/password.txt && /usr/local/bin/geth account new --password /root/.celo/password.txt",
							},
							Env: []corev1.EnvVar{
								{
									Name:  "CELO_ACCOUNT_PASSWORD",
									Value: nodePwd,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "nodes-pvc",
									MountPath: "/root/.celo",
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
		celoDeploy,
		metav1.CreateOptions{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "manager: Manager.CeloNode m.clusterClient.AppsV1().Deployment.Create celo-deployment error")
	}
	arts.Deployments = append(arts.Deployments, celoNameRef)

	log.Printf("[MANAGER] start Celo service <%s> creation", celoNameRef)

	return &NodeInstance{
		ID:        nodeID,
		Name:      nodeName,
		Artifacts: arts,
		Config: &NodeConfig{
			Host:        "0.0.0.0",
			Port:        8545,
			Network:     network,
			Environment: networkEnv,
			CreatedAt:   time.Now(),
		},
	}, nil
}
