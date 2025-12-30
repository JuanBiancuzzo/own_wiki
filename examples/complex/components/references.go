package components

import (
	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

type ReferenceType uint8

const (
	RT_BOOK = iota
	RT_PAPER
	RT_WEB
)

func GetReferencesComponents() []s.ComponentInformation {
	return []s.ComponentInformation{
		s.GetComponentInformation[ReferenceComponent](),

		s.GetComponentInformation[ReferenceBookComponent](),
		s.GetComponentInformation[ReferencePaperComponent](),
		s.GetComponentInformation[ReferenceWebComponent](),
	}
}

type ReferenceComponent struct {
	Type   ReferenceType
	Number uint
}

type ReferenceBookComponent struct {
	ReferenceComponent

	Book *BookComponent
}

type ReferencePaperComponent struct {
	ReferenceComponent

	Paper *PaperComponent
}

type ReferenceWebComponent struct {
	ReferenceComponent

	Url  string
	Date string
}
