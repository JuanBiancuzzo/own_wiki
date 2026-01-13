package api

import (
	"context"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
)

type UserInteraction struct {
	pb.UnimplementedUserInteractionServer
}

func (*UserInteraction) LoadPlugin(context.Context, *pb.LoadPluginRequest) (*pb.Empty, error) {

	return &pb.Empty{}, nil
}

func (*UserInteraction) RegisterComponents(context.Context, *pb.Empty) (*pb.RegisterComponentsResponse, error) {

	return nil, nil
}

func (*UserInteraction) ImportFiles(pb.UserInteraction_ImportFilesServer) error {

	return nil
}

func (*UserInteraction) Render(pb.UserInteraction_RenderServer) error {

	return nil
}
