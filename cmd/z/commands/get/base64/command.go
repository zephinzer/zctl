package base64

import (
	"encoding/base64"
	"fmt"
	"os"
	"zctl/cmd/z/commands/create"
	"zctl/internal/pathutils"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
)

var conf = config.Map{
	"file-input": &config.String{
		Shorthand: "i",
		Usage:     "retrieves the base64 encoding of the file specified here, if defined, ignores the [plaintext] argument",
	},
}

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:     "base64 [plaintext]",
		Aliases: []string{"b64"},
		Short:   "gets base64 encoding of provided input",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileInput := conf.GetString("file-input")
			stringInput := ""
			if fileInput == "" {
				if len(args) == 0 {
					return fmt.Errorf("failed to find a valid <file-path>")
				}
				stringInput = args[0]
			}

			var bytesToEncode []byte
			if fileInput != "" {
				path, err := pathutils.ResolveFullPath(fileInput)
				if err != nil {
					return fmt.Errorf("failed to resolve path[%s]: %s", args[0], err)
				}
				fileContent, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read file[%s]: %s", path, err)
				}
				bytesToEncode = fileContent
			} else {
				bytesToEncode = []byte(stringInput)
			}

			encoded := base64.StdEncoding.EncodeToString(bytesToEncode)
			fmt.Println(encoded)
			return nil
		},
	}
	command.AddCommand(create.GetCommand())
	conf.ApplyToCobra(&command)
	return &command
}
