package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode"

	"github.com/seth-epps/scritch/cmd/cli"
	"github.com/seth-epps/scritch/scratch/util"
	"github.com/spf13/cobra"
)

type cleanOpts struct {
	before string
	dryRun bool
	force  bool
	source string
}

// NewListCommand creates a new `scritch list` command
func NewCleanCommand(cli *cli.CLI) *cobra.Command {
	var opts cleanOpts
	var scratchCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean up generated scratch files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClean(cli, opts)
		},
	}

	// remove additional "[flags]" in usage string
	scratchCmd.DisableFlagsInUseLine = true

	scratchCmd.Flags().StringVarP(&opts.before, "before", "b", "", "cutoff date-modified timestamp (unix or RFC3339) for deleting scratches")
	scratchCmd.Flags().StringVarP(&opts.source, "source", "s", "", "only clean scratches for the provided source name")
	scratchCmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "only output the scratches targetted without deleting anything")
	scratchCmd.Flags().BoolVar(&opts.force, "force", false, "force cleanup without prompting for confirmation")

	return scratchCmd
}

func runClean(cli *cli.CLI, opts cleanOpts) error {
	before, err := parseTime(opts.before)
	if err != nil {
		return err
	}

	scratchPath, err := cli.ResolveScratchPath()
	if err != nil {
		return fmt.Errorf("could not resolve scratch directory: %w", err)
	}

	targets, err := listExisting(scratchPath, opts.source, before)
	if err != nil {
		return err
	}

	if !opts.force {
		c, err := getConfirmation()
		if !c {
			return err
		}
	}

	var deleteErrors error
	for _, target := range targets {
		deleteErrors = errors.Join(deleteErrors, delete(target, opts.dryRun))
	}

	return deleteErrors
}

func parseTime(t string) (*time.Time, error) {
	if t == "" {
		return nil, nil
	}

	var _time time.Time
	var err error

	// attempt to parse as RFC3339
	_time, rfcErr := time.Parse(time.RFC3339, t)
	if rfcErr != nil {
		// try again as unix
		tmp, parseIntErr := strconv.ParseInt(t, 10, 64)
		if parseIntErr != nil {
			err = fmt.Errorf("failed to parse `%s` as RFC3339 or unix time", t)
		}
		_time = time.Unix(tmp, 0)
	}

	return &_time, err
}

func listExisting(scratchPath, source string, before *time.Time) ([]string, error) {
	// only search the specified name path
	if source != "" {
		return util.ListDirs(filepath.Join(scratchPath, source), 1, true, before)
	}
	return util.ListDirs(scratchPath, 2, true, before)
}

func getConfirmation() (bool, error) {
	hocho := '\U0001F52A'
	fmt.Printf("Are you sure you want to delete all your scratch files? [y/%c/n]\n", hocho)
	reader := bufio.NewReader(os.Stdin)
	c, _, err := reader.ReadRune()
	if err != nil {
		return false, fmt.Errorf("could not read input: %w", err)
	}
	if unicode.ToLower(c) == 'y' {
		return true, nil
	}
	return true, nil

}

func delete(target string, dryRun bool) error {
	log.Printf("Deleting %s", target)
	if dryRun {
		// don't actually do anything for dry-run
		return nil
	}

	return os.RemoveAll(target)
}
