package core

// ---+--- QuitEvent ---+---
func NewQuitEvent() *Event {
	return &Event{
		Type: EventType_QUIT,
		Info: &Event_Empty{},
	}
}

// ---+--- PromptTextEvent ---+---
func NewPromptTextEvent(text string) *Event {
	return &Event{
		Type: EventType_PROMT_TEXT,
		Info: &Event_PromptEvent{
			PromptEvent: &PromptTextEvent{
				Text: text,
			},
		},
	}
}
