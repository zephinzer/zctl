package gpgkey

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/zephinzer/go-commander"
)

var conf = config.Map{
	"days": &config.Int{
		Shorthand: "d",
		Usage:     "number of days the gpg key is valid for",
		Default:   365,
	},
	"name": &config.String{
		Shorthand: "n",
		Usage:     "your human-logical name",
		Default:   "z",
	},
	"comment": &config.String{
		Shorthand: "c",
		Usage:     "machine-use name",
		Default:   "z",
	},
	"email": &config.String{
		Shorthand: "r",
		Usage:     "your email address",
		Default:   "z@z.com",
	},
	"overwrite": &config.Bool{
		Usage: "specify this flag to overwrite existing files",
	},
	"state": &config.String{
		Shorthand: "s",
		Usage:     "state specified in tls cert's subject field",
		Default:   "Singapore",
	},
}

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "gpgkey",
		Short: "Creates a GPG key and adds it to your local keyring",
		RunE: func(cmd *cobra.Command, args []string) error {
			validityDays := conf.GetInt("days")
			nameReal := conf.GetString("name")
			nameComment := conf.GetString("comment")
			nameEmail := conf.GetString("email")

			gpgTemplateFileContents := strings.Builder{}
			gpgTemplateFileContents.WriteString("Key-Type: RSA\n")
			gpgTemplateFileContents.WriteString("Key-Length: 4096\n")
			gpgTemplateFileContents.WriteString("Subkey-Type: RSA\n")
			gpgTemplateFileContents.WriteString("Subkey-Length: 4096\n")
			gpgTemplateFileContents.WriteString(fmt.Sprintf("Name-Real: %s\n", nameReal))
			gpgTemplateFileContents.WriteString(fmt.Sprintf("Name-Comment: %s\n", nameComment))
			gpgTemplateFileContents.WriteString(fmt.Sprintf("Name-Email: %s\n", nameEmail))
			gpgTemplateFileContents.WriteString(fmt.Sprintf("Expire-Date: %v\n", validityDays))
			gpgTemplateFileContents.WriteString("%commit\n")

			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %s", err)
			}
			tmpFilePath := path.Join(pwd, "./.z.create.gpgkey.tmp")
			if tmpFilePathInfo, err := os.Lstat(tmpFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to verify path[%s] is available: %s", tmpFilePath, err)
			} else if err == nil && !conf.GetBool("overwrite") {
				return fmt.Errorf("file[%s] already exists at path[%s], use --overwrite to ignore and overwrite the file", tmpFilePathInfo.Name(), tmpFilePath)
			}
			if err := os.WriteFile(tmpFilePath, []byte(gpgTemplateFileContents.String()), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create tmp file at path[%s]: %s", tmpFilePath, err)
			}
			defer os.Remove(tmpFilePath)

			toRun := commander.NewCommand("gpg").
				AddParam("--generate-key").
				AddParam("--batch", tmpFilePath)
			log.Printf("using following content for gpg key generation:\n```%s\n```", gpgTemplateFileContents.String())
			log.Printf("generating gpg key...\n%s", toRun.GetAsString())
			output := toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to generate gpg key: %s", err)
			}

			var gpgKeysListOutputStream bytes.Buffer
			toRun = commander.NewCommand("gpg").
				AddParam("--list-secret-keys").
				AddParam("--keyid-format", "LONG").
				SetStdout(&gpgKeysListOutputStream)
			log.Printf("retrieving key id...\n%s", toRun.GetAsString())
			output = toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to get list of gpg keys: %s", err)
			}
			var keyId string
			gpgKeysList := strings.Split(gpgKeysListOutputStream.String(), "\n")
			for i := len(gpgKeysList) - 1; i >= 0; i-- {
				line := gpgKeysList[i]
				if strings.Contains(line, nameEmail) {
					keyIdLine := gpgKeysList[i-2]
					r := regexp.MustCompile("rsa[0-9]{4,}/([A-F0-9]+)")
					matches := r.FindAllStringSubmatch(keyIdLine, 1)
					keyId = matches[0][1]
					break
				}
			}
			log.Printf("generated key id: %s", keyId)

			toRun = commander.NewCommand("gpg").
				AddParam("--armor").
				AddParam("--export", keyId)
			log.Printf("exporting public key...\n%s", toRun.GetAsString())
			output = toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to export gpg key: %s", err)
			}
			log.Printf("public key is as follows:\n%s", output.Stdout.String())

			return nil
		},
	}
	conf.ApplyToCobra(&command)
	return &command
}
