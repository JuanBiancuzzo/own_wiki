package views

import (
	"github.com/JuanBiancuzzo/own_wiki/src/core/ecv"
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	c "github.com/JuanBiancuzzo/own_wiki/src/core/systems/configuration"
	log "github.com/JuanBiancuzzo/own_wiki/src/core/systems/logging"
	u "github.com/JuanBiancuzzo/own_wiki/src/core/user"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

type MainView struct{}

func NewMainView() *MainView {
	return &MainView{}
}

func (mv *MainView) View(world *v.World, yield v.FnYield) {
	world.Clear()

	mainLayout := world.MainCamera.ScreenLayout
	// Setteariamos los parametros del layout
	// Usariamos la informacion de configuracion para mostrar informacion
	// como el nombre del prjecto, configuracion, path, archivos, etc

	for events := range yield() {
		for event := range events {
			// Ver los eventos y establecer que hacer
			_ = event
		}
	}

	// -----------------------------------
	// Obtenemos la información del usuario para buscar el plugin

	userDirectory := ""

	if err := c.LoadUserConfiguration(userDirectory); err != nil {
		log.Error("Failed to load user configuration, with error: %v", err)

	} else {
		c.LoadDefaultUserConfiguration()
	}
	// Registrar estructura dadas por el usuario, y genera las views
	userStructureData, err := u.GetUserDefineData(userDirectory)
	if err != nil {
		log.Error("Failed to get user define data plugin, with error: '%v'", err)
		// Volvemos a preguntar que proyecto quiere porque no esta bien el proyecto que eligió
	}
	defer userStructureData.Close()

	var ecv *ecv.ECV
	if ecv, err = userStructureData.RegisterStructures(); err != nil {
		log.Error("Failed to get components from user, with error: '%v'", err)
		// Notificamos que el proyecto no puede ser cargado, permitirle hacer un reload del proyecto
		// por si ya lo modifico para que funcione
	}

	userWalker := NewUserPluginWalker(userStructureData.Plugin, ecv)

	userScene := v.NewScene(userWalker)
	mainLayout.Add(userScene)

	for events := range yield() {
		unconsume := []e.Event{}
		for event := range events {
			// Ver los eventos y establecer que hacer

			// Si no se hace nada con estos eventos, mandarlos al siguiente paso
			unconsume = append(unconsume, event)
		}

		userWalker.WalkScene(unconsume)
	}
}
