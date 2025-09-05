package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	Bloque        string
	Parametros    []string
	Informaciones map[string]FnInformacion
}

func NewEndpoint(bloque string, parametros []string, informaciones map[string]FnInformacion) Endpoint {
	return Endpoint{
		Bloque:        bloque,
		Parametros:    parametros,
		Informaciones: informaciones,
	}
}

// Hacer un endpoint para actualizar, insertar y eliminar
func (e Endpoint) GenerarEndpoint(bdd *b.Bdd, viewNombre string) func(ec echo.Context) error {
	nombreVariables := []string{}
	arrayInformacion := []FnInformacion{}
	for nombreValor := range e.Informaciones {
		nombreVariables = append(nombreVariables, nombreValor)
		arrayInformacion = append(arrayInformacion, e.Informaciones[nombreValor])
	}

	bloque := fmt.Sprintf("%s/%s", viewNombre, e.Bloque)

	return func(ec echo.Context) error {
		valores := make([]string, len(e.Parametros))
		for i, requisito := range e.Parametros {
			valores[i] = ec.QueryParam(requisito)
		}

		data := make(d.ConjuntoDato)
		for i, nombreValor := range nombreVariables {
			informacion := arrayInformacion[i]
			if valor, err := informacion(bdd, valores); err != nil {
				fmt.Printf("Error la informacion %s con error: %v\n", nombreValor, err)
				return err

			} else {
				data[nombreValor] = valor
			}
		}

		return ec.Render(200, bloque, data)
	}
}
