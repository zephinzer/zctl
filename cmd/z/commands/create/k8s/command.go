package k8s

import (
	"zctl/cmd/z/commands/create/k8s/deployment"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "k8s",
		Short: "Creates k8s manifest templates in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(deployment.GetCommand())
	return &command
}
