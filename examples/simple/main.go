package plugin

import (
	q "github.com/JuanBiancuzzo/own_wiki/src/core/query"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

// ---+--- Components ---+---
type BookComponent struct {
	Author  string
	Title   string
	Edition int
	Year    int
}

type ChapterComponent struct {
	Book   *BookComponent
	Number int
	Name   s.Option[string]
}

type ReviewComponent struct {
	Book    *BookComponent
	Comment string
}

// ---+--- Entities ---+---
type BookEntity struct {
	Book     BookComponent
	Chapters []ChapterComponent
}

type ReviewEntity struct {
	Book     BookComponent
	Chapters q.Limit[ChapterComponent]
	Review   ReviewComponent
}

// ---+--- Views ---+---

type BookView struct {
	Review ReviewEntity
}

func (bv *BookView) View(world *v.World, creator v.ObjectCreator, yield v.FnYield) v.View {
	return nil
}

type FilterLibraryView struct {
	// This allows the user to get a fix amount of elements
	BookReviews q.Limit[ReviewComponent]
}

func NewFilterLibraryView(amount int) *FilterLibraryView {
	return &FilterLibraryView{
		BookReviews: q.NewLimit([]ReviewComponent{}, amount),
	}
}

func (fv *FilterLibraryView) View(world *v.World, creator v.ObjectCreator, yield v.FnYield) v.View {
	return nil
}

type LibraryView struct {
	// This allows the user to request async an amount of new elements
	Books q.Iterator[BookComponent]
}

func (lv *LibraryView) View(world *v.World, creator v.ObjectCreator, yield v.FnYield) v.View {

	// We create a scene to show a view within this view
	scene := creator.NewScene(NewFilterLibraryView(5), v.DefaultWorldConfiguration())

	// We add the scene to the world os that it is render
	world.MainCamera.ScreenLayout.Add(scene)

	// Esto lo implementamos como un request a la api del sistema
	_ = lv.Books.Request(20)

	for events := range yield() {
		// ...

		if false { // We go to the next view with a given condition
			break
		}

		scene.StepView(events)
	}

	return &LibraryView{
		Books: q.NewIterator([]BookComponent{
			// Vacio sería dejarlo vacio

			// Where condicion simple
			{Author: "Jose"},

			// Caso de un 'or'
			{Author: "Pablo"},
			{Author: "Maria"},

			// Caso de un 'and'
			{Author: "Pablo", Edition: 1},

			// like -> todos los que terminan con 'son'
			{Author: "%son"},

			// Como manejar un rango, y como manejar si quiero recibir todos pero de a poco?
			//  -> Con s.Limit y s.Iterator
		}),
	}
}

// ---+--- Registration and Importing ---+---
type UserDefineStructure struct{}

func (*UserDefineStructure) RegisterComponents() []s.ComponentInformation {
	return []s.ComponentInformation{
		s.GetComponentInformation[BookComponent](),
		s.GetComponentInformation[ChapterComponent](),
		s.GetComponentInformation[ReviewComponent](),
	}
}

func (*UserDefineStructure) RegisterViews() (mainViews []s.ViewInformation, otherViews []s.ViewInformation) {
	mainViews = append(mainViews, []s.ViewInformation{
		s.GetViewInformation[*LibraryView](),
	}...)
	otherViews = append(otherViews, []s.ViewInformation{
		s.GetViewInformation[*BookView](),
		s.GetViewInformation[*FilterLibraryView](),
	}...)
	return mainViews, otherViews
}

func (*UserDefineStructure) ProcessFile(file s.File) []s.Entity {
	// No files are define, or not are important
	return []s.Entity{}
}
