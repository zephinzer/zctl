package tlscert

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/zephinzer/go-commander"
)

var conf = config.Map{
	"common-names": &config.StringSlice{
		Shorthand: "n",
		Usage:     "country specified in tls cert's subject field",
		Default:   []string{"*"},
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
			outputDirectory := args[0]
			if !path.IsAbs(outputDirectory) {
				pwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current working directory: %s", err)
				}
				outputDirectory = path.Join(pwd, outputDirectory)
			}
			if lstatInfo, err := os.Lstat(outputDirectory); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
						return fmt.Errorf("failed to create path[%s]: %s", outputDirectory, err)
					}
				} else {
					return fmt.Errorf("failed to get info about path[%s]: %s", outputDirectory, err)
				}
			} else if !lstatInfo.IsDir() {
				return fmt.Errorf("failed to get a directory at path[%s]", outputDirectory)
			}
			subject := fmt.Sprintf(
				"/C=%s"+
					"/L=%s"+
					"/O=%s"+
					"/OU=%s"+
					"/CN=%s",
				conf.GetString("country"),
				conf.GetString("locality"),
				conf.GetString("org"),
				conf.GetString("org-unit"),
				strings.Join(conf.GetStringSlice("common-names"), ","),
			)

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

			toRun := commander.NewCommand("openssl").
				AddParam("req").
				AddParam("-nodes").
				AddParam("-x509").
				AddParam("-newkey", "rsa:4096").
				AddParam("-days", strconv.Itoa(validityDays)).
				AddParam("-subj", subject).
				AddParam("-keyout", caKeyPath).
				AddParam("-out", caCertPath)
			log.Printf("generating ca cert and key...\n%s", toRun.GetAsString())
			output := toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to generate ca cert and key: %s", output.Error)
			}
			toRun = commander.NewCommand("openssl").
				AddParam("req").
				AddParam("-nodes").
				AddParam("-newkey", "rsa:4096").
				AddParam("-subj", subject).
				AddParam("-keyout", serverKeyPath).
				AddParam("-out", serverCsrPath)
			log.Printf("generating server csr and key...\n%s", toRun.GetAsString())
			output = toRun.Execute()
			if output.Error != nil {
				return fmt.Errorf("failed to generate server csr and key: %s", output.Error)
			}

			toRun = commander.NewCommand("openssl").
				AddParam("x509").
				AddParam("-req").
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

/*

echo '---------------------';
echo 'signing server csr...';
openssl x509 \
  -req \
  -in ./keys2/server-csr.pem \
  -CA ./keys2/ca-cert.pem \
  -CAkey ./keys2/ca-key.pem \
  -CAcreateserial \
  -out ./keys2/server-cert.pem;

echo '-------------------------';
echo 'printing ca cert details:';
openssl x509 \
  -in ./keys2/ca-cert.pem \
  -noout \
  -text;

echo '-----------------------------';
echo 'printing server cert details:';
openssl x509 \
  -in ./keys2/server-cert.pem \
  -noout \
  -text;

*/
