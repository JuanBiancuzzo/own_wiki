package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

// Moverlo al lugar correcto

type MainView struct {
}

func (mv *MainView) Preload(outputEvents v.EventHandler) {
	// We could preload the main view of the user

	// la view definida por el usuario deberia estar en el struct de main view
	// y que deberiamos generar un viewWalker para

}

func (mv *MainView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {
	world.Clear()

	mainLayout := world.MainCamera.ScreenLayout
	// Setteariamos los parametros del layout
	// Usariamos la informacion de configuracion para mostrar informacion
	// como el nombre del prjecto, configuracion, path, archivos, etc

	userWorld := v.NewWorld()
	mainLayout.Add(userWorld)

	for events := range yield(world.Render()) {
		unconsume := []e.Event{}
		for event := range events {
			// Ver los eventos y establecer que hacer

			// Si no se hace nada con estos eventos, mandarlos al siguiente paso
			unconsume = append(unconsume, event)
		}

		// Llamar al iterador necesario para renderizar la view del usuario
		// y mandar la informacion
		// Vamos a hacer una query para obtener la informacion necesaria
	}

	return nil
}
