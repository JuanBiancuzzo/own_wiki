package api

import (
	"context"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
)

type UserInteraction struct {
	pb.UnimplementedUserInteractionServer
}

func (*UserInteraction) LoadPlugin(ctx context.Context, loadPluginRequest *pb.LoadPluginRequest) (*pb.LoadPluginResponse, error) {
	// We would need to find the executable (build as a plugin) and save the reference to that plugin.
	// With the plugin, we should call the necesary function to define the structures given by the user
	// This should be the components, and the views. Then the components should be return

	return nil, nil
}

func (*UserInteraction) ImportFiles(importFileStream pb.UserInteraction_ImportFilesServer) error {
	// We could send via a channel the file paths needed, and via a gorouting be processing them as the
	// user defines. Finally, via another channel, send the component description to the stream

	return nil
}

func (*UserInteraction) Render(renderStream pb.UserInteraction_RenderServer) error {
	// First, we would need to create the SystemInteractionClient to get the Query and SendEvent functions.
	// We would need the view register in the LoadPlugin function, and create the gorouting to continuously
	// be able to send events, and get the scene render

	return nil
}

func (*UserInteraction) Close() {
	// Close the plugin
}
