/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/seth-epps/scritch/cmd/cli"
	"github.com/seth-epps/scritch/scratch"
	"github.com/seth-epps/scritch/scratch/templates"
	"github.com/spf13/cobra"
)

type scratchResult struct {
	Path string `json:"path"`
}

// NewScratchCommand creates a new `scritch scratch` command
func NewScratchCommand(cli *cli.CLI) *cobra.Command {

	var scratchCmd = &cobra.Command{
		Use:   "scratch [source]",
		Short: "Create a scratch for specified source.",
		Long:  `Create a scratch for natively supported source langauge or specify your own source template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("must provide a source")
			}
			return runScratch(cli, args[0])
		},
	}

	// remove additional "[flags]" in usage string
	scratchCmd.DisableFlagsInUseLine = true

	return scratchCmd
}

func runScratch(cli *cli.CLI, source string) error {
	templateProvider, err := getTemplateProvider(cli, source)
	if err != nil {
		log.Fatalf("Could not find template provider: %v", err)
		return nil
	}

	scratchPath, err := cli.ResolveScratchPath()
	if err != nil {
		log.Printf("WARNING: Could not resolve the scratch path, attempting to use %v: %v", scratchPath, err)
	}
	scratch := scratch.NewScratch(scratchPath)
	scratchLocation, err := scratch.GenerateScratch(templateProvider)

	if err != nil {
		log.Fatalf("Failed to generate scratch: %v", err)
		return nil
	}

	printScratch(cli, scratchResult{scratchLocation})
	if err = cli.OpenEditor(scratchLocation); err != nil {
		fmt.Printf("Couldn't open editor: %v\n", err)
	}

	return nil
}

func printScratch(cli *cli.CLI, res scratchResult) error {
	switch cli.OutputFormat {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(res)
	default:
		fmt.Printf("Created scratch at %v\n", res.Path)
		return nil
	}
}

func getTemplateProvider(cli *cli.CLI, source string) (templates.TemplateProvider, error) {
	if slices.Contains(templates.SupportedTemplates, source) {
		return templates.NewEmbeddedTemplateProvider(source), nil
	}

	searchLocations, err := cli.ResolveCustomSourcePaths()
	if err != nil {
		fmt.Printf("WARNING: Some provided paths could not be resolved: %v\n", err)
	}
	return templates.NewFilesystemTemplateProvider(source, searchLocations), nil

}
