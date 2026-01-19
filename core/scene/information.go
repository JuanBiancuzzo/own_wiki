package scene

import (
	ev "github.com/JuanBiancuzzo/own_wiki/core/events"
)

type FrameInformation struct {
	Events     []ev.Event
	FrameCount uint
}

func NewFrameInformation(events []ev.Event, frameCount uint) FrameInformation {
	return FrameInformation{
		Events:     events,
		FrameCount: frameCount,
	}
}
