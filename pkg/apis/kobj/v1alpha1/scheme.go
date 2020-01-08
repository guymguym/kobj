package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// GroupName is the api group name
	GroupName = "kobj.io"

	// GroupVersion is the api group version
	GroupVersion = "v1alpha1"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

// SchemeBuilder may be used to add resources to a Scheme
var SchemeBuilder = runtime.NewSchemeBuilder()

// Register a new type to the scheme builder
func Register(object ...runtime.Object) {
	SchemeBuilder.Register(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(SchemeGroupVersion, object...)
		metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
		return nil
	})
}
