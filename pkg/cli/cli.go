package cli

import (
	"github.com/kobj-io/kobj/pkg/server"
	"github.com/spf13/cobra"
	"k8s.io/component-base/logs"
)

func Run() {
	_ = NewKobjCommand().Execute()
}

func NewKobjCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kobj",
	}
	cmd.AddCommand(NewServerCommand())
	return cmd
}

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "server",
		Long: "run kobj server",
		Run:  RunServer,
	}
	return cmd
}

func RunServer(*cobra.Command, []string) {
	logs.InitLogs()
	defer logs.FlushLogs()
	server.NewServer().Run()
}
