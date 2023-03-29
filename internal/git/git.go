package git

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Action string

const (
	ActionRebase            Action = "rebase"
	ActionAM                Action = "am"
	ActionAMRebase          Action = "am/rebase"
	ActionRebaseInteractive Action = "rebase-i"
	ActionRebaseMerge       Action = "rebase-m"
	ActionMerge             Action = "merge"
	ActionBisect            Action = "bisect"
	ActionCherryPick        Action = "cherry"
	ActionCherryPickSeq     Action = "cherry-seq"
	ActionCherryOrRevert    Action = "cherry-or-revert"
	ActionNone              Action = ""
)

type Info struct {
	Branch     string
	Action     Action
	Dirty      bool
	LeftCount  int // number of commits on head vs upstream
	RightCount int // number of commits on upstream vs head
}

// GetInfo returns information about the git repository at the given path.
func GetInfo(d string) (Info, error) {
	cmd := exec.Command("git", "--no-optional-locks", "status", "--porcelain", "--untracked-files=normal")
	cmd.Dir = d
	dirty := false
	if out, err := cmd.Output(); err == nil {
		dirty = len(out) > 0
	}
	left, right, err := countLeftRight(d)
	if err != nil {
		left = -1
		right = -1
	}
	return Info{
		Branch:     determineBranch(d),
		Action:     determineAction(d),
		Dirty:      dirty,
		LeftCount:  left,
		RightCount: right,
	}, nil
}

func determineBranch(repo string) string {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		cmd = exec.Command("git", "describe", "--all", "--exact-match", "HEAD")
		cmd.Dir = repo
		out, err = cmd.Output()
		if err != nil {
			return "n/a"
		}
		return strings.TrimSpace(string(out))
	}
	return strings.TrimSpace(string(out))
}

func determineAction(repo string) Action {
	gd := path.Join(repo, ".git")
	for _, s := range []string{
		path.Join(gd, "rebase-apply"),
		path.Join(gd, "rebase"),
	} {
		if isDir(s) {
			switch {
			case isFile(path.Join(s, "rebasing")):
				return ActionRebase
			case isFile(path.Join(s, "applying")):
				return ActionAM
			default:
				return ActionAMRebase
			}
		}
	}
	if isFile(path.Join(gd, "rebase-merge", "interactive")) {
		return ActionRebaseInteractive
	}
	if isDir(path.Join(gd, "rebase-merge")) {
		return ActionRebaseMerge
	}
	if isFile(path.Join(gd, "MERGE_HEAD")) {
		return ActionMerge
	}
	if isFile(path.Join(gd, "BISECT_LOG")) {
		return ActionBisect
	}
	if isFile(path.Join(gd, "CHERRY_PICK_HEAD")) {
		if isFile(path.Join(gd, "sequencer")) {
			return ActionCherryPickSeq
		}
		return ActionCherryPick
	}
	if isDir(path.Join(gd, "sequencer")) {
		return ActionCherryOrRevert
	}
	return ActionNone
}

func countLeftRight(repo string) (int, int, error) {
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("rev-list failed: %w", err)
	}
	l, r, found := strings.Cut(strings.TrimSpace(string(out)), "\t")
	if !found {
		return 0, 0, fmt.Errorf("invalid rev-list output: %s", out)
	}
	left, err := strconv.Atoi(l)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid left count: %w", err)
	}
	right, err := strconv.Atoi(r)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid right count: %w", err)
	}
	return left, right, nil
}

func isDir(path string) bool {
	if s, err := os.Stat(path); err == nil {
		return s.IsDir()
	}
	return false
}

func isFile(path string) bool {
	if s, err := os.Stat(path); err == nil {
		return !s.IsDir()
	}
	return false
}
