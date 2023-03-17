package helm

import (
	"zctl/cmd/z/commands/create/helm/deployment"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "helm",
		Short: "Creates k8s manifest templates for use with Helm in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(deployment.GetCommand())
	return &command
}
