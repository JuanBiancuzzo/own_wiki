package api

import (
	"context"
	"fmt"
	"net"

	c "github.com/JuanBiancuzzo/own_wiki/core/systems/configuration"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userInteraction struct {
	pb.UnimplementedUserInteractionServer
}

func newUserInteraction() *userInteraction {
	return &userInteraction{}
}

func (*userInteraction) LoadPlugin(ctx context.Context, loadPluginRequest *pb.LoadPluginRequest) (*pb.LoadPluginResponse, error) {
	// We would need to find the executable (build as a plugin) and save the reference to that plugin.
	// With the plugin, we should call the necesary function to define the structures given by the user
	// This should be the components, and the views. Then the components should be return

	return nil, nil
}

func (*userInteraction) ImportFiles(importFileStream pb.UserInteraction_ImportFilesServer) error {
	// We could send via a channel the file paths needed, and via a gorouting be processing them as the
	// user defines. Finally, via another channel, send the component description to the stream

	return nil
}

func (*userInteraction) Render(renderStream pb.UserInteraction_RenderServer) error {
	// First, we would need to create the SystemInteractionClient to get the Query and SendEvent functions.
	// We would need the view register in the LoadPlugin function, and create the gorouting to continuously
	// be able to send events, and get the scene render

	return nil
}

func (*userInteraction) Close() {
	// Close the plugin
}

type UserInteractionServer struct {
	listener net.Listener
	server   *grpc.Server
}

func NewUserInteractionServer(config c.UserInteractionConfig) (*UserInteractionServer, error) {
	direction := fmt.Sprintf("%s:%d", config.Ip, config.Port)
	lis, err := net.Listen(config.Protocol, direction)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener server for UserInteraction at %s, with error: %v", direction, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserInteractionServer(grpcServer, newUserInteraction())

	return &UserInteractionServer{
		listener: lis,
		server:   grpcServer,
	}, nil
}

func (us *UserInteractionServer) Serve() error {
	if err := us.server.Serve(us.listener); err != nil {
		return fmt.Errorf("failed to serve UserInteraction, with error: %v", err)
	}
	return nil
}

type UserInteractionClient struct {
	conn *grpc.ClientConn
	user pb.UserInteractionClient
}

func NewUserInteractionClient(config c.UserInteractionConfig) (*UserInteractionClient, error) {
	transportCredentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	direction := fmt.Sprintf("%s:%d", config.Ip, config.Port)

	conn, err := grpc.NewClient(direction, transportCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for UserInteraction at %s, with error: %v", direction, err)
	}

	return &UserInteractionClient{
		conn: conn,
		user: pb.NewUserInteractionClient(conn),
	}, nil
}

// TODO: Change LoadPluginResponse to be the structures of components
func (uc *UserInteractionClient) LoadPlugin(ctx context.Context, pluginPath string) (*pb.LoadPluginResponse, error) {
	return uc.user.LoadPlugin(ctx, &pb.LoadPluginRequest{
		PluginPath: pluginPath,
	})
}

// TODO: Change the function to accept a channel of string (send filePaths), and a channel to get the information
// to load to the database
func (uc *UserInteractionClient) ImportFiles(ctx context.Context) (pb.UserInteraction_ImportFilesClient, error) {
	return uc.user.ImportFiles(ctx)
}

// TODO: Change the function to accept a channel for events (as in core/events/Event) to send, and a channel to get
// the scene representation
func (uc *UserInteractionClient) Render(ctx context.Context) (pb.UserInteraction_RenderClient, error) {
	return uc.user.Render(ctx)
}

func (uc *UserInteractionClient) Close() {
	uc.conn.Close()
}
