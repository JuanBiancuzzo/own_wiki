package api

import (
	"context"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
)

type SystemInteraction struct {
	pb.UnimplementedSystemInteractionServer
}

func (*SystemInteraction) Query(ctx context.Context, queryRequest *pb.QueryRequest) (*pb.QueryResponse, error) {
	// We would need the database conextion to responde to this request

	return nil, nil
}

func (*SystemInteraction) SendEvent(ctx context.Context, eventRequest *pb.SendEventRequest) (*pb.Empty, error) {
	// We would need the event queue to send the event

	return &pb.Empty{}, nil
}
