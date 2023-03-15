package z

import (
	"zctl/cmd/z/commands/create"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "z",
		Short: "@zephinzer's developer utility tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	command.AddCommand(create.GetCommand())
	return &command
}
