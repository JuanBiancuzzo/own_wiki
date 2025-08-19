package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"

	"github.com/labstack/echo/v4"
)

type View struct {
	Nombre string
	Bloque string
	Bdd    *b.Bdd

	clavesNecesarias []string
	informaciones    map[string]Informacion
}

type DataView map[string]any

func NewView(bdd *b.Bdd, nombre, bloque string, requisitos []string, informaciones map[string]Informacion) View {
	fmt.Println(nombre)
	return View{
		Nombre:           nombre,
		Bloque:           bloque,
		Bdd:              bdd,
		clavesNecesarias: requisitos,
		informaciones:    informaciones,
	}
}

func (v View) GenerarEndpoint(ruta string, e *echo.Echo) {
	nombreVariables := []string{}
	informaciones := []Informacion{}

	for nombreValor := range v.informaciones {
		nombreVariables = append(nombreVariables, nombreValor)
		informaciones = append(informaciones, v.informaciones[nombreValor])
	}

	e.GET(ruta, func(ec echo.Context) error {
		valoresNecesarios := make(map[string]string)
		for _, requisito := range v.clavesNecesarias {
			valoresNecesarios[requisito] = ec.QueryParam(requisito)
		}

		data := make(DataView)
		for i, nombreValor := range nombreVariables {
			informacion := informaciones[i]
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
