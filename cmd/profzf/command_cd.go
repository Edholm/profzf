package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCdCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cd",
		Short: "Print the cd compatible command",
		Run: func(cmd *cobra.Command, args []string) {
			const bin = "profzf"
			fmt.Printf(
				`cd $(%s list | fzf --delimiter $'\u200b' --nth 1 --tac --no-sort |%s get -il - | jq --raw-output '.path')`,
				bin, bin,
			)
		},
	}
	return cmd
}
