package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

type MainView struct {
	userView v.ViewWalker
}

func (mv *MainView) Preload(outputEvents v.EventHandler) {
	// Se deberia preloadear cosas de la configuracion
}

func (mv *MainView) View(world *v.World, outputEvents v.EventHandler, requestView v.RequestView, yield v.FnYield) v.View {
	world.Clear()

	mainLayout := world.MainCamera.ScreenLayout
	// Setteariamos los parametros del layout
	// Usariamos la informacion de configuracion para mostrar informacion
	// como el nombre del prjecto, configuracion, path, archivos, etc

	// Idea, asociar un world a un walker, de esa forma podemos hacer
	// que directamente se renderise el world, llama al walker
	userScene := v.NewScene(mv.userView)
	mainLayout.Add(userScene)

	for events := range yield() {
		unconsume := []e.Event{}
		for event := range events {
			// Ver los eventos y establecer que hacer

			// Si no se hace nada con estos eventos, mandarlos al siguiente paso
			unconsume = append(unconsume, event)
		}

		mv.userView.WalkScene(unconsume)
	}

	return nil
}
