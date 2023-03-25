package server

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"edholm.dev/profzf/internal/server/db"
	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
)

type Config struct {
	ListenAddr  string   `json:"listenAddr"`
	ProjectDirs []string `json:"projectDirs"`
	IgnoreDirs  []string `json:"ignoreDirs"`
	ConfigDir   string   `json:"configDir"`
}

type Server struct {
	cfg       Config
	db        db.Querier
	dbCleanup db.CleanupFunc

	pb.UnimplementedProjectsServiceServer
}

func New(ctx context.Context, cfg Config) (*Server, error) {
	dbFile := path.Join(cfg.ConfigDir, "profzf.sqlite")
	db, cleanup, err := db.NewQuerier(ctx, dbFile, db.ModeReadWriteCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to setup db: %w", err)
	}
	return &Server{
		cfg:       cfg,
		db:        db,
		dbCleanup: cleanup,
	}, nil
}

// Run the server and blocks until the context is canceled.
func (s *Server) Run(ctx context.Context) error {
	log.Printf("up and running {config=%+v}", s.cfg)
	defer log.Printf("goodbye!")
	s.ScanProjects(ctx)
	s.TrimProjects(ctx)
	ticker := time.NewTicker(33 * time.Second)
	grpcErrCh := s.serveGRPC(ctx, s.cfg.ListenAddr)
	for {
		select {
		case <-ticker.C:
			s.ScanProjects(ctx)
			s.TrimProjects(ctx)
		case err := <-grpcErrCh:
			return fmt.Errorf("grpc server error: %w", err)
		case <-ctx.Done():
			return <-grpcErrCh
		}
	}
}

// Close closes the server and cleans up any resources (e.g. database connections).
func (s *Server) Close() error {
	return s.dbCleanup()
}
