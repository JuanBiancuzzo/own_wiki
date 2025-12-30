package views

import (
	c "github.com/JuanBiancuzzo/own_wiki/examples/complex/components"
	e "github.com/JuanBiancuzzo/own_wiki/examples/complex/entities"

	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

func GetLibaryViews() []s.ViewInformation {
	return []s.ViewInformation{
		s.GetViewInformation[*BookView](),
		s.GetViewInformation[*LibraryView](),
	}
}

type BookView struct {
	// Describimos los componentes que necesitamos como antes definiamos la entidad
	Book e.BookEntity
}

func NewBookView(book c.BookComponent) *BookView {
	return &BookView{
		e.BookEntity{Book: book},
	}
}

func (bv *BookView) Preload(outputEvents v.EventHandler) {}

func (bv *BookView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {
	return nil
}

type LibraryView struct {
	Books s.Iterator[c.BookComponent]
}

func NewLibraryView() *LibraryView {
	return &LibraryView{
		Books: s.NewIterator([]c.BookComponent{}),
	}
}

func (lv *LibraryView) Preload(outputEvents v.EventHandler) {}

func (lv *LibraryView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {
	// Se buscan los libros
	libraryBooks := lv.Books.Request(20)

	// Se muestran

	// Esta la opcion de home
	homeSelected := false
	if homeSelected {
		return &MainView{}
	}

	// Se selecciona alguno, y se cambia a la view de ese libro
	bookSelected := libraryBooks[0]
	return NewBookView(bookSelected)
}
