package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
	"slices"
	"strings"
)

type FnInformacion func(*b.Bdd, []string) (any, error)

type RespuestaInformacion struct {
	Informacion   FnInformacion
	ExtraEndpoint map[string]Endpoint
}

func deepCopyDatos(datos d.ConjuntoDato) d.ConjuntoDato {
	copia := make(d.ConjuntoDato)
	for clave := range datos {
		dato := datos[clave]
		copia[clave] = dato

		if subDatos, ok := dato.(d.ConjuntoDato); ok {
			copia[clave] = deepCopyDatos(subDatos)
		}
	}

	return copia
}

func crearDatosSegunEstructura(infoVariables []d.InformacionClave) (d.ConjuntoDato, []*d.ConjuntoDato, error) {
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
				return resultado, posicionResultados, fmt.Errorf("se construyó mal el conjunto de datos para devolver al usuario")

			} else {
				posicionResultados[i] = &lugar
			}
		}
	}

	return resultado, posicionResultados, nil
}

func NewInformacionFila(query d.QueryDato, parametrosEsperados []string) (FnInformacion, error) {
	indiceIdNecesario := slices.Index(parametrosEsperados, query.Parametros[0])

	datosReferencias := make([]any, len(query.ClaveSelect))
	infoVariables := make([]d.InformacionClave, len(query.ClaveSelect))

	for i, clave := range query.ClaveSelect {
		infoVariables[i] = clave.ObtenerInfoVariable()
		if referencia, err := infoVariables[i].Variable.ObtenerReferencia(); err != nil {
			return nil, err

		} else {
			datosReferencias[i] = referencia
		}
	}

	resultado, posicionResultados, err := crearDatosSegunEstructura(infoVariables)
	if err != nil {
		return nil, err
	}

	return func(bdd *b.Bdd, parametrosDados []string) (any, error) {
		fila := bdd.QueryRow(query.SentenciaQuery, parametrosDados[indiceIdNecesario])

		if err := fila.Scan(datosReferencias...); err != nil {
			return nil, fmt.Errorf("en lectura de una fila, con la query %s, se tuvo el error: %v", query.SentenciaQuery, err)
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

type FnInformacionIterador func(*b.Bdd, []string) (FnIteradorDato, error)
type FnIteradorDato func(yield func(d.ConjuntoDato) bool)

func crearFuncionGenerador(nombreTabla string, query d.QueryDato, parametrosEsperados []string) (FnInformacionIterador, error) {
	datosReferencias := make([]any, len(query.ClaveSelect))
	infoVariables := make([]d.InformacionClave, len(query.ClaveSelect))

	fmt.Println("---")
	for i, clave := range query.ClaveSelect {
		infoVariables[i] = clave.ObtenerInfoVariable()
		if referencia, err := infoVariables[i].Variable.ObtenerReferencia(); err != nil {
			return nil, err

		} else {
			datosReferencias[i] = referencia
		}
	}

	resultado, posicionResultados, err := crearDatosSegunEstructura(infoVariables)
	if err != nil {
		return nil, err
	}
	resultado["Tabla"] = nombreTabla

	indicesParametros := make([]int, len(query.Parametros))
	variablesRequeridas := make([]d.InformacionClave, len(query.Parametros))

	for i, parametro := range query.Parametros {
		if indice := slices.Index(parametrosEsperados, parametro); indice < 0 {
			return nil, fmt.Errorf("no se paso el parametro '%s'", parametro)

		} else {
			indicesParametros[i] = indice
		}

		claveWhere := query.ClaveWhere[i]
		variablesRequeridas[i] = claveWhere.ObtenerInfoVariable()
	}

	return func(bdd *b.Bdd, parametrosDados []string) (FnIteradorDato, error) {
		parametrosRequeridos := make([]any, len(indicesParametros))
		for i, indiceParametro := range indicesParametros {
			if valor, err := variablesRequeridas[i].Variable.ValorPorRepresentacion(parametrosDados[indiceParametro]); err != nil {
				return nil, fmt.Errorf("al obtener los parametros se obtuvo: %v", err)

			} else {
				parametrosRequeridos[i] = valor
			}
		}

		filas, err := bdd.Query(query.SentenciaQuery, parametrosRequeridos...)
		if err != nil {
			return nil, err
		}

		return func(yield func(d.ConjuntoDato) bool) {
			for filas.Next() {

				if err := filas.Scan(datosReferencias...); err != nil {
					fmt.Printf("en lectura de una fila, con la query %s, se tuvo el error: %v\n", query.SentenciaQuery, err)
					return
				}

				for i, info := range infoVariables {
					if valor, err := info.Variable.Desreferenciar(datosReferencias[i]); err != nil {
						fmt.Printf("no se pudo desreferenciar %s.%s por: %v\n", strings.Join(info.Path, "."), info.Nombre, err)
						return

					} else {
						subresultado := posicionResultados[i]
						(*subresultado)[info.Nombre] = valor
					}
				}

				if !yield(deepCopyDatos(resultado)) {
					filas.Close()
					return
				}
			}
			filas.Close()
		}, nil
	}, nil
}

func NewInformacionCompleta(querys map[string]d.QueryDato, parametrosEsperados []string) (FnInformacion, error) {
	tablas := make([]string, len(querys))
	posicionTabla := make(map[string]int)
	generadorIteradores := make([]FnInformacionIterador, len(querys))

	contador := 0
	for tabla := range querys {
		tablas[contador] = tabla
		posicionTabla[tabla] = contador
		if generador, err := crearFuncionGenerador(tabla, querys[tabla], parametrosEsperados); err != nil {
			return nil, err

		} else {
			generadorIteradores[contador] = generador
		}

		contador++
	}

	if len(querys) == 1 {
		generador := generadorIteradores[0]
		return func(bdd *b.Bdd, parametrosDados []string) (any, error) {
			iterador, err := generador(bdd, parametrosDados)
			if err != nil {
				return nil, err
			}

			datos := []d.ConjuntoDato{}
			for dato := range iterador {
				datos = append(datos, dato)
			}
			return datos, nil

		}, nil
	}

	iteradores := make([]FnIteradorDato, len(generadorIteradores))
	return func(bdd *b.Bdd, parametrosDados []string) (any, error) {
		var err error

		for i, generador := range generadorIteradores {
			if iteradores[i], err = generador(bdd, parametrosDados); err != nil {
				return nil, err
			}
		}

		// Cambiarlo para tener en cuenta lo del orden
		datos := []d.ConjuntoDato{}
		for _, iterador := range iteradores {
			for dato := range iterador {
				datos = append(datos, dato)
			}
		}
		return datos, nil
	}, nil
}

func NewInformacionParcial(querys map[string]d.QueryDato, parametrosEsperados []string) (RespuestaInformacion, error) {
	return RespuestaInformacion{
		Informacion:   func(b *b.Bdd, s []string) (any, error) { return []d.ConjuntoDato{}, nil },
		ExtraEndpoint: make(map[string]Endpoint),
	}, nil
}

/*
SELECT Materias.id AS Materias_id, Materias.nombre AS Materias_nombre FROM Materias
	INNER JOIN (
		SELECT Carreras.id AS Carreras_id FROM Carreras WHERE Carreras.id = ?
	) AS temp_0_0
	ON Materias.refCarrera = temp_0_0.Carreras_id
	INNER JOIN (
		SELECT CuatrimestresCarrera.anio AS CuatrimestresCarrera_anio, CuatrimestresCarrera.cuatrimestre AS CuatrimestresCarrera_cuatrimestre FROM CuatrimestresCarrera
	) AS temp_0_1
	ON Materias.refCuatrimestre = temp_0_1.CuatrimestresCarrera_id
*/
