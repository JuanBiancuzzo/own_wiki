package main

import (
	p "plugin"
	"reflect"
	"sync"

	"github.com/hashicorp/go-plugin"

	"github.com/JuanBiancuzzo/own_wiki/src/core/api"
	"github.com/JuanBiancuzzo/own_wiki/src/core/ecv"
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	"github.com/JuanBiancuzzo/own_wiki/src/core/systems/file_loader"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"

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

	Importer *Importer

	World           *v.World
	ObjectCreator   v.ObjectCreator
	EventQueue      chan []e.Event
	WaitInitialView *sync.WaitGroup
}

func NewOwnWiki() *OwnWikiUserStructure {
	var waitGroup sync.WaitGroup
	return &OwnWikiUserStructure{
		Plugin: nil,

		Components: make(map[string]reflect.Type),
		Entities:   make(map[string]reflect.Type),
		Views:      make(map[string]reflect.Type),

		Importer: nil,

		World:           nil,
		ObjectCreator:   nil,
		EventQueue:      nil,
		WaitInitialView: &waitGroup,
	}
}

// ---+--- Register ---+---
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

	mainViews, otherViews := o.Plugin.RegisterViews()
	isMain := []bool{true, false}

	for i, views := range [][]shared.ViewInformation{mainViews, otherViews} {
		for _, view := range views {
			if err := builder.RegisterView(reflect.New(view), isMain[i]); err != nil {
				return api.NewErrorRegisterStructure("Failed to register a view, with error: %v", err), nil
			}

			o.Views[view.Name()] = view
		}
	}

	if !builder.Verify() {
		return api.NewErrorRegisterStructure("Failed to build the system"), nil
	}
	return api.ReturnStructure(*builder), nil
}

// ---+--- Importing ---+---
func (o *OwnWikiUserStructure) InitializeImport(uploader api.Uploader) error {
	o.Importer = NewImporter()

	process := func(file file_loader.File) {
		for _, entity := range o.Plugin.ProcessFile(shared.File(file)) {
			_ = entity // pasarlo a la descripcion de una entidad, tal vez con el builder
			uploader.Upload(ecv.ComponentDescription{})
		}
	}

	file_loader.NewReaderWorker(10, o.Importer.FilePaths, process, o.Importer.WaitGroup)
	return nil
}

func (o *OwnWikiUserStructure) ProcessFile(filePath string) error {
	o.Importer.FilePaths <- filePath
	return nil
}

func (o *OwnWikiUserStructure) FinishImporing() error {
	o.Importer.Close()
	o.Importer = nil
	return nil
}

// ---+--- View Management ---+---
func (o *OwnWikiUserStructure) InitializeViewManeger(worldConfiguration v.WorldConfiguration, system api.OWData) error {
	o.World = v.NewWorld(worldConfiguration)
	o.ObjectCreator = NewObjectCreatorClient(func(view v.View) v.View {
		nameRequestedView := reflect.TypeOf(view).Name()

		// Obtenemos la información necesaria que necesita la view
		entityRequested := o.Views[nameRequestedView]
		entityData, err := system.Query(entityRequested)

		// Agregar esa data a la view dada

		return view
	})

	o.WaitInitialView.Add(1)

	return nil
}

func (o *OwnWikiUserStructure) InitializeView(initialView string, viewData ecv.EntityDescription, system api.OWData) error {
	viewValue := reflect.New(o.Views[initialView])
	view := viewValue.Interface().(v.View) // panics if the view given isnt a view

	o.EventQueue = make(chan []e.Event)
	o.WaitInitialView.Done()

	// Agregar la data a la view
	go func() {
		nextView := view.View(o.World, o.ObjectCreator, func() <-chan []e.Event {
			return o.EventQueue
		})
		o.WaitInitialView.Add(1)

		// Lo cerramos para que las views internas terminen de procesar
		close(o.EventQueue)

		// Request new view
		if nextView != nil {
			system.SendEvent("Load nextView as {nextView}")
		} else {
			system.SendEvent("No new view was created")
		}
	}()

	return nil
}

func (o *OwnWikiUserStructure) WalkScene(events []e.Event) error {
	o.WaitInitialView.Wait()
	o.EventQueue <- events
	return nil
}

func (o *OwnWikiUserStructure) RenderScene() (v.SceneDescription, error) {
	return o.World.Render(), nil
}

// ---+--- Extra functionality ---+---
func (o *OwnWikiUserStructure) Close() {}

// ---+--- Handle plugin ---+---

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
