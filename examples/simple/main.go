package simple_example

import (
	"fmt"

	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
	v "github.com/JuanBiancuzzo/own_wiki/core/views"
)

// Library
type Optional[T any] struct {
	value     T
	has_value bool
}

func NewSome[T any](value T) *Optional[T] {
	return &Optional[T]{
		value:     value,
		has_value: true,
	}
}

func NewNone[T any]() *Optional[T] {
	return &Optional[T]{
		has_value: false,
	}
}

func (o *Optional[T]) Get() (T, bool) {
	return o.value, o.has_value
}

func (o *Optional[T]) Set(value T) {
	o.value = value
	o.has_value = true
}

// Components
type BookComponent struct {
	Name        string
	Author      string
	Category    string
	Year        int
	Description string
}

type ChapterComponent struct {
	Book    *BookComponent
	Number  int
	Name    *Optional[string]
	Summery string
}

// Entities
type LibraryEntity struct {
	Books []BookComponent
}

// View style
var hlineConfig = s.HlineConfig{
	Padding: s.DimConfig{Above: s.ValueEm(.2), Below: s.ValueEm(.5)},
}

// Views
type LibraryView struct {
	Library LibraryEntity
}

func (lv *LibraryView) View(sCtx *s.SceneCtx) (nextView v.View) {
	sCtx.AddMainLayout(s.LayoutConfig{
		Dir: s.VERTICAL_DIR, Padding: s.DimConfig{
			LeftRight: s.ValuePorce(.3), Above: s.ValuePorce(.05),
		},
	}, func() {
		sCtx.AddTitle(s.TitleConfig{Level: 1, Text: "Libreria"})
		sCtx.AddHline(hlineConfig)

		for i, book := range lv.Library.Books {
			if i > 20 {
				break
			}

			if BookBlockButton(sCtx, book) {
				nextView = InitBookView(book)
			}
		}
	})

	if nextView == nil {
		nextView = lv
	}
	return nextView
}

func BookBlockButton(sCtx *s.SceneCtx, book BookComponent) bool {
	return sCtx.AddButton(s.ButtonConfig{
		H: s.Shrink(), W: s.Spand(), RoundedCorners: s.DimConfig{All: s.ValueEm(0.5)},
		Padding: s.DimConfig{All: s.ValueEm(1)}, Margin: s.DimConfig{All: s.ValueEm(1)},
	}, func() {
		BookTitle(sCtx, book, 1)
		sCtx.AddHline(hlineConfig)
		sCtx.AddText(s.TextConfig{Behaviour: s.WRAP_TEXT, Text: book.Description})
	})
}

func BookTitle(sCtx *s.SceneCtx, book BookComponent, level uint) {
	sCtx.AddLayout(s.LayoutConfig{
		H: s.Shrink(), W: s.Spand(), Dir: s.HORIZONTAL_DIR,
	}, func() {
		sCtx.AddTitle(s.TitleConfig{Level: level, Text: fmt.Sprintf("%s by %s", book.Name, book.Author)})
		sCtx.AddBox(s.BoxConfig{H: s.Spand(), W: s.Spand()}, func() {})
		sCtx.AddTitle(s.TitleConfig{Level: level - 2, Text: fmt.Sprintf("at %d", book.Year)})
	})
}

type BookView struct {
	// This would be an anonymous entity
	Book     BookComponent
	Chapters []ChapterComponent
}

func InitBookView(book BookComponent) *BookView {
	return &BookView{
		Book: book,
		Chapters: []ChapterComponent{
			// This is the way to get all the chapter that has as
			// a referece this book
			{Book: &book},
		},
	}
}

func (bv *BookView) View(sCtx *s.SceneCtx) (nextView v.View) {
	sCtx.AddMainLayout(s.LayoutConfig{
		Dir: s.VERTICAL_DIR, Padding: s.DimConfig{
			LeftRight: s.ValuePorce(.3), Above: s.ValueEm(1),
		},
	}, func() {
		BookTitle(sCtx, bv.Book, 1)
		sCtx.AddHline(hlineConfig)
		sCtx.AddText(s.TextConfig{Behaviour: s.WRAP_TEXT, Text: bv.Book.Description})

		for _, chapter := range bv.Chapters {
			ChapterTitle(sCtx, chapter, 2)
			sCtx.AddHline(hlineConfig)
			sCtx.AddText(s.TextConfig{Behaviour: s.WRAP_TEXT, Text: chapter.Summery})
		}
	})

	return bv
}

func ChapterTitle(sCtx *s.SceneCtx, chapter ChapterComponent, level uint) {
	if name, ok := chapter.Name.Get(); ok {
		sCtx.AddTitle(s.TitleConfig{Level: level, Text: fmt.Sprintf("Chapter %d: %s", chapter.Number, name)})
	} else {
		sCtx.AddTitle(s.TitleConfig{Level: level, Text: fmt.Sprintf("Chapter %d", chapter.Number)})
	}
}
