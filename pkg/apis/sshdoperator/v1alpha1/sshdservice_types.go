package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SshdServiceSpec defines the desired state of SshdService
type SshdServiceSpec struct {
        UserName string `json:"username,omitempty"`
}

// SshdServiceStatus defines the observed state of SshdService
type SshdServiceStatus struct {
        Stage string `json:"stage,omitempty"`
        Password string `json:"password,omitempty"`
        NodePort int `json:"nodeport,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SshdService is the Schema for the sshdservices API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sshdservices,scope=Namespaced
type SshdService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SshdServiceSpec   `json:"spec,omitempty"`
	Status SshdServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SshdServiceList contains a list of SshdService
type SshdServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SshdService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SshdService{}, &SshdServiceList{})
}
