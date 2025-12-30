package views

import v "github.com/JuanBiancuzzo/own_wiki/src/core/views"

type MainView struct{}

func (mv *MainView) Preload(outputEvents v.EventHandler) {}

func (mv *MainView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {

	librarySelected := false
	if librarySelected {
		return NewLibraryView()
	}

	return nil
}
