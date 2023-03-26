package main

import (
	"fmt"
	"io"

	"edholm.dev/profzf/internal/client"
	"edholm.dev/profzf/internal/fzf"
	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func newGetCommand() *cobra.Command {
	serverAddr := defaultServerAddr
	increaseUsage := true
	listName := false
	cmd := &cobra.Command{
		Use:   "get [repo name]",
		Short: "Get the repository",
		Long:  "Get the repository based on the repo name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(cmd.Context(), serverAddr)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}
			name := args[0]
			if name == "-" {
				cmd.InOrStdin()
				input, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("failed to read input: %w", err)
				}
				name = string(input)
			}

			if listName {
				name = fzf.ExtractName(name)
			}
			repo, err := c.GetProject(cmd.Context(), &pb.GetProjectRequest{
				Name:               name,
				IncreaseUsageCount: increaseUsage,
			})
			if err != nil {
				return fmt.Errorf("failed to get project: %w", err)
			}
			fmt.Println(protojson.Format(repo))
			return nil
		},
	}
	f := cmd.Flags()
	f.StringVarP(&serverAddr, "addr", "a", serverAddr, "Address to the server, e.g. localhost:9010")
	f.BoolVarP(&increaseUsage, "increase-usage", "i", increaseUsage, "Increase the usage count for the project")
	f.BoolVarP(&listName, "list-name", "l", listName, "The name comes from the list command, i.e. contains git info etc")
	return cmd
}
