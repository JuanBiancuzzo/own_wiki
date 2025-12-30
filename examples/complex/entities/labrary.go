package entities

import (
	c "github.com/JuanBiancuzzo/own_wiki/examples/complex/components"
)

type BookEntity struct {
	Book     c.BookComponent
	Chapters []c.ChapterComponent
}
