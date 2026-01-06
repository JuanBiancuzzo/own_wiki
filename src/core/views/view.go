package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
)

type Renderable interface {
	Render() SceneRepresentation
}

type SceneRepresentation any

type FnYield func() []e.Event

/*
Lo que busco es crear una interfaz de una view, esta deberia recibir un "mundo" el cual
deberia tener una camara por default, llamda main, y en esta podes attachearte para renderizar ui
en dos dimensiones clavada a la camara

Un mundo puede tener dentro de el, otra superficie, la cual puede ser un mundo. De esta forma podemos hacer
que dentro del mismo programa se muestre lo que el usuario queria, agregando la funcionalidad minima alrededor
*/
type View[Data any] interface {
	// When initializing the view, it may be useful to initilize the world necesary layout and
	// references to elements that are needed. And getting the data to use
	Preload(Data)

	// View is the way to render and create a representation of the data
	// A nil return value would be that there isnt a next view
	View(*World, Data, FnYield) View[Data]
}
