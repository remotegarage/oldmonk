package x

import (
	oldmonkv1 "github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1"
	corev1 "k8s.io/api/core/v1"
)

// GetPodNames returns the pod names of the array of pods passed in
func GetPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func GetLabels(m *oldmonkv1.QueueAutoScaler) *oldmonkv1.QueueAutoScaler {
	return m
}

func Contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func Remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
