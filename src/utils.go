package tags

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func expandTilde(path string) string {
    if strings.HasPrefix(path, "~") {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
            return path
        }

        return filepath.Join(homeDir, path[1:])
    } else {
        return path
    }
}
