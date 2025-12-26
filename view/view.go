package view

import (
	"reflect"
)

type SceneRepresentation []any

type Event int
type FnYield func(ops []*SceneOperation) <-chan []Event

type View interface {
	View(scene *Scene, yield FnYield) (nextScene *SceneOperation)
}

// Cambio de escena, Modificacion o eliminacion de componente
type SceneOperation struct {
	// Cambio de escena
	viewName string
	entityID uint64
}

func ChangeScene[V View](entity any) *SceneOperation {
	var view V
	viewInfo := reflect.TypeOf(view)

	return &SceneOperation{
		viewName: viewInfo.Name(),
		entityID: 0, // Hacer que sea el hash de los elementos importantes
	}
}

// -+- Bloques -+-
type Heading struct {
	Level uint8
	Data  string
}

func NewHeading(level uint8, text string) *Heading {
	return &Heading{
		Level: level,
		Data:  text,
	}
}

type Text struct {
	Data string
}

func NewText(text string) *Text {
	return &Text{
		Data: text,
	}
}

func (t *Text) ChangeText(text string) {
	t.Data = text
}

type Link struct {
	data string
}
