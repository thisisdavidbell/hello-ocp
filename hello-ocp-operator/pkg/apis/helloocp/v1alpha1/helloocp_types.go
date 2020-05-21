package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HelloocpSpec defines the desired state of Helloocp
type HelloocpSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// +kubebuilder:validation:Maximum=3
	// Size is the size of the memcached deployment
	Size int32 `json:"size"`

	// +kubebuilder:validation:Enum=Option1;Option2
	// An extra spec field - a string to see what happens
	SomeString string `json:"someString"`
}

// HelloocpStatus defines the observed state of Helloocp
type HelloocpStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Nodes are the names of the memcached pods
	Nodes []string `json:"nodes"`

	// An extra status field - a string to see what happens
	VersionString string `json:"versionString"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Helloocp is the Schema for the helloocps API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=helloocps,scope=Namespaced
type Helloocp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelloocpSpec   `json:"spec,omitempty"`
	Status HelloocpStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelloocpList contains a list of Helloocp
type HelloocpList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Helloocp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Helloocp{}, &HelloocpList{})
}
