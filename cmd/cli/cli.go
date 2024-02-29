package cli

import (
	"errors"
	"os/exec"

	"github.com/seth-epps/scritch/scratch/util"
)

type CLI struct {
	ScratchPath   string   `mapstructure:"scratch-path"`
	EditorCommand string   `mapstructure:"editor-command"`
	CustomSources []string `mapstructure:"custom-sources"`
	OutputFormat  string   `mapstructure:"output-format"`
}

// ResolveScratchPath attempts to do `~` shortcut replacement. If there's any errors
// it will return the original path with the corresponding error
func (cli CLI) ResolveScratchPath() (string, error) {
	return util.ReplaceHomeShortcut(cli.ScratchPath)
}

func (cli CLI) OpenEditor(path string) error {
	if cli.EditorCommand != "" {
		cmd := exec.Command(cli.EditorCommand, path)
		return cmd.Run()
	}
	return nil
}

// ResolveCustomSourcePaths attemps to replace any `~` shortcuts on the provided
// custom source paths. If there's any errors, they are skipped.
func (cli CLI) ResolveCustomSourcePaths() (resolvedPaths []string, err error) {
	for _, customPath := range cli.CustomSources {
		if resolved, replaceErr := util.ReplaceHomeShortcut(customPath); replaceErr != nil {
			errors.Join(err, replaceErr)
		} else {
			resolvedPaths = append(resolvedPaths, resolved)
		}
	}

	return resolvedPaths, err
}
