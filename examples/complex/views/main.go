package views

import (
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

type MainView struct{}

func (mv *MainView) Preload(data s.OWData) {}

func (mv *MainView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	librarySelected := false
	if librarySelected {
		return NewLibraryView()
	}

	return nil
}
