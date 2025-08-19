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
	fmt.Printf("Registrando %s...\n", nombre)
	return View{
		Nombre:           nombre,
		Bloque:           bloque,
		Bdd:              bdd,
		clavesNecesarias: requisitos,
		informaciones:    informaciones,
	}
}

func (v View) GenerarEndpoint(ec echo.Context) error {
	valoresNecesarios := make(map[string]string)
	for _, requisito := range v.clavesNecesarias {
		valoresNecesarios[requisito] = ec.QueryParam(requisito)
	}

	data := make(DataView)
	for nombreValor := range v.informaciones {
		informacion := v.informaciones[nombreValor]
		if valor, err := informacion.ObtenerInformacion(v.Bdd, valoresNecesarios); err != nil {
			fmt.Printf("Error al utilizar endpoint /%s, dado la informacion %s con error: %v\n", v.Nombre, nombreValor, err)
			return err

		} else {
			data[nombreValor] = valor
		}
	}

	return ec.Render(200, v.Bloque, data)
}
