package server

import (
	"github.com/spf13/cobra"
	"k8s.io/component-base/logs"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "server",
		Long: "run object server",
		Run:  RunServer,
	}
	return cmd
}

func RunServer(*cobra.Command, []string) {
	logs.InitLogs()
	defer logs.FlushLogs()
	NewServer().Run()
}
