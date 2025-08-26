package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
	"slices"
	"strings"
)

type FnInformacion func(*b.Bdd, []string) (any, error)

func deepCopyDatos(datos d.ConjuntoDato) d.ConjuntoDato {
	return datos
}

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

	resultado := make(d.ConjuntoDato)
	posicionResultados := make([]*d.ConjuntoDato, len(infoVariables))
	for i, info := range infoVariables {
		posicionResultados[i] = &resultado

		for _, path := range info.Path {
			if dato, ok := (*posicionResultados[i])[path]; !ok {
				nuevoConjunto := make(d.ConjuntoDato)
				ptrNuevoConjunto := &nuevoConjunto

				(*posicionResultados[i])[path] = *ptrNuevoConjunto
				posicionResultados[i] = ptrNuevoConjunto

			} else if lugar, ok := dato.(d.ConjuntoDato); !ok {
				return nil, fmt.Errorf("se construyó mal el conjunto de datos para devolver al usuario")

			} else {
				posicionResultados[i] = &lugar
			}
		}
	}

	return func(bdd *b.Bdd, parametrosDados []string) (any, error) {
		fila := bdd.MySQL.QueryRow(query.Select, parametrosDados[indiceIdNecesario])

		if err := fila.Scan(datosReferencias...); err != nil {
			return nil, fmt.Errorf("en lectura de una fila, con la query %s, se tuvo el error: %v", query.Select, err)
		}

		for i, info := range infoVariables {
			if valor, err := info.Variable.Desreferenciar(datosReferencias[i]); err != nil {
				return nil, fmt.Errorf("no se pudo desreferenciar %s.%s por: %v", strings.Join(info.Path, "."), info.Nombre, err)

			} else {
				subresultado := posicionResultados[i]
				(*subresultado)[info.Nombre] = valor
			}
		}

		// No se necesita usar deepCopy porque solo se usa un elemento a la vez, por lo que
		//   realmente se puede pisar los valores viejos y no deberia haber problema, ya que se crea
		//   una de estas funciones por cada variable, por lo que no se reutiliza en otro lado
		return resultado, nil
	}, nil
}
