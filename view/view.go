package view

import (
	"reflect"
)

type SceneRepresentation []any

type Event int
type FnYield func(ops *SceneOperation) <-chan []Event

type View interface {
	// If not next scene, then nextScene should be nil
	View(scene *Scene, yield FnYield) (nextScene *ChangeSceneOperation)
}

// Cambio de escena, Modificacion o eliminacion de componente
type ChangeSceneOperation struct {
	ViewName string
	EntityID uint64
}

func ChangeSceneOp[V View](entity any) *ChangeSceneOperation {
	var view V
	viewInfo := reflect.TypeOf(view)

	return &ChangeSceneOperation{
		ViewName: viewInfo.Name(),
		EntityID: 0, // Hacer que sea el hash de los elementos importantes
	}
}

type ChangeEntityOperation struct{}

type SceneOperation struct {
	SceneCaracteristics *SceneRepresentation

	ChangeScene *ChangeSceneOperation
	EndScene    bool

	Operations []ChangeEntityOperation
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
