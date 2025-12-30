package components

import (
	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

func GetLibraryComponents() []s.ComponentInformation {
	return []s.ComponentInformation{
		s.GetComponentInformation[BookComponent](),
		s.GetComponentInformation[ChapterComponent](),
		s.GetComponentInformation[PaperComponent](),
	}
}

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

type PaperComponent struct {
	Author    string
	Title     string
	Year      int
	Editorial string
}
