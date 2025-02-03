package main

import (
	"fmt"
	"os"

	"edholm.dev/profzf/internal/client"
	"edholm.dev/profzf/internal/fzf"
	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
	"github.com/spf13/cobra"
)

func newListCommand(common commonOpts) *cobra.Command {
	format := "fzf"
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List git projects",
		Long:    "List all known git projects. Suitable for piping into fzf",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := client.New(common.SocketPath)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			resp, err := c.ListProjects(cmd.Context(), &pb.ListProjectsRequest{})
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}
			fzf.TabPrint(os.Stdout, resp.GetProjects())
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&format, "format", "f", format, "Output format. Currently only 'fzf' is supported")
	return cmd
}
