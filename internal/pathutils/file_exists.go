package pathutils

import (
	"errors"
	"fmt"
	"os"
)

func AssertFileDoesntExist(path string) error {
	if targetFileInfo, err := os.Lstat(path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to get file information on path[%s]: %s", path, err)
		}
	} else {
		return fmt.Errorf("file[%s] at path[%s] already exists", targetFileInfo.Name(), path)
	}
	return nil
}

func EnsureDirectoryExists(path string) error {
	if lstatInfo, err := os.Lstat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create path[%s]: %s", path, err)
			}
		} else {
			return fmt.Errorf("failed to get info about path[%s]: %s", path, err)
		}
	} else if !lstatInfo.IsDir() {
		return fmt.Errorf("failed to get a directory at path[%s]", path)
	}
	return nil
}
