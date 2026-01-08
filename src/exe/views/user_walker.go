package views

import (
	"reflect"

	"github.com/JuanBiancuzzo/own_wiki/src/core/api"

	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"

	log "github.com/JuanBiancuzzo/own_wiki/src/core/systems/logging"
)

type UserPluginWalker struct {
	plugin api.UserStructureData
	data   api.OWData
}

func NewUserPluginWalker(plugin api.UserStructureData, data api.OWData) *UserPluginWalker {
	return &UserPluginWalker{
		plugin: plugin,
		data:   data,
	}
}

func (uw *UserPluginWalker) InitializeView(view v.View[api.OWData]) {
	world := v.NewWorld(v.WorldConfiguration(0))

	viewName := reflect.TypeOf(view).Name()
	if err := uw.plugin.InitializeView(viewName, world, uw.data); err != nil {
		log.Error("Failed to initialize view '%v', with error: %v", view, err)
	}
}

func (uw *UserPluginWalker) WalkScene(events []e.Event) bool {
	if err := uw.plugin.WalkScene(events); err != nil {
		log.Error("Failed to walk the scene in the client, with error: %v", err)
		return false
	}
	return true
}

func (uw *UserPluginWalker) Render() v.SceneRepresentation {
	if scene, err := uw.plugin.RenderScene(); err != nil {
		log.Error("Failed to render a scene, with error: %v", err)
		return nil

	} else {
		return scene
	}
}
