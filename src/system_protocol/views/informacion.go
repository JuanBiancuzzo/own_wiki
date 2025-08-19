package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

type Informacion interface {
	ObtenerInformacion(bdd *b.Bdd, requisitos map[string]string) (any, error)

	NecesitaEndpoint() (bool, string)
}

// Equivale al ParametroElemntos
//
//	-> Hacer una informacion que sea para ir mandando poco a poco
type InformacionTabla struct {
	Ruta        string
	Query       d.FnMultiplesDatos
	Referencias map[string]InformacionReferencia
}

func NewInformacionTabla(ruta string, query d.FnMultiplesDatos, referencias map[string]InformacionReferencia) InformacionTabla {
	return InformacionTabla{
		Ruta:        ruta,
		Query:       query,
		Referencias: referencias,
	}
}

type InformacionFila struct {
	Condicion d.FnUnDato
}

func NewInformacionFila(condicion d.FnUnDato) InformacionFila {
	return InformacionFila{
		Condicion: condicion,
	}
}

type InformacionReferencia struct {
	View       string
	Parametros map[string]string
}

func NewInformacionReferencia(view string, parametros map[string]string) InformacionReferencia {
	return InformacionReferencia{
		View:       view,
		Parametros: parametros,
	}
}

func (i InformacionTabla) ObtenerInformacion(bdd *b.Bdd, requisitos map[string]string) (any, error) {
	listaConjuntos, err := i.Query(bdd, requisitos)
	if err != nil {
		return nil, fmt.Errorf("se tuvo un error obtener valores, con error: %v", err)
	}

	datos := make([]d.ConjuntoDato, len(listaConjuntos))
	for idx, conjuntoDato := range listaConjuntos {
		for nombre := range i.Referencias {
			referencia := i.Referencias[nombre]
			pathView := NewPathView(referencia.View)

			for parametro := range referencia.Parametros {
				claveValor := referencia.Parametros[parametro]

				if valor, ok := conjuntoDato[claveValor]; ok {
					pathView.AgregarParametro(parametro, valor)

				} else if valor, ok := requisitos[claveValor]; ok {
					pathView.AgregarParametro(parametro, valor)

				} else {
					return nil, fmt.Errorf("se necesita valor en %s, y no se consiguio", claveValor)
				}
			}

			conjuntoDato[nombre] = pathView
		}

		datos[idx] = conjuntoDato
	}

	return datos, nil
}

func (i InformacionTabla) NecesitaEndpoint() (bool, string) {
	return true, i.Ruta
}

func (i InformacionFila) ObtenerInformacion(bdd *b.Bdd, requisitos map[string]string) (any, error) {
	return i.Condicion(bdd, requisitos)
}

func (i InformacionFila) NecesitaEndpoint() (bool, string) {
	return false, ""
}

func (i InformacionReferencia) ObtenerInformacion(bdd *b.Bdd, requisitos map[string]string) (any, error) {
	pathView := NewPathView(i.View)

	for parametro := range i.Parametros {
		claveValor := i.Parametros[parametro]

		if valor, ok := requisitos[claveValor]; !ok {
			return nil, fmt.Errorf("se necesita valor en %s, y no se consiguio", claveValor)

		} else {
			pathView.AgregarParametro(claveValor, valor)
		}
	}

	return pathView, nil
}

func (i InformacionReferencia) NecesitaEndpoint() (bool, string) {
	return false, ""
}
