package shared

import (
	"reflect"

	"github.com/JuanBiancuzzo/own_wiki/src/core/systems/file_loader"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

type ComponentInformation reflect.Type

func GetComponentInformation[C any]() ComponentInformation {
	var component C
	return reflect.TypeOf(component)
}

type EntityInformation reflect.Type

func GetEntityInformation[E any]() EntityInformation {
	var entity E
	return reflect.TypeOf(entity)
}

type Entity any

type ViewInformation reflect.Type

func GetViewInformation[V v.View]() ViewInformation {
	var view V
	return reflect.TypeOf(view)
}

type Option[T any] struct{}

type Iterator[T any] struct{}

func NewIterator[T any](elements []T) Iterator[T] {
	return Iterator[T]{}
}

func (r Iterator[T]) Request(amount int) []T {
	return []T{}
}

type Limit[T any] struct {
	request []T
}

func NewLimit[T any](elements []T, amount int) Limit[T] {
	iterator := NewIterator(elements)
	return Limit[T]{
		request: iterator.Request(amount),
	}
}

func (l Limit[T]) Get() []T {
	return l.request
}

type File file_loader.File

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
	RegisterViews() (mainViews []ViewInformation, otherViews []ViewInformation)

	// Given that when importing file there has to be a way to transform them in entities, this
	// is where it happends. This also defines what entity is it wanted to be the main menu. If
	// multiples entities are main menu capable, then it will apear an option to select
	ProcessFile(file File) []Entity
}
