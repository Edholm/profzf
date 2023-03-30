package fzf

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
)

// zero-width space.
const zwsp = '\u200B'

func TabPrint(w io.Writer, repos []*pb.Project) {
	tw := tabwriter.NewWriter(w, 1, 0, 4, ' ', 0)
	for _, repo := range repos {
		gs := repo.GetGitStatus()
		dirty := ""
		if gs.GetDirty() {
			dirty = "*"
		}
		action := gs.GetAction()
		if action != "" {
			action = " (" + action + ")"
		}
		arrows := ""
		if repo.GitStatus.LeftCount > 0 {
			arrows = "⇡"
		}
		if repo.GitStatus.RightCount > 0 {
			arrows += "⇣"
		}
		if len(arrows) > 0 {
			arrows = " " + arrows + " "
		}
		println(arrows)
		// Format: <name>\u200B	<branch><*> <⇣⇡> (<action>)
		_, _ = fmt.Fprintf(tw, "%s%c\t%s%s%s%s\n", repo.GetName(), zwsp, gs.GetBranch(), dirty, arrows, action)
	}
	_ = tw.Flush()
}

// ExtractName extracts the name from the string that was formatted by TabPrint.
func ExtractName(s string) string {
	name, _, found := strings.Cut(s, string(zwsp))
	if !found {
		return s
	}
	return name
}
