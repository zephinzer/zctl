package deps

import (
	"fmt"
	"log"
	"os"
	"zctl/cmd/z/commands/create"
	"zctl/internal/projutils"

	"github.com/spf13/cobra"
	"github.com/usvc/go-config"
	"github.com/zephinzer/go-commander"
)

var conf = config.Map{
	"production": &config.Bool{
		Shorthand: "P",
		Usage:     "when specified, only production dependencies are installed (if applicable for runtime)",
	},
}

func GetCommand() *cobra.Command {
	isProductionInstall := conf.GetBool("production")
	command := cobra.Command{
		Use:     "deps",
		Aliases: []string{"aliases"},
		Short:   "Gets dependencies for the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %s", err)
			}

			/*
				golang
			*/
			if itIs, err := projutils.IsGolangProject(pwd); err != nil {
				return fmt.Errorf("failed to check project type: %s", err)
			} else if itIs {
				log.Print("detected golang project!")
				toRun := commander.NewCommand("go").
					AddParam("mod").
					AddParam("tidy").
					EnableStderr().
					EnableStdout()
				log.Printf("running `%s`...", toRun.GetAsString(true))
				if output := toRun.Execute(); output.Error != nil {
					return fmt.Errorf("failed to get dependencies: %s", output.Error)
				}
				toRun = commander.NewCommand("go").
					AddParam("mod").
					AddParam("vendor").
					EnableStderr().
					EnableStdout()
				log.Printf("running `%s`...", toRun.GetAsString(true))
				if output := toRun.Execute(); output.Error != nil {
					return fmt.Errorf("failed to get dependencies: %s", output.Error)
				}
			}

			/*
				javascript/typescript
			*/
			if itIs, err := projutils.IsJavascriptYarnProjecct(pwd); err != nil {
				return fmt.Errorf("failed to check project type: %s", err)
			} else if itIs {
				log.Print("detected javscript project using yarn!")
				toRun := commander.NewCommand("yarn").
					AddParam("install").
					EnableStderr().
					EnableStdout()
				if isProductionInstall {
					toRun = toRun.AddParam("--production=true")
				}
				log.Printf("running `%s`...", toRun.GetAsString(true))
				if output := toRun.Execute(); output.Error != nil {
					return fmt.Errorf("failed to get dependencies: %s", output.Error)
				}
			} else if itIs, err := projutils.IsJavascriptNpmProjecct(pwd); err != nil {
				return fmt.Errorf("failed to check project type: %s", err)
			} else if itIs {
				log.Print("detected javscript project using npm!")
				toRun := commander.NewCommand("npm").
					AddParam("install").
					EnableStderr().
					EnableStdout()
				if isProductionInstall {
					toRun = toRun.AddParam("--omit=dev")
				}
				log.Printf("running `%s`...", toRun.GetAsString(true))
				if output := toRun.Execute(); output.Error != nil {
					return fmt.Errorf("failed to get dependencies: %s", output.Error)
				}
			}

			if itIs, err := projutils.IsJavascriptYarnProjecct(pwd); err != nil {
			} else if itIs {
				log.Print("detected python project!")
				toRun := commander.NewCommand("pip").
					AddParam("install").
					AddParam("-r", "requirements.txt").
					EnableStderr().
					EnableStdout()
				log.Printf("running `%s`...", toRun.GetAsString(true))
				if output := toRun.Execute(); output.Error != nil {
					return fmt.Errorf("failed to get dependencies: %s", output.Error)
				}
			}
			return nil
		},
	}
	command.AddCommand(create.GetCommand())
	conf.ApplyToCobra(&command)
	return &command
}
