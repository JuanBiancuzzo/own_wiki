package events

type QuitEvent struct{}

func (*QuitEvent) isEvent() {}

func NewQuitEvent() *QuitEvent {
	return &QuitEvent{}
}
