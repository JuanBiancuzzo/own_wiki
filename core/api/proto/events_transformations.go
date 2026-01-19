package core

import (
	"fmt"

	ev "github.com/JuanBiancuzzo/own_wiki/core/events"
)

func ConvertFromSystemEvent(e ev.Event) (event *Event, err error) {
	switch value := e.(type) {
	// System
	case *ev.QuitEvent:
		event = NewQuitEvent()

	// User Interaction
	case *ev.PromptTextEvent:
		event = NewPromptTextEvent(value.Text)

	default:
		err = fmt.Errorf("Invalid event of type: %T", e)
	}

	return event, err
}

func (e *Event) ConvertToSystemEvent() (event ev.Event, err error) {
	switch e.Type {
	// System
	case EventType_QUIT:
		event = ev.NewQuitEvent()

	// User Interaction
	case EventType_PROMT_TEXT:
		promptInfo := e.GetPromptEvent()
		event = ev.NewPromptTextEvent(promptInfo.Text)

	default:
		err = fmt.Errorf("Invalid event type of: %d", e.Type)
	}

	return event, err
}
