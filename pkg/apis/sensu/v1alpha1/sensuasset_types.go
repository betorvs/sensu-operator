package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SensuAssetSpec defines the desired state of SensuAsset
// +k8s:openapi-gen=true
type SensuAssetSpec struct {
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
	// +kubebuilder:printcolumn:name="AssetURL",type="string",JSONPath=".spec.AssetURL",description="Sensu Asset URL to downlaod launched by the Sensu-Operator"
	AssetURL string `json:"asseturl"`
	// +kubebuilder:printcolumn:name="Sha512",type="string",JSONPath=".spec.Sha512",description="Sensu Asset SHA512 launched by the Sensu-Operator"
	Sha512 string `json:"sha512"`
}

// SensuAssetStatus defines the observed state of SensuAsset
// +k8s:openapi-gen=true
type SensuAssetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.Status",description="Status of Asset launched by the Sensu-Operator"
	Status string `json:"status"`
	// +kubebuilder:printcolumn:name="OwnerID",type="string",JSONPath=".status.OwnerID",description="OwnerID to access Sensu API launched by the Sensu-Operator"
	OwnerID string `json:"owner_id"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuAsset is the Schema for the sensuassets API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sensuassets,scope=Namespaced
type SensuAsset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SensuAssetSpec   `json:"spec,omitempty"`
	Status SensuAssetStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuAssetList contains a list of SensuAsset
type SensuAssetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuAsset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SensuAsset{}, &SensuAssetList{})
}
