package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SensuAgentSpec defines the desired state of SensuAgent
// +k8s:openapi-gen=true
type SensuAgentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.Replicas",description="The number of replicas launched by the Sensu-Operator"
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas"`
	// +kubebuilder:printcolumn:name="SensuBackendWebsocket",type="string",JSONPath=".spec.SensuBackendWebsocket",description="Sensu API Websocket launched by the Sensu-Operator"
	// +kubebuilder:validation:MaxLength=60
	// +kubebuilder:validation:MinLength=1
	SensuBackendWebsocket string `json:"sensubackend_websocket,omitempty"`
	// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.Image",description="Sensu Container Image launched by the Sensu-Operator"
	// +kubebuilder:validation:MaxLength=35
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image,omitempty"`
	// +kubebuilder:printcolumn:name="LogLevel",type="string",JSONPath=".spec.LogLevel",description="Sensu Debug Option to startup a new instance with debug launched by the Sensu-Operator"
	LogLevel string `json:"log_level,omitempty"`
	// +kubebuilder:printcolumn:name="CACertificate",type="string",JSONPath=".spec.CACertificate",description="Sensu CA Certificate Secret Name created before not launched by the Sensu-Operator"
	CACertificate string `json:"ca_certificate,omitempty"`
	// +kubebuilder:printcolumn:name="CAFileName",type="string",JSONPath=".spec.CAFileName",description="Sensu CA FileName Secret Name created before not launched by the Sensu-Operator"
	CAFileName string `json:"ca_filename,omitempty"`
	// +kubebuilder:printcolumn:name="Subscriptions",type="set",JSONPath=".spec.Subscriptions",description="Sensu Agent Subscriptions list launched by the Sensu-Operator"
	// +listType=set
	Subscriptions []string `json:"subscriptions,omitempty"`
}

// SensuAgentStatus defines the observed state of SensuAgent
// +k8s:openapi-gen=true
type SensuAgentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.Status",description="Status of Asset launched by the Sensu-Operator"
	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuAgent is the Schema for the sensuagents API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sensuagents,scope=Namespaced
type SensuAgent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SensuAgentSpec   `json:"spec,omitempty"`
	Status SensuAgentStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuAgentList contains a list of SensuAgent
type SensuAgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuAgent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SensuAgent{}, &SensuAgentList{})
}
