package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"
	"github.com/openshift/hypershift/hypershift-operator/controllers/manifests"
	. "github.com/openshift/hypershift/hypershift-operator/controllers/nodepool"

	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	performanceprofilev2 "github.com/openshift/cluster-node-tuning-operator/pkg/apis/performanceprofile/v2"
	"k8s.io/utils/pointer"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type NTOPerformanceProfileStatusTest struct {
	ctx                 context.Context
	managementClient    crclient.Client
	hostedClusterClient crclient.Client
	hostedCluster       *hyperv1.HostedCluster
}

func NewNTOPerformanceProfileStatusTest(ctx context.Context, mgmtClient crclient.Client, hostedCluster *hyperv1.HostedCluster, hcClient crclient.Client) *NTOPerformanceProfileStatusTest {
	return &NTOPerformanceProfileStatusTest{
		ctx:                 ctx,
		hostedCluster:       hostedCluster,
		hostedClusterClient: hcClient,
		managementClient:    mgmtClient,
	}
}

func (mc *NTOPerformanceProfileStatusTest) Setup(t *testing.T) {
	t.Log("Starting test NTOPerformanceProfileStatusTest")
}

func (mc *NTOPerformanceProfileStatusTest) BuildNodePoolManifest(defaultNodepool hyperv1.NodePool) (*hyperv1.NodePool, error) {
	nodePool := &hyperv1.NodePool{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mc.hostedCluster.Name,
			Namespace: mc.hostedCluster.Namespace,
		},
	}
	defaultNodepool.Spec.DeepCopyInto(&nodePool.Spec)
	nodePool.Spec.Replicas = pointer.Int32(1) //Change this to prev

	return nodePool, nil
}

func (mc *NTOPerformanceProfileStatusTest) Run(t *testing.T, nodePool hyperv1.NodePool, nodes []corev1.Node) {
	t.Log("Entering NTO PerformanceProfileStatus test")
	g := NewWithT(t)

	ctx := mc.ctx
	controlPlaneNamespace := manifests.HostedControlPlaneNamespace(mc.hostedCluster.Namespace, mc.hostedCluster.Name)
	t.Logf("Hosted control plane namespace is %s", controlPlaneNamespace)

	// we can't assume the configmap is not there
	// cond := FindStatusCondition(nodePool.Status.Conditions, hyperv1.NodePoolPerformanceProfileAppliedSuccessfullyType)
	// g.Expect(cond).To(BeNil(), "PerformanceProfileAppliedSuccessfully is present although no PerformanceProfileConfigMap applied")

	performanceProfileStatusConfigMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "perfprofile-hostedcluster01-status",
			Namespace: controlPlaneNamespace,
			Labels: map[string]string{
				"hypershift.openshift.io/nto-generated-machine-config": "true", //Should be changed to nto-generated-performance-profile-status
				"hypershift.openshift.io/nodePool":                     nodePool.Name,
				"hypershift.openshift.io/performanceProfileName":       nodePool.Name,
			},
			Annotations: map[string]string{
				"hypershift.openshift.io/nodePool": nodePool.Name,
			},
		},
		Data: map[string]string{
			"status": "{\"conditions\":[{\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"False\",\"type\":\"Available\"},{\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"False\",\"type\":\"Upgradeable\"},{\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"True\",\"type\":\"Progressing\",\"reason\":\"DeploymentStarting\",\"message\":\"Deployment is starting\"},{\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"False\",\"type\":\"Degraded\"}],\"runtimeClass\":\"performance-performance\",\"tuned\":\"openshift-cluster-node-tuning-operator/openshift-node-performance-performance\"}",
		},
	}
	t.Log("Creating the performance profile status configMap with progressing status")

	if err := mc.managementClient.Create(ctx, performanceProfileStatusConfigMap); err != nil {
		if !errors.IsAlreadyExists(err) {
			t.Fatalf("failed to create configmap for PerformanceProfileStatus %v", err)
		}
	}

	expectedCondition := &hyperv1.NodePoolCondition{
		Type:    hyperv1.NodePoolPerformanceProfileAppliedSuccessfullyType,
		Status:  corev1.ConditionFalse,
		Reason:  "DeploymentStarting",
		Message: "Deployment is starting",
	}
	g.Eventually(func(gg Gomega) {
		cond := FindStatusCondition(nodePool.Status.Conditions, hyperv1.NodePoolPerformanceProfileAppliedSuccessfullyType)
		g.Expect(cond).ToNot(BeNil())
		g.Expect(cond.Status).To(Equal(expectedCondition.Status))
		g.Expect(cond.Message).To(Equal(expectedCondition.Message))
		g.Expect(cond.Reason).To(Equal(expectedCondition.Reason))
	}).Within(1 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

	CM := performanceProfileStatusConfigMap.DeepCopy()
	CM.Data = getDegradedStatus()

	t.Log("Patching performance profile status configmap to have degraded status")
	if err := mc.managementClient.Patch(ctx, performanceProfileStatusConfigMap, crclient.MergeFrom(CM)); err != nil {
		t.Fatalf("failed to patch configmap for PerformanceProfileStatus %v", err)
	}

	expectedCondition = &hyperv1.NodePoolCondition{
		Type:    hyperv1.NodePoolPerformanceProfileAppliedSuccessfullyType,
		Status:  corev1.ConditionFalse,
		Reason:  "GettingTunedStatusFailed",
		Message: "Cannot list Tuned Profiles to match with profile perfprofile-hostedcluster01",
	}

	g.Eventually(func(gg Gomega) {
		cond := FindStatusCondition(nodePool.Status.Conditions, hyperv1.NodePoolPerformanceProfileAppliedSuccessfullyType)
		g.Expect(cond).ToNot(BeNil())
		g.Expect(cond.Status).To(Equal(expectedCondition.Status))
		g.Expect(cond.Message).To(Equal(expectedCondition.Message))
		g.Expect(cond.Reason).To(Equal(expectedCondition.Reason))
	}).Within(1 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

	CM = performanceProfileStatusConfigMap.DeepCopy()
	CM.Data = GetAvailableStatus()

	t.Log("Patching performance profile status configmap to have available status")
	if err := mc.managementClient.Patch(ctx, performanceProfileStatusConfigMap, crclient.MergeFrom(CM)); err != nil {
		t.Fatalf("failed to update configmap for PerformanceProfileStatus %v", err)
	}

	expectedCondition = &hyperv1.NodePoolCondition{
		Type:    hyperv1.NodePoolPerformanceProfileAppliedSuccessfullyType,
		Status:  corev1.ConditionTrue,
		Message: "cgroup=v1;",
		Reason:  hyperv1.AsExpectedReason,
	}

	g.Eventually(func(gg Gomega) {
		cond := FindStatusCondition(nodePool.Status.Conditions, hyperv1.NodePoolPerformanceProfileAppliedSuccessfullyType)
		g.Expect(cond).ToNot(BeNil())
		g.Expect(cond.Status).To(Equal(expectedCondition.Status))
		g.Expect(cond.Message).To(Equal(expectedCondition.Message))
		g.Expect(cond.Reason).To(Equal(expectedCondition.Reason))
	}).Within(1 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

	if err := mc.managementClient.Delete(ctx, performanceProfileStatusConfigMap); err != nil {
		t.Logf("failed to delete configmap for PerformanceProfile object: %v", err)
	}

	g.Eventually(func(gg Gomega) {
		cond := FindStatusCondition(nodePool.Status.Conditions, hyperv1.NodePoolPerformanceProfileAppliedSuccessfullyType)
		g.Expect(cond).To(BeNil())
	}).Within(1 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

	t.Log("Ending NTO PerformanceProfileStatus test: OK")

}

func GetAvailPerformanceProfileStatus() *performanceprofilev2.PerformanceProfileStatus {
	// Define the hardcoded values
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


// func GetAvailableStatus() (map[string]string, error) {

// 	availableStatus := GetAvailPerformanceProfileStatus()
// 	encodedStatus, encodeErr := StatusToYAML(availableStatus)
// 	if encodeErr != nil {
// 		return nil, encodeErr
// 	}
// 	data := map[string]string {"status": string(encodedStatus)}
// 	return data, nil
// }

// func getProgressingStatus() map[string]string {
// 	progressingStatus := map[string]string{
// 		"status": "{\"conditions\":[{\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"False\",\"type\":\"Available\"},{\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"False\",\"type\":\"Upgradeable\"},{\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"True\",\"type\":\"Progressing\",\"reason\":\"DeploymentStarting\",\"message\":\"Deployment is starting\"},{\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"False\",\"type\":\"Degraded\"}],\"runtimeClass\":\"performance-performance\",\"tuned\":\"openshift-cluster-node-tuning-operator/openshift-node-performance-performance\"}",
// 	}
// 	return progressingStatus
// }

// func getDegradedStatus() map[string]string {
// 	degradedStatus := map[string]string{
// 		"status": "{\"conditions\":[{\"lastHeartbeatTime\":\"2024-04-18T06:55:45Z\",\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"False\",\"type\":\"Available\"},{\"lastHeartbeatTime\":\"2024-04-18T06:55:45Z\",\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"True\",\"type\":\"Upgradeable\"},{\"lastHeartbeatTime\":\"2024-04-18T06:55:45Z\",\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"False\",\"type\":\"Progressing\"},{\"lastHeartbeatTime\":\"2024-04-18T06:55:45Z\",\"lastTransitionTime\":\"2024-04-18T06:55:45Z\",\"status\":\"True\",\"type\":\"Degraded\",\"reason\":\"GettingTunedStatusFailed\",\"message\":\"Cannot list Tuned Profiles to match with profile perfprofile-hostedcluster01\"}],\"runtimeClass\":\"performance-performance\",\"tuned\":\"openshift-cluster-node-tuning-operator/openshift-node-performance-performance\"}",
// 	}
// 	return degradedStatus
// }


