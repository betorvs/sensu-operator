package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SensuCheckSpec defines the desired state of SensuCheck
// +k8s:openapi-gen=true
type SensuCheckSpec struct {
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
	// +kubebuilder:printcolumn:name="Handlers",type="set",JSONPath=".spec.Handlers",description="Sensu Check Handlers launched by the Sensu-Operator"
	// +listType=set
	Handlers []string `json:"handlers"`
	// +kubebuilder:printcolumn:name="Subscriptions",type="set",JSONPath=".spec.Subscriptions",description="Sensu Check Subscriptions launched by the Sensu-Operator"
	// +listType=set
	Subscriptions []string `json:"subscriptions"`
	// +kubebuilder:printcolumn:name="Interval",type="string",JSONPath=".spec.Interval",description="Sensu Check Interval launched by the Sensu-Operator"
	Interval int `json:"interval"`
	// +kubebuilder:printcolumn:name="Publish",type="string",JSONPath=".spec.Publish",description="Sensu Check Publish launched by the Sensu-Operator"
	Publish bool `json:"publish"`
	// +kubebuilder:printcolumn:name="RuntimeAssets",type="set",JSONPath=".spec.RuntimeAssets",description="Sensu Check RuntimeAssets launched by the Sensu-Operator"
	// +listType=set
	RuntimeAssets []string `json:"runtime_assets,omitempty"`
	// +kubebuilder:printcolumn:name="RoundRobin",type="string",JSONPath=".spec.RoundRobin",description="Sensu Check RoundRobin launched by the Sensu-Operator"
	RoundRobin bool `json:"round_robin,omitempty"`
	// +kubebuilder:printcolumn:name="ProxyEntityName",type="string",JSONPath=".spec.ProxyEntityName",description="Sensu Check ProxyEntityName launched by the Sensu-Operator"
	ProxyEntityName string `json:"proxy_entity_name,omitempty"`
	// +kubebuilder:printcolumn:name="Annotations",type="string",JSONPath=".spec.Annotations",description="Sensu Check Annotations launched by the Sensu-Operator"
	Annotations map[string]string `json:"annotations,omitempty"`
}

// SensuCheckStatus defines the observed state of SensuCheck
// +k8s:openapi-gen=true
type SensuCheckStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.Status",description="Status of Check launched by the Sensu-Operator"
	Status string `json:"status"`
	// +kubebuilder:printcolumn:name="OwnerID",type="string",JSONPath=".status.OwnerID",description="OwnerID to access Sensu API launched by the Sensu-Operator"
	OwnerID string `json:"owner_id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuCheck is the Schema for the sensuchecks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sensuchecks,scope=Namespaced
type SensuCheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SensuCheckSpec   `json:"spec,omitempty"`
	Status SensuCheckStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuCheckList contains a list of SensuCheck
type SensuCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuCheck `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SensuCheck{}, &SensuCheckList{})
}
