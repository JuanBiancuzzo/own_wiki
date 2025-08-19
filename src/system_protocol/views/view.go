package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"

	"github.com/labstack/echo/v4"
)

func NewView(bdd *b.Bdd, bloque string, clavesNecesarias []string, informaciones map[string]Informacion) Endpoint {
	nombreVariables := []string{}
	arrayInformacion := []Informacion{}
	for nombreValor := range informaciones {
		nombreVariables = append(nombreVariables, nombreValor)
		arrayInformacion = append(arrayInformacion, informaciones[nombreValor])
	}

	return func(ec echo.Context) error {
		valoresNecesarios := make(map[string]string)
		for _, requisito := range clavesNecesarias {
			valoresNecesarios[requisito] = ec.QueryParam(requisito)
		}

		data := make(d.ConjuntoDato)
		for i, nombreValor := range nombreVariables {
			informacion := arrayInformacion[i]
			if valor, err := informacion.ObtenerInformacion(bdd, valoresNecesarios); err != nil {
				fmt.Printf("Error la informacion %s con error: %v\n", nombreValor, err)
				return err

			} else {
				data[nombreValor] = valor
			}
		}

		return ec.Render(200, bloque, data)
	}
}
