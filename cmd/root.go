/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	cmd "github.com/seth-epps/scritch/cmd/scratch"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scritch",
	Short: "Generate scratch pads for your favorite programming languages.",
	Long:  "Generate scratch pads for your favorite programming languages.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cmd.NewScratchCommand())
}
