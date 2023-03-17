package initialise

import (
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:     "initialise",
		Aliases: []string{"init", "ini"},
		Short:   "Initialises this CLI tool, use sub-commands to initialise other shiz",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return &command
}
