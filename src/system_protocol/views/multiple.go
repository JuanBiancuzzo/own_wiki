package views

import "github.com/labstack/echo/v4"

/*
	Construir Multiple y InformacionTabla para trabajar en conjunto, donde
	la informacion tiene las variables para preparar el hmtl para obtener los elementos
	y multiple genera un endpoint para ir buscando elementos en la bdd y despues
	ir mandandolas si se le pega a este endpoint
*/
type Multiple struct {
}

func NewMultiple() Multiple {
	return Multiple{}
}

func (m *Multiple) GenerarEndpoint(ec echo.Context) error {

	return nil
}
