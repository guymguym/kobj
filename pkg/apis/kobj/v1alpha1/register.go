package v1alpha1

import (
	"github.com/kobj-io/kobj/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// GroupName is the api group name
	GroupName = "kobj.io"

	// VersionName is the api group version
	VersionName = "v1alpha1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: VersionName}

	// InternalGroupVersion is group version used to register these objects
	InternalGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

	// SchemeBuilder is a list of build functions to apply to a scheme
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// Keep a local pointer to the array because the methods are defined on pointer type
	localSchemeBuilder = &SchemeBuilder

	// AddToScheme adds all Resources to the Scheme
	AddToScheme = localSchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(InternalGroupVersion,
		&Kobj{},
		&KobjList{},
	)
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Kobj{},
		&KobjList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	util.Assert(scheme.SetVersionPriority(SchemeGroupVersion))
	return nil
}
