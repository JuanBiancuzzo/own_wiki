package plugin

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
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

func (bv *BookView) Preload(data s.OWData) {
	if book, ok := data.Query(bv.Review).(ReviewEntity); ok {
		bv.Review = book

	} else {
		data.SendEvent("Failed to get data for book review")
	}
}

func (bv *BookView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	return nil
}

// En el caso de hacer una referencia a componentes, y/u operaciones a entidades
// se toma como una nueva entidad anonima
type FilterLibraryView struct {
	BookReviews q.Limit[ReviewComponent] // Esto se coincidera una entitdad anonima
}

func (fv *FilterLibraryView) Preload(data s.OWData) {
	if books, ok := data.Query(fv.BookReviews).(q.Limit[ReviewComponent]); ok {
		fv.BookReviews = books

	} else {
		data.SendEvent("Failed to get data for book filter review")
	}
}

func (fv *FilterLibraryView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	return nil
}

// Aca el usuario no lo puede llenar para filtrar como dije antes
type LibraryView struct {
	Books q.Iterator[BookComponent]

	reviewView v.View[s.OWData]
}

func (lv *LibraryView) Preload(data s.OWData) {
	if books, ok := data.Query(lv.Books).(q.Iterator[BookComponent]); ok {
		lv.Books = books

	} else {
		data.SendEvent("Failed to get data for book filter review")
	}

	// Por otro lado se podria ir precargando imagenes, o buildeando cosas que fueran necesarias a lo largo de la view
	// Esto implica que deberia haber un evento que defina el inicio de la precarga, ya que sino, se ejecutara este e
	// inmediatamente despues la view

	lv.reviewView = &FilterLibraryView{q.NewLimit([]ReviewComponent{}, 5)}
	lv.reviewView.Preload(data)
}

func (lv *LibraryView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {

	// Deberia chequear si fue preloadeada, pero como hacemos para que este nuevo walker tenga
	// registro de que fue preloadeado
	generalReviews := v.NewLocalWalker(lv.reviewView, world, data)

	// Esto lo implementamos como un request a la api del sistema
	_ = lv.Books.Request(20)

	// Al final se quiere mostrar la view, FilterLibraryView
	generalReviews.WalkScene([]e.Event{})

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

func (*UserDefineStructure) RegisterEntities() []s.EntityInformation {
	return []s.EntityInformation{
		s.GetEntityInformation[BookEntity](),
		s.GetEntityInformation[ReviewEntity](),
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
