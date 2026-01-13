package api

import (
	"context"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
)

type SystemInteraction struct {
	pb.UnimplementedSystemInteractionServer
}

func (*SystemInteraction) Query(context.Context, *pb.QueryRequest) (*pb.QueryResponse, error) {

	return nil, nil
}

func (*SystemInteraction) SendEvent(context.Context, *pb.SendEventRequest) (*pb.Empty, error) {

	return &pb.Empty{}, nil
}
