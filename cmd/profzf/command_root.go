package main

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profzf",
		Short: "List git projects for use with fzf and cd",
		Long:  "Allows you to quickly list git projects and cd into them using fzf",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Version: version(),
	}

	cmd.AddCommand(newServerCommand())
	return cmd
}

func version() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range bi.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return "dev"
}
