package views

import (
	c "github.com/JuanBiancuzzo/own_wiki/examples/complex/components"
	e "github.com/JuanBiancuzzo/own_wiki/examples/complex/entities"

	q "github.com/JuanBiancuzzo/own_wiki/src/core/query"
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

func (bv *BookView) Preload(data s.OWData) {
	if book, ok := data.Query(bv.Book).(e.BookEntity); ok {
		bv.Book = book

	} else {
		data.SendEvent("Failed to get data for book filter review")
	}
}

func (bv *BookView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	return nil
}

type LibraryView struct {
	Books q.Iterator[c.BookComponent]
}

func NewLibraryView() *LibraryView {
	return &LibraryView{
		Books: q.NewIterator([]c.BookComponent{}),
	}
}

func (lv *LibraryView) Preload(data s.OWData) {
	if books, ok := data.Query(lv.Books).(q.Iterator[c.BookComponent]); ok {
		lv.Books = books

	} else {
		data.SendEvent("Failed to get data for book filter review")
	}
}

func (lv *LibraryView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
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
