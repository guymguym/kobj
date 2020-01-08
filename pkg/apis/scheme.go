package apis

import (
	"github.com/kobj-io/kobj/pkg/apis/kobj"
	"k8s.io/apimachinery/pkg/runtime"
)

// SchemeBuilder may be used to add resources to a Scheme
var SchemeBuilder = runtime.NewSchemeBuilder()

func init() {
	SchemeBuilder.Register(kobj.SchemeBuilder...)
}
