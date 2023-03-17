package deployment

import (
	_ "embed"
	"fmt"
	"os"
	"zctl/internal/pathutils"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
)

//go:embed deployment.yaml
var deploymentTemplate []byte

var conf = config.Map{
	"filename": &config.String{
		Usage:     "defines a custom output filename for the resource",
		Shorthand: "f",
		Default:   "deployment.yaml",
	},
	"overwrite": &config.Bool{
		Usage: "specify this flag to overwrite existing files",
	},
}

func GetCommand() *cobra.Command {
	command := cobra.Command{
		Use:   "deployment",
		Short: "Creates a k8s Deployment resource manifest in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := conf.GetString("filename")
			outputPath, err := pathutils.ResolveFullPath(filename)
			if err != nil {
				return fmt.Errorf("failed to resolve path[%s]", filename)
			}
			if err := pathutils.AssertFileDoesntExist(outputPath); err != nil {
				if !conf.GetBool("overwrite") {
					return fmt.Errorf("file[%s] already exists, use the --overwrite flag to force an overwrite (not recommended)", outputPath)
				}
			}
			if err := os.WriteFile(outputPath, deploymentTemplate, os.ModePerm); err != nil {
				return fmt.Errorf("failed to write file to path[%s]: %s", outputPath, err)
			}
			return nil
		},
	}
	return &command
}
