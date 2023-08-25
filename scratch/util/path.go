package util

import (
	"os"
	"path/filepath"
	"strings"
)

func ReplaceHomeShortcut(path string) (string, error) {
	switch {
	case path == "~":
		if home, err := os.UserHomeDir(); err != nil {
			return path, err
		} else {
			return home, err
		}
	case strings.HasPrefix(path, "~/"):
		if home, err := os.UserHomeDir(); err != nil {
			return path, err
		} else {
			return filepath.Join(home, path[2:]), err
		}
	}
	return path, nil
}
