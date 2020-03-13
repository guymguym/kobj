package server

import (
	"fmt"
	"net"
	goruntime "runtime"

	"github.com/kobj-io/kobj/pkg/apis"
	kobjv1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"
	genopenapi "github.com/kobj-io/kobj/pkg/generated/openapi"
	"github.com/kobj-io/kobj/pkg/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/registry/rest"
	apiserver "k8s.io/apiserver/pkg/server"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/klog"
)

type KobjServer struct {
	APIServer *apiserver.GenericAPIServer
	Scheme    *runtime.Scheme
}

const KobjResource = "kobjs"

var KobjGroupResource = kobjv1.Resource(KobjResource)

// NewServer creates a new instance of the kobj api server
func NewServer() *KobjServer {
	klog.Info("KOBJ: Starting ...")

	scheme := runtime.NewScheme()
	util.Assert(apis.AddToScheme(scheme))
	// we need to add the options to empty v1
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})
	// TODO: keep the generic API server from wanting this
	scheme.AddUnversionedTypes(schema.GroupVersion{Group: "", Version: "v1"},
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)

	codecs := serializer.NewCodecFactory(scheme)
	config := apiserver.NewRecommendedConfig(codecs)
	delegate := apiserver.NewEmptyDelegate()
	localhostIP := net.ParseIP("127.0.0.1")

	options := apiserveroptions.NewServerRunOptions()
	options.AdvertiseAddress = localhostIP
	// options.Authorization.RemoteKubeConfigFileOptional = true
	// options.Authentication.RemoteKubeConfigFileOptional = true
	util.Assert(options.Validate()...)
	util.Assert(options.ApplyTo(&config.Config))

	sso := apiserveroptions.NewSecureServingOptions().WithLoopback()
	sso.BindPort = 8443
	util.Assert(sso.MaybeDefaultWithSelfSignedCerts(
		"localhost", nil, []net.IP{localhostIP},
	))
	util.Assert(sso.Validate()...)
	util.Assert(sso.ApplyTo(&config.SecureServing, &config.LoopbackClientConfig))

	var (
		// these come from ldflags
		gitVersion   = "v0.0.0-master+$Format:%h$"
		gitCommit    = "$Format:%H$"          // sha1 from git, output of $(git rev-parse HEAD)
		gitTreeState = ""                     // state of git tree, either "clean" or "dirty"
		buildDate    = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	)
	config.Version = &version.Info{
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    goruntime.Version(),
		Compiler:     goruntime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", goruntime.GOOS, goruntime.GOARCH),
	}

	config.OpenAPIConfig = apiserver.DefaultOpenAPIConfig(
		genopenapi.GetOpenAPIDefinitions,
		openapi.NewDefinitionNamer(scheme),
	)
	config.OpenAPIConfig.Info.Title = "Kobj"
	config.OpenAPIConfig.Info.Version = "1.0.0"

	completedConfig := config.Complete()

	s, err := completedConfig.New("kobj", delegate)
	util.Assert(err)

	groupInfo := apiserver.NewDefaultAPIGroupInfo(kobjv1.GroupName, scheme, runtime.NewParameterCodec(scheme), codecs)
	groupInfo.VersionedResourcesStorageMap[kobjv1.VersionName] = map[string]rest.Storage{
		// KobjResource: NewKobjMemRESTStorage(),
		KobjResource: NewKobjStorage(),

		// TODO split to data subresource
		// see https://github.com/operator-framework/operator-lifecycle-manager/blob/34f3888376fb24c75bc97a8f4e2b99b6a78ba801/pkg/package-server/apiserver/generic/storage.go#L68-L71
	}
	groupInfo.MetaGroupVersion = &kobjv1.SchemeGroupVersion
	_ = s.InstallAPIGroup(&groupInfo)

	return &KobjServer{
		APIServer: s,
		Scheme:    scheme,
	}
}

// Run the api server
func (s *KobjServer) Run() error {
	klog.Info("KOBJ: Running API server ...")
	stopChan := apiserver.SetupSignalHandler()
	return s.APIServer.PrepareRun().Run(stopChan)
}
