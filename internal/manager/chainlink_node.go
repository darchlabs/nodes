package manager

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

//"context"
//"fmt"
//"k8s.io/client-go/kubernetes"
//"k8s.io/client-go/tools/clientcmd"
//"k8s.io/client-go/util/homedir"
//"path/filepath"

//corev1 "k8s.io/api/core/v1"
//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//"k8s.io/apimachinery/pkg/util/intstr"
//appsv1 "k8s.io/api/apps/v1"

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
	arts := &Artifacts{}
	nodeID := m.idGenerator()
	nodeName := m.nameGenerator.Generate()
	// check chainlink basics. Create folders and files
	networkDir := fmt.Sprintf("%s/%s/%s/%s", m.basePathDB, network, env["ENVIRONMENT"], nodeName)
	err := m.chainlinkBasics(networkDir, env["PASSWORD"], env["NODE_EMAIL"], env["NODE_EMAIL_PWD"])
	if err != nil {
		return nil, errors.Wrap(err, "manager: m.ChainlinkNode m.createChainlinkBasics error")
	}

	envVars := make([]corev1.EnvVar, 0)
	// user provided ev vars
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

	// handled env vars
	// TODO: use dynamic password
	dbUser := "postgres"
	dbPass := "ThisPasswordIsSecure"
	dbURL := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", dbUser, dbPass, nodeName)
	envVars = append(envVars, []corev1.EnvVar{
		{
			Name:  "DATABASE_URL",
			Value: dbURL,
		},
		{
			Name:  "ROOT",
			Value: networkDir,
		},
	}...)

	fmt.Printf("THIS IS THE ENV VAR \n%+v\n", envVars)

	return &NodeInstance{
		ID:        nodeID,
		Name:      nodeName,
		Artifacts: arts,
		Config: &NodeConfig{
			Host:              "0.0.0.0",
			Network:           network,
			Environment:       env["ENVIRONMENT"],
			BaseChainDataPath: networkDir,
			CreatedAt:         time.Now(),
		},
	}, nil
}

func (m *Manager) chainlinkBasics(networkDir, nodePwd, email, emailPwd string) error {
	// 1. check folder exist
	//		create folder for chainlink if not exist
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

	creds := fmt.Sprintf("%s\n%s", email, emailPwd)
	err = ioutil.WriteFile(fmt.Sprintf("%s/creds.txt", networkDir), []byte(creds), 0644)
	if err != nil {
		return errors.Wrap(err, "manager: Manager.createChainlinkBasics ioutil.WriteFile creds.txt error")
	}

	return nil
}
