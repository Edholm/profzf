package server

import (
	"context"
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
		projects[i] = &pb.Project{
			Path: repo.Path,
			Name: repo.Name,
			GitStatus: &pb.GitStatus{
				Branch: repo.GitBranch,
				Dirty:  repo.GitDirty,
			},
		}
	}
	return &pb.ListProjectsResponse{
		Projects: projects,
	}, nil
}
