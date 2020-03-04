package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SensuNamespaceSpec defines the desired state of SensuNamespace
// +k8s:openapi-gen=true
type SensuNamespaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Name",type="string",JSONPath=".spec.Name",description="Sensu Namespace launched by the Sensu-Operator"
	// +kubebuilder:validation:MaxLength=30
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// +kubebuilder:printcolumn:name="SensuBackendAPI",type="string",JSONPath=".spec.SensuBackendAPI",description="Sensu Backend API"
	SensuBackendAPI string `json:"sensu_backend_api"`
}

// SensuNamespaceStatus defines the observed state of SensuNamespace
// +k8s:openapi-gen=true
type SensuNamespaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.Status",description="Status of namespace launched by the Sensu-Operator"
	Status string `json:"status"`
	// +kubebuilder:printcolumn:name="OwnerID",type="string",JSONPath=".status.OwnerID",description="OwnerID to access Sensu API launched by the Sensu-Operator"
	OwnerID string `json:"owner_id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuNamespace is the Schema for the sensunamespaces API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sensunamespaces,scope=Namespaced
type SensuNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SensuNamespaceSpec   `json:"spec,omitempty"`
	Status SensuNamespaceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuNamespaceList contains a list of SensuNamespace
type SensuNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuNamespace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SensuNamespace{}, &SensuNamespaceList{})
}
