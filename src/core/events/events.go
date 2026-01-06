package events

type Event any

func NewCloseEvent(motive string) Event { return nil }

func Copy(events []Event) []Event {
	return events
}
