package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CLI struct {
	ScratchPath   string   `mapstructure:"scratch-path"`
	EditorCommand string   `mapstructure:"editor-command"`
	CustomSources []string `mapstructure:"custom-sources"`
}

// ResolveScratchPath attempts to do `~` shortcut replacement. If there's any errors
// it will return the original path with the corresponding error
func (cli CLI) ResolveScratchPath() (string, error) {
	switch {
	case cli.ScratchPath == "~":
		if home, err := os.UserHomeDir(); err != nil {
			return cli.ScratchPath, err
		} else {
			return home, err
		}
	case strings.HasPrefix(cli.ScratchPath, "~/"):
		if home, err := os.UserHomeDir(); err != nil {
			return cli.ScratchPath, err
		} else {
			return filepath.Join(home, cli.ScratchPath[2:]), err
		}
	}
	return cli.ScratchPath, nil
}

func (cli CLI) OpenEditor(path string) error {
	if cli.EditorCommand != "" {
		cmd := exec.Command(cli.EditorCommand, path)
		return cmd.Run()
	}
	return nil
}
