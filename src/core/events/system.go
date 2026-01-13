package events

type CloseEvent struct {
	Motive string
}

func NewCloseEvent(motive string) Event {
	return CloseEvent{
		Motive: motive,
	}
}

/* func (CloseEvent) EqualType(other Event) bool {
	_, ok := other.(CloseEvent)
	return ok
} */

type CreateViewEvent struct {
	ViewName   string
	EntityData any
}

func NewCreateViewEvent(viewName string, entityData any) Event {
	return CreateViewEvent{
		ViewName:   viewName,
		EntityData: entityData,
	}
}

/* func (CreateViewEvent) EqualType(other Event) bool {
	_, ok := other.(CreateViewEvent)
	return ok
} */
