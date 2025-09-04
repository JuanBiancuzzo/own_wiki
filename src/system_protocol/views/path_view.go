package views

import (
	"fmt"
	"slices"
)

type PathView struct {
	Views []string
}

func NewPathView() *PathView {
	return &PathView{
		Views: []string{},
	}
}

func (pv *PathView) AgregarView(view string) error {
	if slices.Contains(pv.Views, view) {
		return fmt.Errorf("ya se cargo esa view")
	}

	pv.Views = append(pv.Views, view)
	return nil
}

func (pv *PathView) CreateURLPathView(view string) string {
	if !slices.Contains(pv.Views, view) {
		return "ERROR - No existe view"
	}

	return fmt.Sprintf("/%s?redirect=true", view)
}
