package tlscert

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"zctl/internal/pathutils"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/zephinzer/go-commander"
)

var conf = config.Map{
	"alt-names": &config.StringSlice{
		Shorthand: "N",
		Usage:     "subject alt names for the tls cert",
		Default:   []string{},
	},
	"common-name": &config.String{
		Shorthand: "n",
		Usage:     "common name specified in tls cert's subject field",
		Default:   "z.com",
	},
	"country": &config.String{
		Shorthand: "c",
		Usage:     "country specified in tls cert's subject field",
		Default:   "SG",
	},
	"days": &config.Int{
		Shorthand: "d",
		Usage:     "number of days the tls cert is valid for",
		Default:   365,
	},
	"locality": &config.String{
		Shorthand: "l",
		Usage:     "locality specified in tls cert's subject field",
		Default:   "Singapore",
	},
	"org": &config.String{
		Shorthand: "o",
		Usage:     "organisation specified in tls cert's subject field",
		Default:   "z",
	},
	"org-unit": &config.String{
		Shorthand: "u",
		Usage:     "organisation unit specified in tls cert's subject field",
		Default:   "z",
	},
	"overwrite": &config.Bool{
		Usage: "specify this flag to overwrite existing files",
	},
	"password": &config.String{
		Shorthand: "p",
		Usage:     "password used to protect private key",
	},
	"state": &config.String{
		Shorthand: "s",
		Usage:     "state specified in tls cert's subject field",
		Default:   "Singapore",
	},
}

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "tlscert <target-directory>",
		Short: "Creates a CA and server certificate and key for use with TLS",
		RunE: func(cmd *cobra.Command, args []string) error {
			validityDays := conf.GetInt("days")
			if len(args) == 0 {
				return fmt.Errorf("specify the output directory as the first argument")
			}
			outputDirectory, err := pathutils.ResolveFullPath(args[0])
			if err != nil {
				return fmt.Errorf("failed to resolve path[%s]", args[0])
			}
			if err := pathutils.EnsureDirectoryExists(outputDirectory); err != nil {
				return fmt.Errorf("failed to ensure directory exists: %s", err)
			}

			caCertPath := path.Join(outputDirectory, "ca-cert.pem")
			caKeyPath := path.Join(outputDirectory, "ca-key.pem")
			serverCsrPath := path.Join(outputDirectory, "server-csr.pem")
			serverKeyPath := path.Join(outputDirectory, "server-key.pem")
			serverCertPath := path.Join(outputDirectory, "server-cert.pem")

			for _, targetPath := range []string{caCertPath, caKeyPath, serverCsrPath, serverKeyPath, serverCertPath} {
				if targetFileInfo, err := os.Lstat(targetPath); err != nil {
					if !errors.Is(err, os.ErrNotExist) {
						return fmt.Errorf("failed to get file information on path[%s]: %s", targetPath, err)
					}
				} else if !conf.GetBool("overwrite") {
					return fmt.Errorf("file[%s] at path[%s] already exists, use --overwrite to overwrite them", targetFileInfo.Name(), targetPath)
				}
			}

			openSslTemplate := map[string]map[string]string{
				"req": {
					"default_bits":       "4096",
					"distinguished_name": "distinguished_name",
					"prompt":             "no",
				},
				"distinguished_name": {
					"countryName":            conf.GetString("country"),
					"stateOrProvinceName":    conf.GetString("state"),
					"localityName":           conf.GetString("locality"),
					"organizationName":       conf.GetString("org"),
					"organizationalUnitName": conf.GetString("org-unit"),
					"commonName":             conf.GetString("common-name"),
				},
			}
			altNames := conf.GetStringSlice("alt-names")
			if len(altNames) > 0 {
				openSslTemplate["req"]["req_extensions"] = "req_extensions"
				openSslTemplate["req_extensions"] = map[string]string{
					"subjectAltName": "@alt_names",
				}
				openSslTemplate["alt_names"] = map[string]string{}
				for index, altName := range altNames {
					openSslTemplate["alt_names"][fmt.Sprintf("DNS.%v", index+1)] = altName
				}
			}
			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %s", err)
			}
			tmpFilePath := path.Join(pwd, "./.z.create.tlscert.tmp")
			if tmpFilePathInfo, err := os.Lstat(tmpFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to verify path[%s] is available: %s", tmpFilePath, err)
			} else if err == nil && !conf.GetBool("overwrite") {
				return fmt.Errorf("file[%s] already exists at path[%s], use --overwrite to ignore and overwrite the file", tmpFilePathInfo.Name(), tmpFilePath)
			}
			openSslTemplateContentsDelimited := []string{}
			for header, keyValuePair := range openSslTemplate {
				openSslTemplateContentsDelimited = append(openSslTemplateContentsDelimited, "[ "+header+" ]")
				for key, value := range keyValuePair {
					openSslTemplateContentsDelimited = append(openSslTemplateContentsDelimited, key+" = "+value)
				}
			}
			openSslTemplateContents := strings.Join(openSslTemplateContentsDelimited, "\n") + "\n\n"
			if err := os.WriteFile(tmpFilePath, []byte(openSslTemplateContents), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create tmp file at path[%s]: %s", tmpFilePath, err)
			}
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-signals
				if err := os.Remove(tmpFilePath); err != nil {
					log.Printf("failed to remove tmp file at path[%s]: %s", tmpFilePath, err)
				}
			}()
			defer os.Remove(tmpFilePath)

			log.Printf("using following openssl config:\n%s", openSslTemplateContents)
			toRun := commander.NewCommand("openssl").
				AddParam("req").
				AddParam("-nodes").
				AddParam("-x509").
				AddParam("-newkey", "rsa:4096").
				AddParam("-days", strconv.Itoa(validityDays)).
				AddParam("-config", tmpFilePath).
				AddParam("-keyout", caKeyPath).
				AddParam("-out", caCertPath)
			if conf.GetString("password") != "" {
				toRun.SetEnvironment("OPENSSL_CERT_PASSWORD", conf.GetString("password"))
				toRun.AddParam("passout", "env:OPENSSL_CERT_PASSWORD")
			}
			log.Printf("generating ca cert and key...\n%s", toRun.GetAsString())
			output := toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to generate ca cert and key: %s", output.Error)
			}

			toRun = commander.NewCommand("openssl").
				AddParam("req").
				AddParam("-nodes").
				AddParam("-config", tmpFilePath).
				AddParam("-keyout", serverKeyPath).
				AddParam("-out", serverCsrPath)
			if conf.GetString("password") != "" {
				toRun.SetEnvironment("OPENSSL_CERT_PASSWORD", conf.GetString("password"))
				toRun.AddParam("passout", "env:OPENSSL_CERT_PASSWORD")
			}
			log.Printf("generating server csr and key...\n%s", toRun.GetAsString())
			output = toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to generate server csr and key: %s", output.Error)
			}

			toRun = commander.NewCommand("openssl").
				AddParam("x509").
				AddParam("-req").
				AddParam("-days", strconv.Itoa(validityDays)).
				AddParam("-in", serverCsrPath).
				AddParam("-CA", caCertPath).
				AddParam("-CAkey", caKeyPath).
				AddParam("-CAcreateserial").
				AddParam("-out", serverCertPath)
			log.Printf("signing server csr...\n%s", toRun.GetAsString())
			output = toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to sign server csr: %s", output.Error)
			}

			return nil
		},
	}
	conf.ApplyToCobra(&command)
	return &command
}
