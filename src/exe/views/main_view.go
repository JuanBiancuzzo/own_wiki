package views

import (
	"github.com/JuanBiancuzzo/own_wiki/src/core/api"
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

type MainView struct {
	UserView v.ViewWalker[api.OWData]
}

func (mv *MainView) Preload(data api.OWData) {
	// Se deberia preloadear cosas de la configuracion
}

func (mv *MainView) View(world *v.World, data api.OWData, yield v.FnYield) v.View[api.OWData] {
	world.Clear()

	mainLayout := world.MainCamera.ScreenLayout
	// Setteariamos los parametros del layout
	// Usariamos la informacion de configuracion para mostrar informacion
	// como el nombre del prjecto, configuracion, path, archivos, etc

	// Idea, asociar un world a un walker, de esa forma podemos hacer
	// que directamente se renderise el world, llama al walker
	userScene := v.NewScene(mv.UserView)
	mainLayout.Add(userScene)

	for events := range yield() {
		unconsume := []e.Event{}
		for event := range events {
			// Ver los eventos y establecer que hacer

			// Si no se hace nada con estos eventos, mandarlos al siguiente paso
			unconsume = append(unconsume, event)
		}

		mv.UserView.WalkScene(unconsume)
	}

	return nil
}
