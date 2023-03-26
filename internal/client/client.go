package client

import (
	"context"
	"fmt"

	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(ctx context.Context, addr string) (pb.ProjectsServiceClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", addr, err)
	}
	return pb.NewProjectsServiceClient(conn), nil
}
