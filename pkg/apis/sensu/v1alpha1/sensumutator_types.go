package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SensuMutatorSpec defines the desired state of SensuMutator
// +k8s:openapi-gen=true
type SensuMutatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Name",type="string",JSONPath=".spec.Name",description="Sensu Namespace launched by the Sensu-Operator"
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// +kubebuilder:printcolumn:name="SensuBackendAPI",type="string",JSONPath=".spec.SensuBackendAPI",description="Sensu Backend API"
	SensuBackendAPI string `json:"sensu_backend_api"`
	// +kubebuilder:printcolumn:name="Namespace",type="string",JSONPath=".spec.Namespace",description="Sensu Namespace launched by the Sensu-Operator"
	Namespace string `json:"namespace"`
	// +kubebuilder:printcolumn:name="Command",type="string",JSONPath=".spec.Command",description="Sensu Check Command launched by the Sensu-Operator"
	Command string `json:"command"`
	// +kubebuilder:printcolumn:name="RuntimeAssets",type="set",JSONPath=".spec.RuntimeAssets",description="Sensu Check RuntimeAssets launched by the Sensu-Operator"
	// +listType=set
	RuntimeAssets []string `json:"runtime_assets,omitempty"`
	// +kubebuilder:printcolumn:name="EnvVars",type="set",JSONPath=".spec.EnvVars",description="Sensu Check EnvVars launched by the Sensu-Operator"
	// +listType=set
	EnvVars []string `json:"env_vars,omitempty"`
	// +kubebuilder:printcolumn:name="Timeout",type="string",JSONPath=".spec.Timeout",description="Sensu Check Timeout launched by the Sensu-Operator"
	Timeout int `json:"timeout,omitempty"`
}

// SensuMutatorStatus defines the observed state of SensuMutator
// +k8s:openapi-gen=true
type SensuMutatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.Status",description="Status of Check launched by the Sensu-Operator"
	Status string `json:"status"`
	// +kubebuilder:printcolumn:name="OwnerID",type="string",JSONPath=".status.OwnerID",description="OwnerID to access Sensu API launched by the Sensu-Operator"
	OwnerID string `json:"owner_id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuMutator is the Schema for the sensumutators API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sensumutators,scope=Namespaced
type SensuMutator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SensuMutatorSpec   `json:"spec,omitempty"`
	Status SensuMutatorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuMutatorList contains a list of SensuMutator
type SensuMutatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuMutator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SensuMutator{}, &SensuMutatorList{})
}
