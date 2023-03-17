package get

import (
	"zctl/cmd/z/commands/get/base64"
	"zctl/cmd/z/commands/get/deps"
	"zctl/cmd/z/commands/get/md5"
	"zctl/cmd/z/commands/get/sha256"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "get",
		Short: "Gets stuff",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(base64.GetCommand())
	command.AddCommand(deps.GetCommand())
	command.AddCommand(md5.GetCommand())
	command.AddCommand(sha256.GetCommand())
	return &command
}
