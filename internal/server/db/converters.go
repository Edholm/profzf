package db

import (
	pb "edholm.dev/profzf/proto/gen/edholm/profzf/projects/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r Repository) Proto() *pb.Project {
	return &pb.Project{
		Name: r.Name,
		Path: r.Path,
		GitStatus: &pb.GitStatus{
			Branch:     r.GitBranch,
			Dirty:      r.GitDirty,
			Action:     r.GitAction,
			LeftCount:  int32(r.GitCountLeft),
			RightCount: int32(r.GitCountRight),
		},
		UsageCount:  int32(r.UsageCount),
		LastUpdated: timestamppb.New(r.UpdateTime),
	}
}
