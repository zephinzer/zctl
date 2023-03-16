package create

import (
	"zctl/cmd/z/commands/create/gpgkey"
	"zctl/cmd/z/commands/create/sshkey"
	"zctl/cmd/z/commands/create/tlscert"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "create",
		Short: "Creates stuff",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(gpgkey.GetCommand())
	command.AddCommand(sshkey.GetCommand())
	command.AddCommand(tlscert.GetCommand())
	return &command
}
