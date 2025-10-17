package views

import (
	b "github.com/JuanBiancuzzo/own_wiki/system_protocol/base_de_datos"

	"github.com/labstack/echo/v4"
)

type InfoViews struct {
	RecursosEstaticos map[string]string
	ViewManager       *ViewManager
}

func NewInfoViews(recursosEstaticos map[string]string) (*InfoViews, error) {
	return &InfoViews{
		RecursosEstaticos: recursosEstaticos,
		ViewManager:       NewViewManager(),
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
	for path := range iv.RecursosEstaticos {
		e.Static(path, iv.RecursosEstaticos[path])
	}

	iv.ViewManager.GenerarEndpoints(e, bdd)
}
