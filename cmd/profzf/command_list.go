package main

import (
	"fmt"
	"os"

	"edholm.dev/profzf/internal/client"
	"edholm.dev/profzf/internal/fzf"
	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	serverAddr := defaultServerAddr
	format := "fzf"
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List git projects",
		Long:    "List all known git projects. Suitable for piping into fzf",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(cmd.Context(), serverAddr)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			resp, err := c.ListProjects(cmd.Context(), &pb.ListProjectsRequest{})
			if err != nil {
				return fmt.Errorf("failed to list projects: %w", err)
			}
			// buffer := bytes.NewBuffer(nil)
			fzf.TabPrint(os.Stdout, resp.GetProjects())
			// fmt.Printf("%s", buffer)
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&serverAddr, "addr", "a", serverAddr, "Address to the server, e.g. localhost:9010")
	f.StringVarP(&format, "format", "f", format, "Output format. Currently only 'fzf' is supported")
	return cmd
}
