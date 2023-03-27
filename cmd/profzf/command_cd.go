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
				//nolint:lll
				`cddir=$(%s list | fzf --delimiter $'\u200b' --nth 1 --tac --no-sort --preview 'git -C $(%s get -l -i=false {} | jq --raw-output .path) log -10' --preview-label="git log" |%s get -il - | jq --raw-output '.path') && cd "$cddir"`,
				bin, bin, bin,
			)
		},
	}
	return cmd
}
