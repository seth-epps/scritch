/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/seth-epps/scritch/cmd/cli"
	list_cmd "github.com/seth-epps/scritch/cmd/list"
	scratch_cmd "github.com/seth-epps/scritch/cmd/scratch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "scritch",
		Short: "Generate scratch pads for your favorite programming languages.",
		Long:  "Generate scratch pads for your favorite programming languages.",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var cli cli.CLI
	cobra.OnInitialize(func() {
		initializeConfig(&cli)
	})

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ~/.scritch/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&cli.ScratchPath, "scratch-path", "~/.scritch/scratch", "path to generate scratches at")
	rootCmd.PersistentFlags().StringVar(&cli.EditorCommand, "editor-command", "", "command to open an editor after scratch is generated")
	rootCmd.PersistentFlags().StringArrayVar(&cli.CustomSources, "custom-sources", []string{"~/.scritch/templates"}, "list of paths to search for custom source tempaltes")
	rootCmd.PersistentFlags().StringVarP(&cli.OutputFormat, "output-format", "o", "", "output format of result")

	viper.BindPFlag("scratch-path", rootCmd.PersistentFlags().Lookup("scratch-path"))
	viper.BindPFlag("editor-command", rootCmd.PersistentFlags().Lookup("editor-command"))
	viper.BindPFlag("custom-sources", rootCmd.PersistentFlags().Lookup("custom-sources"))

	rootCmd.AddCommand(scratch_cmd.NewScratchCommand(&cli))
	rootCmd.AddCommand(list_cmd.NewListCommand(&cli))
}

func initializeConfig(cli *cli.CLI) {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".scritch"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Could not read configuration; using defaults.")
		return
	}

	if err := viper.Unmarshal(cli); err != nil {
		cobra.CheckErr(fmt.Errorf("Couldn't construct configuration: %w", err))
	}

}
