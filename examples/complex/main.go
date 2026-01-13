package plugin

import (
	c "github.com/JuanBiancuzzo/own_wiki/examples/complex/components"
	e "github.com/JuanBiancuzzo/own_wiki/examples/complex/entities"
	v "github.com/JuanBiancuzzo/own_wiki/examples/complex/views"

	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

type UserDefineStructure struct{}

// ---+--- Registration ---+---
func (*UserDefineStructure) RegisterComponents() (components []s.ComponentInformation) {
	components = append(components, c.GetLibraryComponents()...)
	components = append(components, c.GetReferencesComponents()...)

	return components
}

func (*UserDefineStructure) RegisterEntities() []s.EntityInformation {
	return []s.EntityInformation{
		s.GetEntityInformation[e.BookEntity](),
	}
}

func (*UserDefineStructure) RegisterViews() (mainViews []s.ViewInformation, otherViews []s.ViewInformation) {
	mainViews = append(mainViews, s.GetViewInformation[*v.MainView]())

	otherViews = append(otherViews, v.GetLibaryViews()...)
	otherViews = append(otherViews, v.GetReferencesViews()...)

	return mainViews, otherViews
}

// ---+--- Importing ---+---
func (*UserDefineStructure) ProcessFile(file s.File) []s.Entity {
	// No files are define, or not are important
	return []s.Entity{}
}
