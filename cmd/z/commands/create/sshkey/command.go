package sshkey

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"zctl/internal/pathutils"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/zephinzer/go-commander"
)

var conf = config.Map{
	"password": &config.String{
		Shorthand: "p",
		Usage:     "password for the ssh key file",
		Default:   "",
	},
	"overwrite": &config.Bool{
		Usage: "specify this flag to overwrite existing files",
	},
	"ssh-dir": &config.String{
		Shorthand: "d",
		Usage:     "directory to save the keys to",
		Default:   "~/.ssh/",
	},
	"state": &config.String{
		Shorthand: "s",
		Usage:     "state specified in tls cert's subject field",
		Default:   "Singapore",
	},
}

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "sshkey <key-name>",
		Short: "Creates an SSH key-pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputDirectory, err := pathutils.ResolveFullPath(conf.GetString("ssh-dir"))
			if err != nil {
				return fmt.Errorf("failed to resolve path[%s]", conf.GetString("ssh-dir"))
			}
			if err := pathutils.EnsureDirectoryExists(outputDirectory); err != nil {
				return fmt.Errorf("failed to ensure directory exists: %s", err)
			}
			if len(args) == 0 {
				return fmt.Errorf("failed to receive a key name in the arguments")
			}
			keyName := args[0]
			nameRegexString := "[a-zA-Z0-9_]+"
			nameRegex := regexp.MustCompile(nameRegexString)
			if !nameRegex.MatchString(keyName) {
				return fmt.Errorf("provided name should match regex '%s'", nameRegexString)
			}
			keyPath := path.Join(outputDirectory, fmt.Sprintf("id_rsa_%s", keyName))
			if err := pathutils.AssertFileDoesntExist(keyPath); err != nil {
				if !conf.GetBool("overwrite") {
					return fmt.Errorf("file[%s] already exists, use the --overwrite flag to force an overwrite (not recommended)", keyPath)
				}
			}
			keyPassword := conf.GetString("password")
			isUsingPasswordString := "no"
			isUsingPassword := len(keyPassword) > 0
			if isUsingPassword {
				isUsingPasswordString = "yes"
			}

			toRun := commander.NewCommand("ssh-keygen").
				AddParam("-b", "8192").
				AddParam("-t", "rsa").
				AddParam("-f", keyPath).
				AddParam("-q").
				AddParam("-N", keyPassword)
			commandString := strings.ReplaceAll(toRun.GetAsString(), keyPassword, "[REDACTED]")
			log.Printf("generating ssh keys at path[%s] (using password: %s)...\n%s", keyPath, isUsingPasswordString, commandString)
			output := toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to generate ca cert and key: %s", output.Error)
			}

			publicKey, err := os.ReadFile(keyPath + ".pub")
			if err != nil {
				return fmt.Errorf("failed to read the public key file: %s", err)
			}
			log.Printf("public key follows:\n%s", string(publicKey))

			return nil
		},
	}
	conf.ApplyToCobra(&command)
	return &command
}
