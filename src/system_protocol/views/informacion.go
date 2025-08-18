package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

type Informacion interface {
	ObtenerInformacion(bdd *b.Bdd, requisitos map[string]string) (any, error)
}

// Equivale al ParametroElemntos
type InformacionTabla struct {
	Tabla       *d.DescripcionTabla
	Condicion   d.Condicion
	Referencias map[string]InformacionReferencia
}

func NewInformacionTabla(tabla *d.DescripcionTabla, condicion d.Condicion, referencias map[string]InformacionReferencia) InformacionTabla {
	return InformacionTabla{
		Tabla:       tabla,
		Condicion:   condicion,
		Referencias: referencias,
	}
}

type InformacionFila struct {
	Tabla     *d.DescripcionTabla
	Condicion d.Condicion
}

func NewInformacionFila(tabla *d.DescripcionTabla, condicion d.Condicion) InformacionFila {
	return InformacionFila{
		Tabla:     tabla,
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

type InformacionArray struct {
	Elementos []Informacion
}

func (i InformacionTabla) ObtenerInformacion(bdd *b.Bdd, requisitos map[string]string) (any, error) {
	datos := []d.ConjuntoDato{}

	iterador, err := i.Tabla.QueryAll(bdd, i.Condicion, requisitos)
	if err != nil {
		return nil, fmt.Errorf("se tuvo un error al intentar iterar sobre la tabla: %s, con error: %v", i.Tabla.NombreTabla, err)
	}

	for conjuntoDato := range iterador {
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
					return nil, fmt.Errorf("se necesita valor en %s de la tabla %s, y no se consiguio", claveValor, i.Tabla.NombreTabla)
				}
			}

			conjuntoDato[nombre] = pathView
		}

		datos = append(datos, conjuntoDato)
	}

	return datos, nil
}

func (i InformacionFila) ObtenerInformacion(bdd *b.Bdd, requisitos map[string]string) (any, error) {
	return i.Tabla.QueryElemento(bdd, i.Condicion, requisitos)
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

func (i InformacionArray) ObtenerInformacion(bdd *b.Bdd, requisitos map[string]string) (any, error) {
	datos := make([]any, len(i.Elementos))
	for i, informacion := range i.Elementos {
		if dato, err := informacion.ObtenerInformacion(bdd, requisitos); err != nil {
			return datos, fmt.Errorf("en el elemento %d, se tuvo el error: %v", i, err)

		} else {
			datos[i] = dato
		}
	}
	return datos, nil
}
