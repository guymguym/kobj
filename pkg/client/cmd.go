package client

import (
	"io/ioutil"
	"os"

	kobjclient "github.com/kobj-io/kobj/pkg/apis/clientset/versioned"
	kobjv1 "github.com/kobj-io/kobj/pkg/apis/kobj/v1alpha1"
	"github.com/kobj-io/kobj/pkg/util"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
)

func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Long: "get object",
		Run:  RunGet,
	}
	cmd.Flags().StringP("namespace", "n", "", "Set Namespace")
	return cmd
}

func NewPutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "put",
		Long: "put object",
		Run:  RunPut,
	}
	cmd.Flags().StringP("namespace", "n", "", "Set Namespace")
	return cmd
}

func RunGet(cmd *cobra.Command, args []string) {
	logs.InitLogs()
	defer logs.FlushLogs()

	client := kobjclient.NewForConfigOrDie(util.KubeConfig()).KobjV1alpha1()
	ns, _ := cmd.Flags().GetString("namespace")
	name := args[0]

	ko, err := client.Kobjs(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	klog.Infof("Get: %s/%s received %d bytes", ko.Namespace, ko.Name, len(ko.Data))
	if ko.Data != nil {
		os.Stdout.Write(ko.Data)
	}
}

func RunPut(cmd *cobra.Command, args []string) {
	logs.InitLogs()
	defer logs.FlushLogs()

	client := kobjclient.NewForConfigOrDie(util.KubeConfig()).KobjV1alpha1()
	ns, _ := cmd.Flags().GetString("namespace")
	name := args[0]

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		klog.Fatal(err)
	}

	ko := &kobjv1.Kobj{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Data: data,
	}

	klog.Infof("Put: %s/%s sending %d bytes", ko.Namespace, ko.Name, len(ko.Data))
	ko, err = client.Kobjs(ns).Create(ko)
	if err != nil {
		klog.Fatal(err)
	}

	p := printers.YAMLPrinter{}
	p.PrintObj(ko, os.Stdout)
}
