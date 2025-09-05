package views

import (
	b "own_wiki/system_protocol/bass_de_datos"

	"github.com/labstack/echo/v4"
)

type InfoViews struct {
	PathCss      string
	PathImagenes string
	Views        []View

	PathEndpoint *PathEndpoint
	PathView     *PathView
}

func NewInfoViews(views []View, pathCss, pathImagenes string) (*InfoViews, error) {
	pathEndpoint := NewPathEndpoint()
	pathView := NewPathView()
	for _, view := range views {
		if err := pathView.AgregarView(view.Nombre); err != nil {
			return nil, err
		}

		if err := view.RegistrarEndpoints(pathEndpoint); err != nil {
			return nil, err
		}
	}

	return &InfoViews{
		Views:        views,
		PathCss:      pathCss,
		PathImagenes: pathImagenes,

		PathEndpoint: pathEndpoint,
		PathView:     pathView,
	}, nil
}

func (t *InfoViews) RegistrarRenderer(e *echo.Echo, carpetaRoot string) error {
	var err error
	if e.Renderer, err = NewTemplate(t.Views, t.PathView, t.PathEndpoint); err != nil {
		return err
	}
	return nil
}

func (t *InfoViews) GenerarEndpoints(e *echo.Echo, bdd *b.Bdd) {
	for _, view := range t.Views {
		view.GenerarEndpoints(e, bdd)
	}
}
