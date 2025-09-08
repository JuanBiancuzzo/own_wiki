package views

import (
	b "own_wiki/system_protocol/base_de_datos"

	"github.com/labstack/echo/v4"
)

type InfoViews struct {
	PathCss      string
	PathImagenes string
	ViewManager  *ViewManager
}

func NewInfoViews(pathCss, pathImagenes string) (*InfoViews, error) {
	return &InfoViews{
		PathCss:      pathCss,
		PathImagenes: pathImagenes,
		ViewManager:  NewViewManager(),
	}, nil
}

func (iv *InfoViews) AgregarView(view View) error {
	return iv.ViewManager.Agregar(view)
}

func (iv *InfoViews) RegistrarRenderer(e *echo.Echo, carpetaRoot string) (err error) {
	e.Renderer, err = NewTemplate(iv.ViewManager)
	return err
}

func (iv *InfoViews) GenerarEndpoints(e *echo.Echo, bdd *b.Bdd) {
	iv.ViewManager.GenerarEndpoints(e, bdd)
}
