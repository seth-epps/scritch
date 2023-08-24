package cli

type CLI struct {
	ScratchPath   string   `mapstructure:"scratch-path"`
	EditorCommand string   `mapstructure:"editor-command"`
	CustomSources []string `mapstructure:"custom-sources"`
}
