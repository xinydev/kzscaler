package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ZeroScaledObject struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec holds the desired state of the ZeroScaledObject (from the client).
	// +optional
	Spec ZeroScaledObjectSpec `json:"spec,omitempty"`

	// Status communicates the observed state of the ZeroScaledObject (from the controller).
	// +optional
	Status ZeroScaledObjectStatus `json:"status,omitempty"`
}

// Verify that ZeroScaledObject adheres to the appropriate interfaces.
var (
	// Check that ZeroScaledObject can be validated and can be defaulted.
	_ apis.Validatable = (*ZeroScaledObject)(nil)
	_ apis.Defaultable = (*ZeroScaledObject)(nil)

	// Check that we can create OwnerReferences to a ZeroScaledObject.
	_ kmeta.OwnerRefable = (*ZeroScaledObject)(nil)

	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*ZeroScaledObject)(nil)
)

const (
	// PodAutoscalerConditionReady is set when the revision is starting to materialize
	// runtime resources, and becomes true when those resources are ready.
	ZeroScaledObjectConditionReady = apis.ConditionReady
)

type ZeroScaledObjectSpec struct {
	Service  corev1.ObjectReference `json:"service"`
	Workload corev1.ObjectReference `json:"workload"`
	Rule     *ScaleRuleSpec         `json:"rule,omitempty"`
}

type ScaleRuleSpec struct {
	StableWindow *int `json:"stable-window,omitempty"` // seconds,default 300s
}

type ZeroScaledObjectStatus struct {
	duckv1.Status `json:",inline"`
	Replicas      *int32 `json:"replicas,omitempty"`
}

// ZeroScaledObjectList is a list of ZeroScaledObject resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ZeroScaledObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ZeroScaledObject `json:"items"`
}

// GetStatus retrieves the status of the PodAutoscaler. Implements the KRShaped interface.
func (zs *ZeroScaledObject) GetStatus() *duckv1.Status {
	return &zs.Status.Status
}
