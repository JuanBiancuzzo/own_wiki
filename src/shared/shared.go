package shared

import (
	"reflect"

	"github.com/JuanBiancuzzo/own_wiki/src/core/systems/file_loader"
)

type ComponentInformation reflect.Type

type EntityInformation reflect.Type
type Entity struct {
	// Entity struct register with each component fill
	Entity     any
	IsMainManu bool
}

type ViewInformation reflect.Type

/*
This interface lets the user define the components, each entity and view for the project.
*/
type UserDefineStructure interface {
	// The components are the smallest data storage given by the system. They can depende on
	// each other, but there has to be a way to constructe them with out an infinite loop
	RegisterComponents() []ComponentInformation

	// Entities are the composite of components, which are capable of being shown
	RegisterEntities() []EntityInformation

	// Views are the representation of a entity to be shown by the program in the platform
	// define at compilation time
	RegisterViews() map[EntityInformation]ViewInformation

	// Given that when importing file there has to be a way to transform them in entities, this
	// is where it happends. This also defines what entity is it wanted to be the main menu. If
	// multiples entities are main menu capable, then it will apear an option to select
	ProcessFile(file file_loader.File) []Entity
}
