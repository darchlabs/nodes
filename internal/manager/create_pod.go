package manager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type CreatePodOptions struct {
	Image   string
	EnvVars map[string]string
}

func (m *Manager) CreatePod(opts *CreatePodOptions) (*NodeInstance, error) {
	envVars := make([]corev1.EnvVar, 0)

	for k, v := range opts.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	containerName := fmt.Sprintf("ethereum-%s", m.nameGenerator.Generate())

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
					Image:           "darchlabs/node-ethereum-dev:0.0.2",
					ImagePullPolicy: corev1.PullIfNotPresent,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8545,
							HostPort:      8545,
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
		return nil, errors.Wrap(err, "manager: Manager.CreatePod m.clusterClient.CoreV1().Pods().Create")
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
					Port:       8545,
					TargetPort: intstr.FromInt(8545),
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
		fmt.Println("--ERROR-- ", err.Error())
		return nil, errors.Wrap(err, "manager: Manager.CreatePod m.clusterClient.CoreV1().Services().Create")
	}

	return &NodeInstance{
		ID:     m.idGenerator(),
		Name:   containerName,
		Config: &NodeConfig{},
	}, nil
}

func intPtr32(i int32) *int32 {
	return &i
}
