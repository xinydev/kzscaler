package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ZeroScaler struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec holds the desired state of the ZeroScaler (from the client).
	// +optional
	Spec ZeroScalerSpec `json:"spec,omitempty"`

	// Status communicates the observed state of the ZeroScaler (from the controller).
	// +optional
	Status ZeroScalerStatus `json:"status,omitempty"`
}

// Verify that ZeroScaler adheres to the appropriate interfaces.
var (
	// Check that ZeroScaler can be validated and can be defaulted.
	_ apis.Validatable = (*ZeroScaler)(nil)
	_ apis.Defaultable = (*ZeroScaler)(nil)

	// Check that we can create OwnerReferences to a ZeroScaler.
	_ kmeta.OwnerRefable = (*ZeroScaler)(nil)

	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*ZeroScaler)(nil)
)

const (
	// PodAutoscalerConditionReady is set when the revision is starting to materialize
	// runtime resources, and becomes true when those resources are ready.
	ZeroScalerConditionReady = apis.ConditionReady
)

type ZeroScalerSpec struct {
	Service  duckv1.KReference `json:"service"`
	Workload duckv1.KReference `json:"workload"`
}

type ZeroScalerStatus struct {
	duckv1.Status `json:",inline"`
	Replicas      *int32 `json:"replicas,omitempty"`
}

// ZeroScalerList is a list of ZeroScaler resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ZeroScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ZeroScaler `json:"items"`
}

// GetStatus retrieves the status of the PodAutoscaler. Implements the KRShaped interface.
func (zs *ZeroScaler) GetStatus() *duckv1.Status {
	return &zs.Status.Status
}
