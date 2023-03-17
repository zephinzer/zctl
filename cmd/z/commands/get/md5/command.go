package md5

import (
	"crypto/md5"
	"encoding/hex"
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
		Usage:     "retrieves the md5 hash of the file specified here, if defined, ignores the [plaintext] argument",
	},
}

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "md5 [plaintext]",
		Short: "gets md5 hash of provided input",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileInput := conf.GetString("file-input")
			stringInput := ""
			if fileInput == "" {
				if len(args) == 0 {
					return fmt.Errorf("failed to find a valid <file-path>")
				}
				stringInput = args[0]
			}

			var bytesToHash []byte
			if fileInput != "" {
				path, err := pathutils.ResolveFullPath(fileInput)
				if err != nil {
					return fmt.Errorf("failed to resolve path[%s]: %s", args[0], err)
				}
				fileContent, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read file[%s]: %s", path, err)
				}
				bytesToHash = fileContent
			} else {
				bytesToHash = []byte(stringInput)
			}

			hash := md5.Sum(bytesToHash)
			fmt.Println(hex.EncodeToString(hash[:]))
			return nil
		},
	}
	command.AddCommand(create.GetCommand())
	conf.ApplyToCobra(&command)
	return &command
}
