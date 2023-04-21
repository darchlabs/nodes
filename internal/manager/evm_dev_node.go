package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (m *Manager) EvmDevNode(network string, env map[string]string) (*NodeInstance, error) {
	envVars := make([]corev1.EnvVar, 0)

	if _, ok := env["FROM_BLOCK_NUMBER"]; ok {
		env["NETWORK_URL"] = fmt.Sprintf("%s@%s", env["NETWORK_URL"], env["FROM_BLOCK_NUMBER"])
		delete(env, "FROM_BLOCK_NUMBER")
	}

	for k, v := range env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	// main receipe
	containerName := fmt.Sprintf("%s-%s", network, m.nameGenerator.Generate())
	arts := &Artifacts{
		Pods:     []string{containerName},
		Services: []string{containerName},
	}

	// TODO: define this as deployment so k8s can be in charge of restart nodes
	// Define the new pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: containerName,
			Labels: map[string]string{
				"role": containerName,
				"app":  containerName,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: containerName,
					//Image:           opts.Image,
					Image:           m.MainConfig.Images["evm"],
					ImagePullPolicy: corev1.PullIfNotPresent,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8544, // TODO: make it dynamic
							HostPort:      8545, // TODO: make it dynamic
							Name:          "http",
							Protocol:      corev1.ProtocolTCP,
						},
					},
					Env: envVars,
				},
			},
		},
	}

	_, err := m.clusterClient.CoreV1().Pods("default").Create(
		context.Background(),
		pod,
		metav1.CreateOptions{},
	)
	if err != nil {
		fmt.Println("--ERROR-- ", err.Error())
		return nil, errors.Wrap(err, "manager: evmNode m.clusterClient.CoreV1().Pods().Create")
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: containerName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"role": containerName,
				"app":  containerName,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       8545,                 // TODO: make it dynamic
					TargetPort: intstr.FromInt(8545), // TODO: make it dynamic
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// Create the new pod
	_, err = m.clusterClient.CoreV1().Services("default").Create(
		context.Background(),
		service,
		metav1.CreateOptions{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "manager: evmNode m.clusterClient.CoreV1().Services().Create")
	}

	return &NodeInstance{
		ID:        m.idGenerator(),
		Name:      containerName,
		Artifacts: arts,
		Config: &NodeConfig{
			Host:      containerName,
			Network:   network,
			Port:      8545, // TODO: make it dynamic
			CreatedAt: time.Now(),
		},
	}, nil
}
