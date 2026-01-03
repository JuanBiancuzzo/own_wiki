package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
)

type Renderable interface {
	Render() SceneRepresentation
}

type EventHandler interface {
	// This method is to be capable to send modifications to the system in a view
	PushEvent(event e.Event) error
}

type ViewId uint64

type RequestView interface {
	// La dataView generada no puede haber sido preloadeada porque se
	// actualizan los valores ocultos, por ende no son mandados
	Request(requestedView View) (uid ViewId, dataView View)
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
type View interface {
	// When initializing the view, it may be useful to initilize the world necesary layout and
	// references to elements that are needed. And getting the data to use
	Preload(outputEvents EventHandler)

	// View is the way to render and create a representation of the data
	// A nil return value would be that there isnt a next view
	View(world *World, outputEvents EventHandler, yield FnYield) View
}
