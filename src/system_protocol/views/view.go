package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"

	"github.com/labstack/echo/v4"
)

type View struct {
	Nombre string
	Bloque string
	Bdd    *b.Bdd

	clavesNecesarias []string
	nombreVariables  []string
	informaciones    []Informacion
	multiples        map[string]Endpoint
}

func NewView(bdd *b.Bdd, nombre, bloque string, clavesNecesarias []string, informaciones map[string]Informacion, multiples map[string]Endpoint) View {
	nombreVariables := []string{}
	arrayInformacion := []Informacion{}
	for nombreValor := range informaciones {
		nombreVariables = append(nombreVariables, nombreValor)
		arrayInformacion = append(arrayInformacion, informaciones[nombreValor])
	}

	return View{
		Nombre: nombre,
		Bloque: bloque,
		Bdd:    bdd,

		clavesNecesarias: clavesNecesarias,
		nombreVariables:  nombreVariables,
		informaciones:    arrayInformacion,
		multiples:        multiples,
	}
}

func (v View) GenerarEndpoint(ruta string, e *echo.Echo) {
	for subPath := range v.multiples {
		multiple := v.multiples[subPath]

		e.GET(fmt.Sprintf("%s/%s", ruta, subPath), multiple.GenerarEndpoint)
	}

	e.GET(ruta, func(ec echo.Context) error {
		valoresNecesarios := make(map[string]string)
		for _, requisito := range v.clavesNecesarias {
			valoresNecesarios[requisito] = ec.QueryParam(requisito)
		}

		data := make(d.ConjuntoDato)
		for i, nombreValor := range v.nombreVariables {
			informacion := v.informaciones[i]
			if valor, err := informacion.ObtenerInformacion(v.Bdd, valoresNecesarios); err != nil {
				fmt.Printf("Error al utilizar endpoint /%s, dado la informacion %s con error: %v\n", v.Nombre, nombreValor, err)
				return err

			} else {
				data[nombreValor] = valor
			}
		}

		return ec.Render(200, v.Bloque, data)
	})
}
