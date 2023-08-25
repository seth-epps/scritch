/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/seth-epps/scritch/cmd/cli"
	"github.com/seth-epps/scritch/scratch"
	"github.com/seth-epps/scritch/scratch/templates"
	"github.com/spf13/cobra"
)

type scratchOptions struct {
	variant string
	source  string
}

// NewScratchCommand creates a new `scritch scratch` command
func NewScratchCommand(cli *cli.CLI) *cobra.Command {
	opts := scratchOptions{}

	var scratchCmd = &cobra.Command{
		Use:   "scratch [language]",
		Short: "Create a scratch for specified supported langauge.",
		Long: `Create a scratch for specified langauge or specify your own source 
template.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScratch(cli, getLanguageFromArgs(args), opts)
		},
	}

	// remove additional "[flags]" in usage string
	scratchCmd.DisableFlagsInUseLine = true

	scratchCmd.Flags().StringVarP(&opts.variant, "variant", "v", "default", "Language template variant.")
	scratchCmd.Flags().StringVarP(&opts.source, "source", "s", "", "Custom source template directory to generate scratch workspace for if language is not specified.")

	return scratchCmd
}

func getLanguageFromArgs(args []string) string {
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func runScratch(cli *cli.CLI, language string, options scratchOptions) error {
	templateProvider, err := getTemplateProvider(cli, language, options)
	if err != nil {
		log.Fatalf("Could not find template provider: %v", err)
		return nil
	}

	scratchPath, err := cli.ResolveScratchPath()
	if err != nil {
		fmt.Printf("WARNING: Could not resolve the scratch path, attempting to use %v: %v\n", scratchPath, err)
	}
	scratch := scratch.NewScratch(scratchPath)
	scratchLocation, err := scratch.GenerateScratch(templateProvider)

	if err != nil {
		log.Fatalf("Failed to generate scratch: %v", err)
		return nil
	}
	log.Printf("Created scratch at %v\n", scratchLocation)
	if err = cli.OpenEditor(scratchLocation); err != nil {
		fmt.Printf("Couldn't open editor: %v\n", err)
	}

	return nil
}

func getTemplateProvider(cli *cli.CLI, language string, options scratchOptions) (templates.TemplateProvider, error) {
	if language != "" {
		return templates.NewEmbeddedTemplateProvider(language, options.variant), nil
	}

	if options.source != "" {
		searchLocations, err := cli.ResolveCustomSourcePaths()
		if err != nil {
			fmt.Printf("WARNING: Some provided paths could not be resolved: %v\n", err)
		}
		return templates.NewFilesystemTemplateProvider(options.source, searchLocations), nil
	}

	return nil, errors.New("Must provide custom template source path if language not specified.")
}
