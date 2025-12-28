package main

import (
	p "plugin"
	"reflect"

	"github.com/hashicorp/go-plugin"

	"github.com/JuanBiancuzzo/own_wiki/src/core/api"
	"github.com/JuanBiancuzzo/own_wiki/src/core/ecv"

	"github.com/JuanBiancuzzo/own_wiki/src/shared"
)

// go build -buildmode=plugin -o plugin/plugin.so ./path/to/plugin
const PLUGIN_PATH string = "plugin/plugin.so"
const PLUGIN_STRUCTURE_NAME string = "UserDefineStructure"

type OwnWikiUserStructure struct {
	Plugin shared.UserDefineStructure

	Components map[string]reflect.Type
	Entities   map[string]reflect.Type
	Views      map[string]reflect.Type
}

func NewOwnWiki() *OwnWikiUserStructure {
	return &OwnWikiUserStructure{
		Plugin: nil,

		Components: make(map[string]reflect.Type),
		Entities:   make(map[string]reflect.Type),
		Views:      make(map[string]reflect.Type),
	}
}

func (o *OwnWikiUserStructure) LoadPlugin(path string) (api.ErrorLoadPath, error) {
	if userPlugin, err := p.Open(path); err != nil {
		return api.NewErrorLoadPath("Plugin not found, with error: %v", err), nil

	} else if userDefineStructure, err := userPlugin.Lookup(PLUGIN_STRUCTURE_NAME); err != nil {
		return api.NewErrorLoadPath(
			"In the plugin there is no %s structure defining the structure of the program, with error: %v", PLUGIN_STRUCTURE_NAME, err,
		), nil

	} else if userStructure, ok := userDefineStructure.(shared.UserDefineStructure); !ok {
		return api.NewErrorLoadPath(
			"The %s structure does not define the interfaces needed, with error: %v", PLUGIN_STRUCTURE_NAME, err,
		), nil

	} else {
		o.Plugin = userStructure
		return api.NoErrorLoadPath(), nil
	}
}

func (o *OwnWikiUserStructure) RegisterStructures() (api.ReturnRegisterStructure, error) {
	if o.Plugin == nil {
		return api.NewErrorRegisterStructure("The plugin was not loaded"), nil
	}

	builder := ecv.NewECVBuilder()

	for _, component := range o.Plugin.RegisterComponents() {
		if err := builder.RegisterComponent(reflect.New(component)); err != nil {
			return api.NewErrorRegisterStructure("Failed to register a component, with error: %v", err), nil
		}

		o.Components[component.Name()] = component
	}

	for _, entity := range o.Plugin.RegisterEntities() {
		if err := builder.RegisterEntity(reflect.New(entity)); err != nil {
			return api.NewErrorRegisterStructure("Failed to register an entity, with error: %v", err), nil
		}

		o.Entities[entity.Name()] = entity
	}

	for entity, view := range o.Plugin.RegisterViews() {
		if err := builder.RegisterView(reflect.New(entity), reflect.New(view)); err != nil {
			return api.NewErrorRegisterStructure("Failed to register a view, with error: %v", err), nil
		}

		o.Views[view.Name()] = view
	}

	if system, err := builder.BuildECV(); err != nil {
		return api.NewErrorRegisterStructure("Failed to build the system, with error: %v", err), nil

	} else {
		return api.ReturnStructure(system), nil
	}
}

// View

func (o *OwnWikiUserStructure) Close() {}

func main() {
	ownWiki := NewOwnWiki()
	defer ownWiki.Close()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: api.Handshake,
		Plugins: map[string]plugin.Plugin{
			"userStructureData": &api.UserStructureDataPlugin{Impl: ownWiki},
		},
	})
}
