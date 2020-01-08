// Package kobj ...
package kobj

import (
	"github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

// SchemeBuilder may be used to add resources to a Scheme
var SchemeBuilder = runtime.NewSchemeBuilder()

func init() {
	SchemeBuilder.Register(v1alpha1.SchemeBuilder...)
}
