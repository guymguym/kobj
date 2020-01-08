package main

import (
	"net"

	apis "github.com/kobj-io/kobj/pkg/apis"
	kobjapi "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"
	genopenapi "github.com/kobj-io/kobj/pkg/apis/openapi"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/discovery"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	apiserver "k8s.io/apiserver/pkg/server"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	klog.Info("KOBJ: Starting ...")

	scheme := NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	server := NewServer(scheme, codecs)

	groupVersion := metav1.GroupVersionForDiscovery{
		GroupVersion: kobjapi.SchemeGroupVersion.String(),
		Version:      kobjapi.GroupVersion,
	}
	apiGroup := metav1.APIGroup{
		Name:             kobjapi.GroupName,
		Versions:         []metav1.GroupVersionForDiscovery{groupVersion},
		PreferredVersion: groupVersion,
	}
	// groupInfo := apiserver.NewDefaultAPIGroupInfo(GroupName, scheme, runtime.NewParameterCodec(scheme), codecs)
	server.DiscoveryGroupManager.AddGroup(apiGroup)
	container := server.Handler.GoRestfulContainer
	container.Add(discovery.NewAPIGroupHandler(server.Serializer, apiGroup).WebService())

	klog.Info("KOBJ: Running API server ...")
	stopChan := apiserver.SetupSignalHandler()
	server.PrepareRun().Run(stopChan)
}

func NewScheme() *runtime.Scheme {

	scheme := runtime.NewScheme()

	apis.SchemeBuilder.AddToScheme(scheme)

	// metav1.AddToGroupVersion(scheme, kobjapi.SchemeGroupVersion)
	// scheme.AddKnownTypes(kobjapi.SchemeGroupVersion,
	// 	&kobjapi.Kobj{},
	// 	&kobjapi.KobjList{},
	// )

	// // we need to add the options to empty v1
	// // TODO fix the server code to avoid this
	// metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})

	// // TODO: keep the generic API server from wanting this
	// unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	// scheme.AddUnversionedTypes(unversioned,
	// 	&metav1.Status{},
	// 	&metav1.APIVersions{},
	// 	&metav1.APIGroupList{},
	// 	&metav1.APIGroup{},
	// 	&metav1.APIResourceList{},
	// )

	return scheme
}

func NewServer(scheme *runtime.Scheme, codecs serializer.CodecFactory) *apiserver.GenericAPIServer {
	config := apiserver.NewRecommendedConfig(codecs)
	delegate := apiserver.NewEmptyDelegate()

	options := apiserveroptions.NewServerRunOptions()
	options.AdvertiseAddress = net.ParseIP("127.0.0.1")
	// options.Authorization.RemoteKubeConfigFileOptional = true
	// options.Authentication.RemoteKubeConfigFileOptional = true

	if errs := options.Validate(); len(errs) > 0 {
		panic(errs)
	}

	if err := options.ApplyTo(&config.Config); err != nil {
		panic(err)
	}

	sso := apiserveroptions.NewSecureServingOptions().WithLoopback()
	sso.BindPort = 8443

	// TODO have a "real" external address (have an AdvertiseAddress?)
	if err := sso.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		panic(err)
	}

	if errs := sso.Validate(); len(errs) > 0 {
		panic(errs)
	}

	if err := sso.ApplyTo(&config.SecureServing, &config.LoopbackClientConfig); err != nil {
		panic(err)
	}

	config.OpenAPIConfig = apiserver.DefaultOpenAPIConfig(genopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(scheme))
	config.OpenAPIConfig.Info.Title = "Kobj"
	config.OpenAPIConfig.Info.Version = "1.0.0"

	completedConfig := config.Complete()
	server, err := completedConfig.New("kobj", delegate)
	if err != nil {
		panic(err)
	}

	return server
}

func KubeConfig() *rest.Config {
	kubeConfig, err :=
		clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		).ClientConfig()
	if err != nil {
		panic(err)
	}

	return kubeConfig
}
