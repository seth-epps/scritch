/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"text/tabwriter"

	"github.com/seth-epps/scritch/cmd/cli"
	"github.com/seth-epps/scritch/scratch/templates"
	"github.com/seth-epps/scritch/scratch/util"
	"github.com/spf13/cobra"
)

type listItem struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// NewListCommand creates a new `scritch list` command
func NewListCommand(cli *cli.CLI) *cobra.Command {
	var scratchCmd = &cobra.Command{
		Use:   "list",
		Short: "List available scratch sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cli)
		},
	}

	// remove additional "[flags]" in usage string
	scratchCmd.DisableFlagsInUseLine = true

	return scratchCmd
}

func runList(cli *cli.CLI) error {
	var res []listItem
	for _, source := range templates.SupportedTemplates {
		res = append(res, listItem{source, "NATIVE"})
	}

	searchLocations, err := cli.ResolveCustomSourcePaths()
	if err != nil {
		fmt.Printf("WARNING: Some provided paths could not be resolved: %v\n", err)
	}

	fsSources := listSources(searchLocations)
	for sourceLocation, sources := range fsSources {
		for _, source := range sources {
			// native templates take precedence over custom templates
			if !slices.Contains(templates.SupportedTemplates, source) {
				res = append(res, listItem{source, sourceLocation})
			}
		}
	}

	return printList(cli, res)
}

func printList(cli *cli.CLI, l []listItem) error {
	switch cli.OutputFormat {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(l)
	default:
		tw := tabwriter.NewWriter(os.Stdout, 1, 3, 3, ' ', 0)
		fmt.Fprintln(tw, "Name\tSource")
		for _, item := range l {
			fmt.Fprintf(tw, "%s\t%s\n", item.Name, item.Source)
		}
		tw.Flush()
		return nil
	}
}

func listSources(locations []string) map[string][]string {
	sourceMap := make(map[string][]string)

	for _, searchPath := range locations {
		contents, err := util.ListDirs(searchPath, 1, false, nil)
		if err != nil {
			// intentionally skip directories when we encounter errors
			// todo (seth) maybe log?
			continue
		}

		sourceMap[searchPath] = contents
	}

	return sourceMap
}
