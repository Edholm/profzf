package server

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitInfo struct {
	Branch string
	Dirty  bool
}

// GitStatus returns the git status of the given directory.
func GitStatus(d string) (GitInfo, error) {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	cmd.Dir = d
	ref, err := cmd.Output()
	if err != nil {
		return GitInfo{}, fmt.Errorf("git symbolic-ref HEAD: %w", err)
	}
	branch := strings.TrimSpace(string(ref))

	cmd = exec.Command("git", "--no-optional-locks", "status", "--porcelain", "--untracked-files=normal")
	cmd.Dir = d

	dirty := false
	if out, err := cmd.Output(); err == nil {
		dirty = len(out) > 0
	}
	return GitInfo{
		Branch: branch,
		Dirty:  dirty,
	}, nil
}
