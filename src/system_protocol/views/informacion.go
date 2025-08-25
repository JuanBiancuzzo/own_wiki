package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
	"slices"
	"strings"
)

type FnInformacion func(*b.Bdd, []string) (any, error)

func NewInformacionFila(query d.QueryDato, parametrosEsperados []string, pametroParaId string) (FnInformacion, error) {
	indiceIdNecesario := slices.Index(parametrosEsperados, pametroParaId)

	datosReferencias := make([]any, len(query.Claves))
	infoVariables := make([]d.InformacionClave, len(query.Claves))

	for i, clave := range query.Claves {
		infoVariables[i] = clave.ObtenerInfoVariable()
		if referencia, err := infoVariables[i].Variable.ObtenerReferencia(); err != nil {
			return nil, err

		} else {
			datosReferencias[i] = referencia
		}
	}

	return func(bdd *b.Bdd, parametrosDados []string) (any, error) {
		fila := bdd.MySQL.QueryRow(query.Select, parametrosDados[indiceIdNecesario])

		if err := fila.Scan(datosReferencias...); err != nil {
			return nil, fmt.Errorf("en lectura de una fila, con la query %s, se tuvo el error: %v", query.Select, err)
		}

		resultado := make(d.ConjuntoDato)
		for i, info := range infoVariables {
			subresultado := &resultado

			for _, path := range info.Path {
				if dato, ok := (*subresultado)[path]; !ok {
					nuevoConjunto := make(d.ConjuntoDato)
					ptrNuevoConjunto := &nuevoConjunto

					(*subresultado)[path] = *ptrNuevoConjunto
					subresultado = ptrNuevoConjunto

				} else if lugar, ok := dato.(d.ConjuntoDato); !ok {
					return nil, fmt.Errorf("se construyó mal el conjunto de datos para devolver al usuario")

				} else {
					subresultado = &lugar
				}
			}

			if valor, err := info.Variable.Desreferenciar(datosReferencias[i]); err != nil {
				return nil, fmt.Errorf("no se pudo desreferenciar %s.%s por: %v", strings.Join(info.Path, "."), info.Nombre, err)

			} else {
				(*subresultado)[info.Nombre] = valor
			}
		}

		return resultado, nil
	}, nil
}

/*
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
		datos[idx] = conjuntoDato
	}

	return datos, nil
}

*/
