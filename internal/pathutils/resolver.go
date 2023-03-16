package pathutils

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func ResolveFullPath(input string) (string, error) {
	output := input
	if strings.Index(input, "~") == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %s", err)
		}
		output = path.Join(homeDir, strings.Replace(input, "~", "", 1))
	}
	if !path.IsAbs(output) {
		pwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %s", err)
		}
		output = path.Join(pwd, output)
	}
	return output, nil
}
