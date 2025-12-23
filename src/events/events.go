package events

type Event any

type CharacterEvent struct {
	Char rune
}

func NewCharacterEvent(char rune) CharacterEvent {
	return CharacterEvent{
		Char: char,
	}
}
