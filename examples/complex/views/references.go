package views

import (
	"fmt"

	c "github.com/JuanBiancuzzo/own_wiki/examples/complex/components"

	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
	s "github.com/JuanBiancuzzo/own_wiki/src/shared"
)

func GetReferencesViews() []s.ViewInformation {
	return []s.ViewInformation{
		s.GetViewInformation[*ReferenceView](),
		s.GetViewInformation[*ReferencesView](),

		s.GetViewInformation[*ReferencesBookView](),
		s.GetViewInformation[*ReferencesPaperView](),
		s.GetViewInformation[*ReferencesWebView](),
	}
}

type ReferenceView struct {
	Reference c.ReferenceComponent
}

func NewReferenceView(reference c.ReferenceComponent) *ReferenceView {
	return &ReferenceView{
		Reference: reference,
	}
}

func (rv *ReferenceView) Preload(data s.OWData) {
	if reference, ok := data.Query(rv.Reference).(c.ReferenceComponent); ok {
		rv.Reference = reference

	} else {
		data.SendEvent("Failed to get data for reference")
	}
}

func (rv *ReferenceView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	switch rv.Reference.Type {
	case c.RT_BOOK:
		return NewReferencesBookView(rv.Reference)

	case c.RT_PAPER:
		return NewReferencesPaperView(rv.Reference)

	case c.RT_WEB:
		return NewReferencesWebView(rv.Reference)

	default:
		// Send error via outputEvent, to be handle by the system represented in the view
		data.SendEvent(fmt.Sprintf("Failed to load references view, the reference %q doesnt exists", rv.Reference.Type.String()))

		// Maybe make a view to show the error
		return nil
	}
}

type ReferencesView struct {
	References []c.ReferenceComponent

	referencesViews []*ReferenceView
}

func NewReferencesView(references []c.ReferenceComponent) *ReferencesView {
	return &ReferencesView{
		References:      references,
		referencesViews: make([]*ReferenceView, len(references)),
	}
}

func (rv *ReferencesView) Preload(data s.OWData) {
	for i, reference := range rv.References {
		rv.referencesViews[i] = NewReferenceView(reference)
		rv.referencesViews[i].Preload(data)
	}
}

func (rv *ReferencesView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	walkers := make([]*v.LocalWalker[s.OWData], len(rv.referencesViews))
	for i, reference := range rv.referencesViews {
		walkers[i] = v.NewLocalWalker(v.View[s.OWData](reference), world, data)
	}

	for events := range yield() {
		for _, walker := range walkers {
			walker.WalkScene(e.Copy(events))
		}
	}

	return nil
}

// ---+--- Book ---+---
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

func (rv *ReferencesBookView) Preload(data s.OWData) {
	if reference, ok := data.Query(rv.Reference).(c.ReferenceBookComponent); ok {
		rv.Reference = reference

	} else {
		data.SendEvent("Failed to get data for book reference")
	}
}

func (rv *ReferencesBookView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	for range yield() {
	}
	return nil
}

// ---+--- Paper ---+---
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

func (rv *ReferencesPaperView) Preload(data s.OWData) {
	if reference, ok := data.Query(rv.Reference).(c.ReferencePaperComponent); ok {
		rv.Reference = reference

	} else {
		data.SendEvent("Failed to get data for paper reference")
	}
}

func (rv *ReferencesPaperView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	for range yield() {
	}
	return nil
}

// ---+--- Web ---+---
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

func (rv *ReferencesWebView) Preload(data s.OWData) {
	if reference, ok := data.Query(rv.Reference).(c.ReferenceWebComponent); ok {
		rv.Reference = reference

	} else {
		data.SendEvent("Failed to get data for web reference")
	}
}

func (rv *ReferencesWebView) View(world *v.World, data s.OWData, yield v.FnYield) v.View[s.OWData] {
	for range yield() {
	}
	return nil
}
