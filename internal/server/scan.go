package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"edholm.dev/profzf/internal/git"
	"edholm.dev/profzf/internal/server/db"
)

const workerCount = 10

func (s *Server) ScanProjects(ctx context.Context) {
	// Spawn workers
	repo := make(chan string)
	var wg sync.WaitGroup
	count := 0
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go s.upsertRepoWorker(ctx, repo, &wg)
	}
	for _, rootDir := range s.cfg.ProjectDirs {
		gitDirs, err := s.findGitRepos(rootDir)
		if err != nil {
			fmt.Printf("failed to find git repos in %q: %v", rootDir, err)
			continue
		}
		for _, gitDir := range gitDirs {
			repo <- gitDir
			count++
		}
	}
	close(repo)
	wg.Wait() // Wait for all workers to finish
}

func (s *Server) upsertRepoWorker(ctx context.Context, repo <-chan string, wg *sync.WaitGroup) {
	for gitDir := range repo {
		status, err := git.GetInfo(gitDir)
		if err != nil {
			log.Printf("failed to get git status for %q: %v", gitDir, err)
		}
		if err := s.db.UpsertRepository(ctx, db.UpsertRepositoryParams{
			Path:          gitDir,
			Name:          path.Base(gitDir),
			GitBranch:     status.Branch,
			GitDirty:      status.Dirty,
			GitAction:     string(status.Action),
			GitCountLeft:  int64(status.LeftCount),
			GitCountRight: int64(status.RightCount),
		}); err != nil {
			log.Printf("failed to insert repo %q: %v", gitDir, err)
		}
	}
	wg.Done()
}

// TrimProjects removes projects that no longer exist or should be skipped by the server.
func (s *Server) TrimProjects(ctx context.Context) {
	repos, err := s.db.ListRepositories(ctx)
	if err != nil {
		log.Printf("failed to trim projects: unable to list repos: %v", err)
		return
	}
	for _, repo := range repos {
		if !s.shouldKeep(repo.Path) {
			log.Printf("trimming %q", repo.Path)
			if err := s.db.DeleteRepository(ctx, repo.Path); err != nil {
				log.Printf("failed to trim projects: unable to delete repo %q: %v", repo.Path, err)
			}
		}
	}
}

func (s *Server) findGitRepos(rootDir string) ([]string, error) {
	gitDirs := make([]string, 0, 25)
	lastProject := ""

	abs, err := filepath.Abs(rootDir)
	if err != nil {
		return gitDirs, fmt.Errorf("failed to get absolute path for %q: %w", rootDir, err)
	}

	if err := filepath.WalkDir(abs, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if lastProject != "" && strings.HasPrefix(p, lastProject) || s.skipDir(p) {
			return filepath.SkipDir
		}
		if _, err := os.Lstat(path.Join(p, ".git")); err == nil {
			absP, err := filepath.Abs(p)
			if err == nil {
				p = absP
			}
			gitDirs = append(gitDirs, p)
			lastProject = p
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
		return gitDirs, fmt.Errorf("failed to walk %q: %w", rootDir, err)
	}
	return gitDirs, nil
}

func (s *Server) skipDir(p string) bool {
	gopath := os.Getenv("GOPATH")
	if gopath != "" &&
		strings.HasPrefix(p, path.Join(gopath, "bin")) ||
		strings.HasPrefix(p, path.Join(gopath, "pkg")) {
		return true
	}
	name := path.Base(p)
	for _, ignoreDir := range s.cfg.IgnoreDirs {
		abs, err := filepath.Abs(ignoreDir)
		if err != nil {
			continue
		}
		if p == abs {
			return true
		}
	}

	switch name {
	case "node_modules", "vendor", "dist", "build", "bin", "out", "target", ".idea", ".vscode", ".svn", ".gradle":
		return true
	default:
		return false
	}
}

// shouldKeep returns true if the given path should be kept in the database
// based on whether it exists and whether it should be ignored.
func (s *Server) shouldKeep(p string) bool {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return false
	}
	// TODO: not a part of the project dirs
	for _, ignoreDir := range s.cfg.IgnoreDirs {
		if strings.HasPrefix(p, ignoreDir) {
			return false
		}
	}
	return true
}
