package main

import (
	"net"

	// force compilation of packages we'll later rely upon
	// _ "k8s.io/kube-aggregator/pkg/apis/apiregistration/install"
	// _ "k8s.io/kube-aggregator/pkg/apis/apiregistration/validation"
	// _ "k8s.io/kube-aggregator/pkg/client/listers/apiregistration/v1beta1"
	// _ "k8s.io/kube-aggregator/pkg/client/listers/apiregistration/internalversion"
	// _ "k8s.io/kube-aggregator/pkg/client/clientset_generated/internalclientset"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/discovery"
	apiserver "k8s.io/apiserver/pkg/server"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
)

const GroupName = "kobj.io"
const GroupVersion = "v1alpha1"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

type Kobj struct {
	metav1.TypeMeta `json:",inline"`
	// metav1.ObjectMeta `json:"metadata,omitempty"`
	Key   string `json:"key" protobuf:"bytes,1,name=key"`
	Value string `json:"value" protobuf:"bytes,2,name=value"`
}

type KobjList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Kobj `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func (obj *Kobj) DeepCopyObject() runtime.Object {
	out := *obj
	out.TypeMeta = obj.TypeMeta
	// obj.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return &out
}

func (obj *KobjList) DeepCopyObject() runtime.Object {
	out := *obj
	out.TypeMeta = obj.TypeMeta
	obj.ListMeta.DeepCopyInto(&out.ListMeta)
	if obj.Items != nil {
		out.Items = make([]Kobj, len(obj.Items))
		for i := range obj.Items {
			out.Items[i] = *obj.Items[i].DeepCopyObject().(*Kobj)
		}
	}
	return &out
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	klog.Info("KOBJ: Starting ...")

	scheme := NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	server := NewServer(scheme, codecs)

	groupVersion := metav1.GroupVersionForDiscovery{
		GroupVersion: SchemeGroupVersion.String(),
		Version:      GroupVersion,
	}
	apiGroup := metav1.APIGroup{
		Name:             GroupName,
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

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Kobj{},
		&KobjList{},
	)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)

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
