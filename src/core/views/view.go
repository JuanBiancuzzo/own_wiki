package views

type View interface {
}

// This structure is capable of waking the state machine define by the sequence
// of views
type ViewWaker struct{}
