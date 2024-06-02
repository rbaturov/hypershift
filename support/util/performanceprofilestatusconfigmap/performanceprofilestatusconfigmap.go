package performanceprofilestatusconfigmap

import (
	"encoding/json"
	"time"

	performanceprofilev2 "github.com/openshift/cluster-node-tuning-operator/pkg/apis/performanceprofile/v2"
	"sigs.k8s.io/yaml"

	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func New(controlPlaneNamespace, userClustersNamespace, nodePoolName string, opts ...func(cm *corev1.ConfigMap) error) (*corev1.ConfigMap, error) {
	cm := getTestPerformanceProfileStatusConfigMap(controlPlaneNamespace, userClustersNamespace, nodePoolName)
	for _, opt := range opts {
		if err := opt(cm); err != nil {
			return nil, err
		}
	}
	return cm, nil
}

func getTestPerformanceProfileStatusConfigMap(controlPlaneNamespace, userClustersNamespace, nodePoolName string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "perfprofile-" + nodePoolName + "-status",
			Namespace: controlPlaneNamespace,
			Labels: map[string]string{
				"hypershift.openshift.io/nto-generated-performance-profile-status": "true",
				"hypershift.openshift.io/nodePool":                                 nodePoolName,
				"hypershift.openshift.io/performanceProfileName":                   nodePoolName,
			},
			Annotations: map[string]string{
				"hypershift.openshift.io/nodePool": nodePoolName,
			},
		},
	}
}

func WithStatus(status *performanceprofilev2.PerformanceProfileStatus) func(*corev1.ConfigMap) error {
	return func(cm *corev1.ConfigMap) error {
		if err := UpdateStatus(cm, status); err != nil {
			return err
		}
		return nil
	}
}

func UpdateStatus(cm *corev1.ConfigMap, status *performanceprofilev2.PerformanceProfileStatus) error{
	encodedStatus, encodeErr := StatusToYAML(status)
	if encodeErr != nil {
		return encodeErr
	}
	data := map[string]string{"status": string(encodedStatus)}
	cm.Data = data
	return nil
}

func AvailablePerformanceProfileStatus() *performanceprofilev2.PerformanceProfileStatus {
	lastHeartbeatTime := "2024-04-18T06:55:45Z"
	lastTransitionTime := "2024-04-18T06:55:45Z"

	heartbeatTime, _ := time.Parse(time.RFC3339, lastHeartbeatTime)
	transitionTime, _ := time.Parse(time.RFC3339, lastTransitionTime)

	conditions := []conditionsv1.Condition{
		{
			LastHeartbeatTime:  metav1.Time{Time: heartbeatTime},
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Message:            "cgroup=v1;",
			Status:             "True",
			Type:               "Available",
		},
		{
			LastHeartbeatTime:  metav1.Time{Time: heartbeatTime},
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "True",
			Type:               "Upgradeable",
		},
		{
			LastHeartbeatTime:  metav1.Time{Time: heartbeatTime},
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "False",
			Type:               "Progressing",
		},
		{
			LastHeartbeatTime:  metav1.Time{Time: heartbeatTime},
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "False",
			Type:               "Degraded",
		},
	}

	runtimeClass := "performance-performance"
	tuned := "openshift-cluster-node-tuning-operator/openshift-node-performance-performance"

	return &performanceprofilev2.PerformanceProfileStatus{
		Conditions:   conditions,
		Tuned:        &tuned,
		RuntimeClass: &runtimeClass,
	}
}



func ProgressingPerformanceProfileStatus() *performanceprofilev2.PerformanceProfileStatus {
	lastTransitionTime := "2024-04-18T06:55:45Z"
	transitionTime, _ := time.Parse(time.RFC3339, lastTransitionTime)

	conditions := []conditionsv1.Condition{
		{
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "False",
			Type:               "Available",
		},
		{
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "False",
			Type:               "Upgradeable",
		},
		{
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "True",
			Type:               "Progressing",
			Reason:             "DeploymentStarting",
			Message:            "Deployment is starting",
		},
		{
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "False",
			Type:               "Degraded",
		},
	}

	runtimeClass := "performance-performance"
	tuned := "openshift-cluster-node-tuning-operator/openshift-node-performance-performance"

	return &performanceprofilev2.PerformanceProfileStatus{
		Conditions:   conditions,
		RuntimeClass: &runtimeClass,
		Tuned:        &tuned,
	}
}

func DegradedPerformanceProfileStatus() *performanceprofilev2.PerformanceProfileStatus {
	lastHeartbeatTime := "2024-04-18T06:55:45Z"
	lastTransitionTime := "2024-04-18T06:55:45Z"

	heartbeatTime, _ := time.Parse(time.RFC3339, lastHeartbeatTime)
	transitionTime, _ := time.Parse(time.RFC3339, lastTransitionTime)

	conditions := []conditionsv1.Condition{
		{
			LastHeartbeatTime:  metav1.Time{Time: heartbeatTime},
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "False",
			Type:               "Available",
		},
		{
			LastHeartbeatTime:  metav1.Time{Time: heartbeatTime},
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "True",
			Type:               "Upgradeable",
		},
		{
			LastHeartbeatTime:  metav1.Time{Time: heartbeatTime},
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "False",
			Type:               "Progressing",
		},
		{
			LastHeartbeatTime:  metav1.Time{Time: heartbeatTime},
			LastTransitionTime: metav1.Time{Time: transitionTime},
			Status:             "True",
			Type:               "Degraded",
			Reason:             "GettingTunedStatusFailed",
			Message:            "Cannot list Tuned Profiles to match with profile perfprofile-hostedcluster01",
		},
	}

	runtimeClass := "performance-performance"
	tuned := "openshift-cluster-node-tuning-operator/openshift-node-performance-performance"

	return &performanceprofilev2.PerformanceProfileStatus{
		Conditions:   conditions,
		RuntimeClass: &runtimeClass,
		Tuned:        &tuned,
	}
}

func StatusToYAML(status *performanceprofilev2.PerformanceProfileStatus) ([]byte, error) {
	jsonData, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}
	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, err
	}
	return yamlData, nil
}

func StatusFromYAML(yamlData []byte, status *performanceprofilev2.PerformanceProfileStatus) error {
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, status)
	if err != nil {
		return err
	}
	return nil
}
