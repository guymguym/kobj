package util

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Assert is used to check for errors in calls that are not expected to fail
func Assert(errs ...error) {
	if len(errs) > 1 {
		panic(errs)
	}
	if len(errs) > 0 && errs[0] != nil {
		panic(errs[0])
	}
}

// KubeConfig returns the local config
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
