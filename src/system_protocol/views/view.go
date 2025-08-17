package views

import (
	b "own_wiki/system_protocol/bass_de_datos"

	"github.com/labstack/echo/v4"
)

type View struct {
	Nombre string
	Bdd    *b.Bdd
}

func NewView(nombre string) View {
	return View{
		Nombre: nombre,
	}
}

func (v View) GenerarEndpoint(ec echo.Context) error {
	return nil
}
