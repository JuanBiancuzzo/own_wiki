package views

import (
	"fmt"
	b "own_wiki/system_protocol/base_de_datos"
	d "own_wiki/system_protocol/dependencias"

	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	Bloque        string
	Parametros    []string
	Informaciones map[string]FnInformacion
}

func NewEndpoint(nombreView, bloque string, parametros []string, informaciones map[string]FnInformacion) Endpoint {
	return Endpoint{
		Bloque:        fmt.Sprintf("%s/%s", nombreView, bloque),
		Parametros:    parametros,
		Informaciones: informaciones,
	}
}

// Hacer un endpoint para actualizar, insertar y eliminar
func (e Endpoint) GenerarEndpoint(bdd *b.Bdd) func(ec echo.Context) error {
	nombreVariables := []string{}
	arrayInformacion := []FnInformacion{}
	for nombreValor := range e.Informaciones {
		nombreVariables = append(nombreVariables, nombreValor)
		arrayInformacion = append(arrayInformacion, e.Informaciones[nombreValor])
	}

	valores := make([]string, len(e.Parametros))
	data := make(d.ConjuntoDato)

	return func(ec echo.Context) (err error) {
		for i, requisito := range e.Parametros {
			valores[i] = ec.QueryParam(requisito)
		}

		for i, nombreValor := range nombreVariables {
			informacion := arrayInformacion[i]
			if data[nombreValor], err = informacion(bdd, valores); err != nil {
				fmt.Printf("Error la informacion %s con error: %v\n", nombreValor, err)
				return err
			}
		}

		if err = bdd.Checkpoint(b.TC_PASSIVE); err != nil {
			fmt.Printf("Error al hacer checkpoint con error: %v\n", err)
			return err
		}

		return ec.Render(200, e.Bloque, data)
	}
}
