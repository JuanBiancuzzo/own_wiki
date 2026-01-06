package views

import (
	"github.com/JuanBiancuzzo/own_wiki/src/core/api"

	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"

	log "github.com/JuanBiancuzzo/own_wiki/src/core/systems/logging"
)

type UserPluginWalker struct {
	plugin       api.UserStructureData
	outputEvents v.EventHandler
	request      v.RequestView
}

func (uw *UserPluginWalker) InitializeView(view v.View, world *v.World) {
	if err := uw.plugin.InitializeView(view, world, uw.outputEvents, uw.request); err != nil {
		log.Error("Failed to initialize view '%v', with error: %v", view, err)
	}
}

func (uw *UserPluginWalker) Preload(uid v.ViewId, view v.View) {
	if err := uw.plugin.Prelaod(uid, view); err != nil {
		log.Error("Failed to preload the view '%v' of uid: %d, with error: %v", view, uid, err)
	}
}

func (uw *UserPluginWalker) WalkScene(events []e.Event) {
	if err := uw.plugin.WalkScene(events); err != nil {
		log.Error("Failed to walk the scene in the client, with error: %v", err)
	}
}

func (uw *UserPluginWalker) Render() v.SceneRepresentation {
	if scene, err := uw.plugin.RenderScene(); err != nil {
		log.Error("Failed to render a scene, with error: %v", err)
		return nil

	} else {
		return scene
	}
}
