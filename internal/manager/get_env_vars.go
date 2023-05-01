package manager

import (
	corev1 "k8s.io/api/core/v1"
)

func getEnvVars(envList []string, env map[string]string, envVars []corev1.EnvVar) ([]corev1.EnvVar, error) {
	for _, key := range envList {
		value, ok := env[key]
		if !ok {
			return nil, ErrKeyNotFound
		}

		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	return envVars, nil
}
