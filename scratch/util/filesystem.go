package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func ListDirs(startingPath string, targetDepth int, includePath bool, before *time.Time) ([]string, error) {
	var targets []string
	err := filepath.WalkDir(startingPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		depthFromRoot := strings.Count(path, string(os.PathSeparator)) - strings.Count(startingPath, string(os.PathSeparator))
		if depthFromRoot < targetDepth {
			// keep going and do nothing till we're looking
			// in the right place
			return nil
		}

		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("could not read directory info: %w", err)
			}

			switch {
			case depthFromRoot > targetDepth:
				// we've gone too far and need to bail out of the walk
				return fs.SkipDir
			case before == nil || (before != nil && info.ModTime().Before(*before)):
				// ahhhh this is too much complexity
				if includePath {
					targets = append(targets, path)
				} else {
					targets = append(targets, d.Name())
				}
			}
		}
		return nil
	})

	return targets, err
}
