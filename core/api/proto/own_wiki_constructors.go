package core

import (
	"fmt"

	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

func NewRenderRequest(frameInformation s.FrameInformation) (*RenderRequest, error) {
	events := make([]*Event, len(frameInformation.Events))
	for i, event := range frameInformation.Events {
		var err error
		if events[i], err = ConvertFromSystemEvent(event); err != nil {
			return nil, fmt.Errorf("Error while parsing events, with error: %v", err)
		}
	}

	return &RenderRequest{
		Events: events,
		Frame: &FrameInformation{
			FrameCount: uint32(frameInformation.FrameCount),
		},
	}, nil
}
