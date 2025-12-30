package plugin

import (
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

type LibraryEntity struct {
	Books s.Iterator[BookComponent]
}

func NewLibrary(books []BookComponent) LibraryEntity {
	return LibraryEntity{
		Books: s.NewIterator(books),
	}
}

// ---+--- Views ---+---

// Le calculamos el uid del componente, y lo podemos rellenar
// con la informacion que tenemos
type BookView struct {
	// Describimos los componentes que necesitamos como antes definiamos la entidad
	Book ReviewEntity
}

func (bv *BookView) Preload() {}

// Lo puedo hacer como parte de mi cliente
func (bv *BookView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) {

}

// Aca el usuario no lo puede llenar para filtrar como dije antes
type LibraryView struct {
	Library LibraryEntity
}

func (lv *LibraryView) Preload() {
	// Aca podrias ordenar, tendrias q filtrar para el between

	// Por otro lado se podria ir precargando imagenes, o buildeando cosas que fueran necesarias a lo largo de la view
	// Esto implica que deberia haber un evento que defina el inicio de la precarga, ya que sino, se ejecutara este e
	// inmediatamente despues la view
}

func (lv *LibraryView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {

	// Esto lo implementamos como un request a la api del sistema
	_ = lv.Library.Books.Request()

	return &LibraryView{
		NewLibrary([]BookComponent{
			// All seria vacio

			// Where condicion simple
			BookComponent{Author: "Jose"},

			// Caso de un 'or'
			BookComponent{Author: "Pablo"},
			BookComponent{Author: "Maria"},

			// Caso de un 'and'
			BookComponent{Author: "Pablo", Edition: 1},

			// like -> todos los que terminan con 'son'
			BookComponent{Author: "%son"},

			// Como manejar un rango, y como manejar si quiero recibir todos pero de a poco?
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
		s.GetEntityInformation[LibraryEntity](),
	}
}

func (*UserDefineStructure) RegisterViews() (mainViews []s.ViewInformation, otherViews []s.ViewInformation) {
	mainViews = append(mainViews, []s.ViewInformation{
		s.GetViewInformation[LibraryView](),
	})
	otherViews = append(otherViews, []s.ViewInformation{
		s.GetViewInformation[BookView](),
	})

	return mainViews, otherViews
}

func (*UserDefineStructure) ProcessFile(file s.File) []s.Entity {
	// No files are define, or not are important
	return []s.Entity{}
}
