/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/seth-epps/scritch/scratch"
	"github.com/seth-epps/scritch/scratch/templates"
	"github.com/spf13/cobra"
)

type scratchOptions struct {
	variant string
	source  string
}

// NewScratchCommand creates a new `scritch scratch` command
func NewScratchCommand() *cobra.Command {
	opts := scratchOptions{}

	var scratchCmd = &cobra.Command{
		Use:   "scratch [language]",
		Short: "Create a scratch for specified supported langauge.",
		Long: `Create a scratch for specified langauge or (TODO) specify your own source 
templates.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runScratch(getLanguageFromArgs(args), opts)
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

func runScratch(language string, options scratchOptions) error {
	templateProvider, err := getTemplateProvider(language, options)
	if err != nil {
		log.Fatalf("Could not find template provider: %v", err)
		return nil
	}

	scratcher := scratch.NewDefaultScratcher()
	scratchLocation, err := scratcher.GenerateScratch(templateProvider)
	if err != nil {
		log.Fatalf("Failed to generate scratch: %v", err)
		return nil
	}
	log.Printf("Created scratch at %v", scratchLocation)
	openEditor(scratchLocation)
	return nil
}

func getTemplateProvider(language string, options scratchOptions) (templates.TemplateProvider, error) {
	if language != "" {
		return templates.NewEmbeddedTemplateProvider(language, options.variant), nil
	}
	if options.source != "" {
		return templates.NewFilesystemTemplateProvider(options.source), nil
	}

	return nil, errors.New("Must provide custom template source path if language not specified.")
}

func openEditor(path string) {
	cmd := exec.Command("code", path)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Couldn't open editor: %v", err)
	}
}
