package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"edholm.dev/profzf/internal/server"
	"github.com/spf13/cobra"
)

func newServerCommand() *cobra.Command {
	opts := server.Config{
		ListenAddr:  "localhost:9910",
		ProjectDirs: nil,
		IgnoreDirs:  nil,
		ConfigDir:   "~/.config/profzf",
	}
	cmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"srv", "s"},
		Short:   "Start the server",
		Long:    "Start the server that keeps track of git projects and continuously updates them in the background",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			configDir, err := normalizePath(opts.ConfigDir)
			if err != nil {
				return fmt.Errorf("failed to normalize config dir: %w", err)
			}
			opts.ConfigDir = configDir

			if err := prepareConfigDir(opts.ConfigDir); err != nil {
				return fmt.Errorf("failed to prepare config dir: %w", err)
			}

			for i, dir := range opts.IgnoreDirs {
				p, err := normalizePath(dir)
				if err != nil {
					return fmt.Errorf("failed to normalize ignore dir: %w", err)
				}
				opts.IgnoreDirs[i] = p
			}
			for i, dir := range opts.ProjectDirs {
				p, err := normalizePath(dir)
				if err != nil {
					return fmt.Errorf("failed to normalize project dir: %w", err)
				}
				opts.ProjectDirs[i] = p
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := server.New(cmd.Context(), opts)
			if err != nil {
				return err
			}
			defer func() {
				if err := srv.Close(); err != nil {
					log.Printf("server.Close: %v", err)
				}
			}()
			return srv.Run(cmd.Context())
		},
	}
	f := cmd.Flags()
	f.StringSliceVarP(&opts.IgnoreDirs, "ignore-dir", "i", opts.IgnoreDirs, "Directory to ignore")
	f.StringSliceVarP(&opts.ProjectDirs, "project-dir", "p", opts.ProjectDirs, "Directory to search for git projects")
	f.StringVarP(&opts.ConfigDir, "config-dir", "c", opts.ConfigDir, "Directory to store config files")
	f.StringVarP(&opts.ListenAddr, "listen", "l", opts.ListenAddr, "Address to listen on")
	return cmd
}

func normalizePath(p string) (string, error) {
	if strings.ContainsRune(p, '~') {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", err)
		}
		p = strings.Replace(p, "~", home, 1)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs: %w", err)
	}
	return path.Clean(abs), nil
}

func prepareConfigDir(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}
	return nil
}
