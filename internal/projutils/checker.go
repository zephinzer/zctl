package projutils

import (
	"fmt"
	"os"
)

const (
	TypeGo   = "golang"
	TypeJava = "java"
	TypeJs   = "javascript"
	TypePy   = "python"
)

func IsJavascriptNpmProjecct(directory string) (bool, error) {
	return checkForFiles(directory, []string{
		"package-lock.json",
	})
}

func IsJavascriptYarnProjecct(directory string) (bool, error) {
	return checkForFiles(directory, []string{
		"yarn.lock",
	})
}

func IsGolangProject(directory string) (bool, error) {
	return checkForFiles(directory, []string{
		"go.mod",
		"go.sum",
	})
}

func IsPythonProject(directory string) (bool, error) {
	return checkForFiles(directory, []string{
		"requirements.txt",
	})
}

func checkForFiles(directory string, signatureFiles []string) (bool, error) {
	requiredFiles := map[string]bool{}
	for _, signatureFile := range signatureFiles {
		requiredFiles[signatureFile] = false
	}
	listings, err := os.ReadDir(directory)
	if err != nil {
		return false, fmt.Errorf("failed to list contents of directory[%s]: %s", directory, err)
	}
	for _, listing := range listings {
		if _, exists := requiredFiles[listing.Name()]; exists {
			requiredFiles[listing.Name()] = true
		}
	}
	for _, exists := range requiredFiles {
		if !exists {
			return false, nil
		}
	}
	return true, nil
}
