package views

import (
	b "own_wiki/system_protocol/bass_de_datos"

	"github.com/labstack/echo/v4"
)

type InfoViews struct {
	PathCss      string
	PathImagenes string
	PathView     *PathView
	Views        []View
}

func NewInfoViews(views []View, pathCss, pathImagenes string) *InfoViews {
	pathView := NewPathView()
	for _, view := range views {
		view.RegistrarEndpoints(pathView)
	}

	return &InfoViews{
		Views:        views,
		PathCss:      pathCss,
		PathImagenes: pathImagenes,
		PathView:     pathView,
	}
}

func (t *InfoViews) RegistrarRenderer(e *echo.Echo, carpetaRoot string) error {
	var err error
	if e.Renderer, err = NewTemplate(t.Views, t.PathView); err != nil {
		return err
	}
	return nil
}

func (t *InfoViews) GenerarEndpoints(e *echo.Echo, bdd *b.Bdd) {
	for _, view := range t.Views {
		view.GenerarEndpoints(e, bdd)
	}
}
