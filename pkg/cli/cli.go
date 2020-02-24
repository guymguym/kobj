package cli

import (
	"github.com/kobj-io/kobj/pkg/client"
	"github.com/kobj-io/kobj/pkg/server"
	"github.com/spf13/cobra"
)

func Run() error {
	return NewCommand().Execute()
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kobj",
	}
	cmd.AddCommand(server.NewCommand())
	cmd.AddCommand(client.NewGetCommand())
	cmd.AddCommand(client.NewPutCommand())
	return cmd
}
