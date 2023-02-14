package manager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type CreatePodOptions struct {
	Image   string
	EnvVars map[string]string
}

func (m *Manager) CreatePod(opts *CreatePodOptions) error {
	// Create a new Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "manager: CreatePod.CreatePod rest.InClusterConfig error")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "manager: CreatePod.CreatePod kubernetes.NewForConfig error")
	}
	fmt.Print("------- clientset", clientset.CoreV1())
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "manager: CreatePod.CreatePod clientset.CoreV1().Pods().List error")
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	name := m.nameGenerator.Generate()

	envVars := make([]corev1.EnvVar, 0)

	for k, v := range opts.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	// Define the new pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            fmt.Sprintf("ethereum-%s", name),
					Image:           "darchlabs/node-ethereum-dev",
					ImagePullPolicy: corev1.PullNever,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8545,
						},
					},
					Env: envVars,
				},
			},
		},
	}

	// Create the new pod
	_, err = clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("--ERROR-- ", err.Error())
		return err
	}

	return nil
}
