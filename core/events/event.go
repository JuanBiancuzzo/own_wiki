package events

// The types that holds the interface
// System
// 	*QuitEvent
//
// User interface:
// 	*PromptTextEvent
type Event isEvent

type isEvent interface {
	isEvent()
}
