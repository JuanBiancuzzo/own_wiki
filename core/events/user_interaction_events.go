package events

type PromptTextEvent struct {
	Text string
}

func (*PromptTextEvent) isEvent() {}

func NewPromptTextEvent(text string) *PromptTextEvent {
	return &PromptTextEvent{
		Text: text,
	}
}
