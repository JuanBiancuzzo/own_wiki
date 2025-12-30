package views

import (
	c "github.com/JuanBiancuzzo/own_wiki/examples/complex/components"

	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

func GetReferencesViews() []s.ViewInformation {
	return []s.ViewInformation{
		s.GetViewInformation[*ReferencesView](),
	}
}

type ReferencesView struct {
	Reference c.ReferenceComponent
}

func (rv *ReferencesView) Preload(outputEvents v.EventHandler) {
	switch rv.Reference.Type {
	case c.RT_BOOK:
		// Preload ReferencesBookView

	case c.RT_PAPER:
		// Preload ReferencesPaperView

	case c.RT_WEB:
		// Preload ReferencesWebView
	}
}

func (rv *ReferencesView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {
	switch rv.Reference.Type {
	case c.RT_BOOK:
		return NewReferencesBookView(rv.Reference)

	case c.RT_PAPER:
		return NewReferencesPaperView(rv.Reference)

	case c.RT_WEB:
		return NewReferencesWebView(rv.Reference)
	}

	return nil
}

type ReferencesBookView struct {
	Reference c.ReferenceBookComponent
}

func NewReferencesBookView(reference c.ReferenceComponent) *ReferencesBookView {
	return &ReferencesBookView{
		Reference: c.ReferenceBookComponent{
			ReferenceComponent: reference,
		},
	}
}

func (rv *ReferencesBookView) Preload(outputEvents v.EventHandler) {}

func (rv *ReferencesBookView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {
	return nil
}

type ReferencesPaperView struct {
	Reference c.ReferencePaperComponent
}

func NewReferencesPaperView(reference c.ReferenceComponent) *ReferencesPaperView {
	return &ReferencesPaperView{
		Reference: c.ReferencePaperComponent{
			ReferenceComponent: reference,
		},
	}
}

func (rv *ReferencesPaperView) Preload(outputEvents v.EventHandler) {}

func (rv *ReferencesPaperView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {
	return nil
}

type ReferencesWebView struct {
	Reference c.ReferenceWebComponent
}

func NewReferencesWebView(reference c.ReferenceComponent) *ReferencesWebView {
	return &ReferencesWebView{
		Reference: c.ReferenceWebComponent{
			ReferenceComponent: reference,
		},
	}
}

func (rv *ReferencesWebView) Preload(outputEvents v.EventHandler) {}

func (rv *ReferencesWebView) View(world *v.World, outputEvents v.EventHandler, yield v.FnYield) v.View {
	return nil
}
