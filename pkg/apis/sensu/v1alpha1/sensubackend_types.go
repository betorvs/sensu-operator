package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SensuBackendSpec defines the desired state of SensuBackend
// +k8s:openapi-gen=true
type SensuBackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.Replicas",description="The number of replicas launched by the Sensu-Operator"
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas"`
	// +kubebuilder:printcolumn:name="SensuBackendURL",type="string",JSONPath=".spec.SensuBackendURL",description="Sensu API URL launched by the Sensu-Operator"
	// +kubebuilder:validation:MaxLength=60
	// +kubebuilder:validation:MinLength=1
	SensuBackendURL string `json:"sensubackendurl"`
	// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.Image",description="Sensu Container Image launched by the Sensu-Operator"
	// +kubebuilder:validation:MaxLength=35
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`
	// +kubebuilder:printcolumn:name="DebugSensu",type="string",JSONPath=".spec.DebugSensu",description="Sensu Debug Option to startup a new instance with debug launched by the Sensu-Operator"
	DebugSensu bool `json:"debug"`
}

// SensuBackendStatus defines the observed state of SensuBackend
// +k8s:openapi-gen=true
type SensuBackendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Nodes are the names of pods
	// +listType=set
	Nodes []string `json:"nodes,omitempty"`
	// Services are the names of svcs
	// +listType=set
	Services []string `json:"services,omitempty"`
	// Token to connect to Sensu API
	Token string `json:"token"`
	// AdminToken to access Sensu API
	AdminToken string `json:"admin_token"`
	// OperatorToken to access Sensu API
	OperatorToken string `json:"operator_token"`
	// ClusterID to access Sensu API
	ClusterID string `json:"cluster_id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuBackend is the Schema for the sensubackends API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sensubackends,scope=Namespaced
type SensuBackend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SensuBackendSpec   `json:"spec,omitempty"`
	Status SensuBackendStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuBackendList contains a list of SensuBackend
type SensuBackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuBackend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SensuBackend{}, &SensuBackendList{})
}
