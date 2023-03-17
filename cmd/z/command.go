package main

import (
	"zctl/cmd/z/commands/create"
	"zctl/cmd/z/commands/get"
	"zctl/cmd/z/commands/initialise"

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
	command.AddCommand(get.GetCommand())
	command.AddCommand(initialise.GetCommand())
	return &command
}
