package client

import (
	"context"
	"fmt"
	"net"

	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New(socketPath string) (pb.ProjectsServiceClient, error) {
	conn, err := grpc.NewClient(
		socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(
			func(_ context.Context, addr string) (net.Conn, error) {
				return net.Dial("unix", addr)
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", socketPath, err)
	}
	return pb.NewProjectsServiceClient(conn), nil
}
