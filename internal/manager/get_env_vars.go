package manager

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

func getFromMap(envs map[string]string) []corev1.EnvVar {
	arr := make([]corev1.EnvVar, 0)
	for key, value := range envs {
		arr = append(arr, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	return arr
}

func getEnvVars(envList []string, env map[string]string) ([]corev1.EnvVar, error) {
	envVars := make([]corev1.EnvVar, len(envList))

	for _, key := range envList {
		value, ok := env[key]
		if !ok {
			return nil, errors.Wrapf(ErrKeyNotFound, "manager: getEnvVars %s error", key)
		}

		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	return envVars, nil
}

func mergeEnvVars(arrs ...[]corev1.EnvVar) []corev1.EnvVar {
	merged := make([]corev1.EnvVar, 0)
	for _, arr := range arrs {
		merged = append(merged, arr...)
	}

	return merged
}
