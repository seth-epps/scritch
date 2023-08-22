/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"log"

	"github.com/seth-epps/scritch/scratch"
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
	if language != "" {
		scratchLocation, err := scratch.GenerateScratch(language, options.variant)
		if err != nil {
			log.Fatalf("Failed to generate `%v (%v)` scratch: %v", language, options.variant, err)
			return nil
		}
		log.Printf("Created scratch at %v", scratchLocation)
		return nil
	}

	//TODO actually do something with this...
	if options.source == "" {
		return errors.New("Must provide custom template source path if language not specified.")
	}

	return nil
}
