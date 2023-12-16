package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ExpandTilde(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			return "", fmt.Errorf("Could not expand tilde in path:\n%s", err)
		}

        fullPath := filepath.Join(homeDir, path[1:])

		return fullPath, nil
	} else {
		return path, nil
	}
}
