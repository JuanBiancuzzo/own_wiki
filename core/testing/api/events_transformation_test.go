package api_test

import (
	"testing"

	"github.com/go-test/deep"

	pb "github.com/JuanBiancuzzo/own_wiki/core/api/proto"
	ev "github.com/JuanBiancuzzo/own_wiki/core/events"
)

// ---+--- tests ---+---

func TestSystemEvents(t *testing.T) {
	systemEvent := ev.NewQuitEvent()

	expectEvent := pb.NewQuitEvent()

	if event, err := pb.ConvertFromSystemEvent(systemEvent); err != nil {
		t.Errorf("While converting to Event, got the error: %v", err)

	} else if diff := deep.Equal(expectEvent, event); diff != nil {
		t.Error(diff)

	} else if systemEventGen, err := event.ConvertToSystemEvent(); err != nil {
		t.Errorf("While converting to System Event, got the error: %v", err)

	} else if diff := deep.Equal(systemEvent, systemEventGen); diff != nil {
		t.Error(diff)
	}
}

func TestUserInteractionEvents(t *testing.T) {
	systemEvent := ev.NewPromptTextEvent("this is a prompt text")

	expectEvent := pb.NewPromptTextEvent("this is a prompt text")

	if event, err := pb.ConvertFromSystemEvent(systemEvent); err != nil {
		t.Errorf("While converting to Event, got the error: %v", err)

	} else if diff := deep.Equal(expectEvent, event); diff != nil {
		t.Error(diff)

	} else if systemEventGen, err := event.ConvertToSystemEvent(); err != nil {
		t.Errorf("While converting to System Event, got the error: %v", err)

	} else if diff := deep.Equal(systemEvent, systemEventGen); diff != nil {
		t.Error(diff)
	}
}
