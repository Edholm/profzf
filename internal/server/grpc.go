package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"

	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func (s *Server) serveGRPC(ctx context.Context, addr string) chan error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			errCh <- fmt.Errorf("failed to listen on %q: %w", addr, err)
			return
		}
		grpcServer := grpc.NewServer()
		pb.RegisterProjectsServiceServer(grpcServer, s)
		reflection.Register(grpcServer)

		go func() {
			<-ctx.Done()
			grpcServer.GracefulStop()
		}()
		log.Printf("serving grpc on %q", addr)
		defer log.Printf("grpc server going down")
		if err := grpcServer.Serve(listener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return
			}
			errCh <- fmt.Errorf("failed to serve grpc: %w", err)
			return
		}
	}()
	return errCh
}

func (s *Server) ListProjects(ctx context.Context, _ *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error) {
	repos, err := s.db.ListRepositories(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list projects: %v", err)
	}
	projects := make([]*pb.Project, len(repos))
	for i, repo := range repos {
		projects[i] = repo.Proto()
	}
	return &pb.ListProjectsResponse{
		Projects: projects,
	}, nil
}

func (s *Server) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.Project, error) {
	repo, err := s.db.GetByName(ctx, req.GetName())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no such project %q", req.GetName())
		}
	}
	if req.GetIncreaseUsageCount() {
		if err := s.db.IncRepoUsageCount(ctx, repo.Path); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to increase usage count: %v", err)
		}
	}
	repo.UsageCount++
	return repo.Proto(), nil
}
