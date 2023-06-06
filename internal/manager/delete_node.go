package manager

import (
	"context"
	"fmt"
	"log"

	"github.com/darchlabs/nodes/internal/command"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Manager) DeleteNode(ni *NodeInstance) error {
	ctx := context.Background()

	// Delete Services
	for _, svc := range ni.Artifacts.Services {
		err := m.clusterClient.CoreV1().Services("default").Delete(ctx, svc, metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrap(err, "manager: Manager.DeleteArtifacts m.clusterClient.CoreV1().Services().Delete error")
		}
		log.Printf("[MANAGER] Delete Service artifact <%s>", svc)
	}

	// Delete Pods
	for _, pod := range ni.Artifacts.Pods {
		err := m.clusterClient.CoreV1().Pods("default").Delete(ctx, pod, metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrap(err, "manager: Manager.DeleteArtifacts m.clusterClient.CoreV1().Pods().Delete")
		}
		log.Printf("[MANAGER] Delete Pod artifact <%s>", pod)
	}

	// Delete Deployments
	for _, deployment := range ni.Artifacts.Deployments {
		err := m.clusterClient.AppsV1().Deployments("default").Delete(ctx, deployment, metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrap(err, "manager: Manager.DeleteArtifacts m.clusterClient.AppsV1().Deployments().Delete")
		}
		log.Printf("[MANAGER] Delete Deployment artifact <%s>", deployment)
	}

	// Delete folders
	pathDir := fmt.Sprintf("%s%s/%s/%s", m.basePathDB, ni.Config.Network, ni.Config.Environment, ni.Name)
	rmDir := command.New("rm", "-rf", pathDir)

	err := rmDir.Start()
	if err != nil {
		return errors.Wrap(err, "manager: Manager.DeleteNode rmDir.Start error")
	}
	log.Printf("[MANAGER] Delete folder <%s>", pathDir)

	return nil
}
