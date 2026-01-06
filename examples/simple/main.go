package plugin

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
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
// Conjunto de components que definen la busqueda
type BookEntity struct {
	Book     BookComponent
	Chapters []ChapterComponent // Si esta vacio se toma todos los casos que acompañen el libro
}

type ReviewEntity struct {
	Book     BookComponent
	Chapters s.Limit[ChapterComponent]
	Review   ReviewComponent
}

// ---+--- Views ---+---

// Le calculamos el uid del componente, y lo podemos rellenar
// con la informacion que tenemos
type BookView struct {
	// Describimos los componentes que necesitamos como antes definiamos la entidad
	Book ReviewEntity
}

func (bv *BookView) Preload(outputEvents v.EventHandler) {}

// Lo puedo hacer como parte de mi cliente
func (bv *BookView) View(world *v.World, outputEvents v.EventHandler, requestView v.RequestView, yield v.FnYield) v.View {
	return nil
}

// En el caso de hacer una referencia a componentes, y/u operaciones a entidades
// se toma como una nueva entidad anonima
type FilterLibraryView struct {
	BookReviews s.Limit[ReviewComponent]
}

func (fv *FilterLibraryView) Preload(outputEvents v.EventHandler) {}

func (fv *FilterLibraryView) View(world *v.World, outputEvents v.EventHandler, requestView v.RequestView, yield v.FnYield) v.View {
	return nil
}

// Aca el usuario no lo puede llenar para filtrar como dije antes
type LibraryView struct {
	Books s.Iterator[BookComponent]

	reviewView v.View
}

func (lv *LibraryView) Preload(outputEvents v.EventHandler) {
	// Aca podrias ordenar, tendrias q filtrar para el between

	// Por otro lado se podria ir precargando imagenes, o buildeando cosas que fueran necesarias a lo largo de la view
	// Esto implica que deberia haber un evento que defina el inicio de la precarga, ya que sino, se ejecutara este e
	// inmediatamente despues la view

	lv.reviewView = &FilterLibraryView{s.NewLimit([]ReviewComponent{}, 5)}
	// mandar evento preload esta view
	// outputEvents.PushEvent(e.PreloadView(lv.reviewView))
}

func (lv *LibraryView) View(world *v.World, outputEvents v.EventHandler, requestView v.RequestView, yield v.FnYield) v.View {

	// Deberia chequear si fue preloadeada, pero como hacemos para que este nuevo walker tenga
	// registro de que fue preloadeado
	generalReviews := v.NewLocalWalker(lv.reviewView, world, outputEvents, requestView)

	// Esto lo implementamos como un request a la api del sistema
	_ = lv.Books.Request(20)

	// Al final se quiere mostrar la view, FilterLibraryView
	generalReviews.WalkScene([]e.Event{})

	return &LibraryView{
		Books: s.NewIterator([]BookComponent{
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
		// s.GetViewInformation[*FilterLibraryView](),
	}...)
	return mainViews, otherViews
}

func (*UserDefineStructure) ProcessFile(file s.File) []s.Entity {
	// No files are define, or not are important
	return []s.Entity{}
}
