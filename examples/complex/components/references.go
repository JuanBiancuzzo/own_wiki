package components

import (
	"fmt"

	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

type ReferenceType uint8

const (
	RT_BOOK = iota
	RT_PAPER
	RT_WEB
)

func (rt ReferenceType) String() string {
	switch rt {
	case RT_BOOK:
		return "Book"

	case RT_PAPER:
		return "Paper"

	case RT_WEB:
		return "Web"

	default:
		return fmt.Sprintf("Error %d is not a ReferenceType possible value", rt)
	}
}

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
