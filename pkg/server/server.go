package server

import (
	"net"

	"github.com/kobj-io/kobj/pkg/apis"
	kobjv1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"
	genopenapi "github.com/kobj-io/kobj/pkg/apis/openapi"
	"github.com/kobj-io/kobj/pkg/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	apiserver "k8s.io/apiserver/pkg/server"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
)

type KobjServer struct {
	APIServer *apiserver.GenericAPIServer
	Scheme    *runtime.Scheme
}

const KobjResource = "kobjs"

var KobjGroupResource = kobjv1.Resource(KobjResource)

func NewServer() *KobjServer {
	klog.Info("KOBJ: Starting ...")

	scheme := NewScheme()
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

	config.OpenAPIConfig = apiserver.DefaultOpenAPIConfig(
		genopenapi.GetOpenAPIDefinitions,
		openapi.NewDefinitionNamer(scheme),
	)
	config.OpenAPIConfig.Info.Title = "Kobj"
	config.OpenAPIConfig.Info.Version = "1.0.0"
	completedConfig := config.Complete()

	s, err := completedConfig.New("kobj", delegate)
	util.Assert(err)

	storage := NewKobjMemStorage()
	strategy := &KobjStrategy{}
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &kobjv1.Kobj{} },
		NewListFunc:              func() runtime.Object { return &kobjv1.KobjList{} },
		DefaultQualifiedResource: KobjGroupResource,
		PredicateFunc:            KobjMatcher,
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
		TableConvertor: printerstorage.TableConvertor{
			TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers),
		},
		Storage: genericregistry.DryRunnableStorage{
			Storage: storage,
			Codec:   codecs.LegacyCodec(kobjv1.SchemeGroupVersion),
		},
	}
	util.Assert(store.CompleteWithOptions(&generic.StoreOptions{
		RESTOptions: generic.RESTOptions{ResourcePrefix: KobjResource},
		AttrFunc:    KobjGetAttrs,
	}))

	groupInfo := apiserver.NewDefaultAPIGroupInfo(kobjv1.GroupName, scheme, runtime.NewParameterCodec(scheme), codecs)
	groupInfo.VersionedResourcesStorageMap[kobjv1.VersionName] = map[string]rest.Storage{KobjResource: store}
	groupInfo.MetaGroupVersion = &kobjv1.SchemeGroupVersion
	_ = s.InstallAPIGroup(&groupInfo)

	return &KobjServer{
		APIServer: s,
		Scheme:    scheme,
	}
}

func (s *KobjServer) Run() error {
	klog.Info("KOBJ: Running API server ...")
	stopChan := apiserver.SetupSignalHandler()
	return s.APIServer.PrepareRun().Run(stopChan)
}

func NewScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	util.Assert(apis.AddToScheme(scheme))
	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})
	return scheme
}
