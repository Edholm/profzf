package main

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

type commonOpts struct {
	SocketPath string
}

func newRootCommand() *cobra.Command {
	common := commonOpts{
		SocketPath: "~/.config/profzf/profzf.sock",
	}
	if path, err := normalizePath(common.SocketPath); err == nil {
		common.SocketPath = path
	}

	cmd := &cobra.Command{
		Use:           "profzf",
		Short:         "List git projects for use with fzf and cd",
		Long:          "Allows you to quickly list git projects and cd into them using fzf",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
		Version: version(),
	}

	cmd.PersistentFlags().StringVarP(&common.SocketPath, "socket", "s", common.SocketPath, "Path to the socket file")

	cmd.AddCommand(newServerCommand(common))
	cmd.AddCommand(newListCommand(common))
	cmd.AddCommand(newGetCommand(common))
	cmd.AddCommand(newCdCommand())
	return cmd
}

//nolint:gochecknoglobals // This is set at build time, in sage.
var LDDVersion string

func version() string {
	if LDDVersion != "" {
		return LDDVersion
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range bi.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return "dev"
}
