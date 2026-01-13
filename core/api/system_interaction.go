package api

import (
	"context"
	"fmt"

	c "github.com/JuanBiancuzzo/own_wiki/core/systems/configuration"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

type SystemInteractionClient struct {
	conn   *grpc.ClientConn
	system pb.SystemInteractionClient
}

func NewSystemInteractionClient(config c.SystemInteractionConfig) (*SystemInteractionClient, error) {
	transportCredentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	direction := fmt.Sprintf("%s:%d", config.Ip, config.Port)

	conn, err := grpc.NewClient(direction, transportCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for SystemInteraction at %s, with error: %v", direction, err)
	}

	return &SystemInteractionClient{
		conn:   conn,
		system: pb.NewSystemInteractionClient(conn),
	}, nil
}

// TODO: Change "in *pb.QueryRequest" and "*pb.QueryResponse", to be the simplest representation, and this function
// should converted to request and response, respectively
func (sc *SystemInteractionClient) Query(ctx context.Context, in *pb.QueryRequest) (*pb.QueryResponse, error) {
	return sc.system.Query(ctx, in)
}

// TODO: Change "in *pb.SendEventRequest", to be a core/events/Event, and this function should change it to the request
func (sc *SystemInteractionClient) SendEvent(ctx context.Context, in *pb.SendEventRequest) error {
	_, err := sc.system.SendEvent(ctx, in)
	return err
}

func (sc *SystemInteractionClient) Close() {
	sc.conn.Close()
}
