package dependencias

import (
	"fmt"
	"slices"
	"strings"
)

type QueryDato struct {
	Select string
	Claves []*HojaClave
}

/*
	"Carreras": {
		"tipo": "ElementoUnico",
		"tabla": "Carreras",
		"id": "idCarrera",
		"claves": [ "nombre" ]
	}

Resulta en la query

	"SELECT nombre FROM Carreras WHERE id = ?"

Para esto, usamos la variable de "nombre", y se tiene que buscar el indice del parametro "idCarrera" para
pasarlo en el lugar de ?

Otro ejemplo:

	"Materia": {
		"tipo": "ElementoUnico",
		"tabla": "Materias",
		"id": "idMateria",
		"claves": [ "id", "nombre", "refCarrera:id", "refCarrera:nombre" ]
	},

Resulta en la query

	"SELECT res.Materias_id, res.Materias_nombre, res.Carreras_id, res.Carreras_nombre FROM (
		SELECT Materias.id AS Materias_id, Materias.nombre AS Materias_nombre, Carreras.id AS Carreras_id, Carreras.nombre AS Carreras_nombre
		FROM Materias INNER JOIN Carreras ON Materias.refCarrera == Carreras.id
	) AS res WHERE res.Materias_id = ?"

Otro ejemplo:

	"Tema": {
	    "tipo": "ElementoUnico",
	    "tabla": "TemasMateria",
	    "id": "idTema",
	    "claves": [ "nombre", "refMateria:id", "refMateria:nombre", "refMateria:refCarrera:id", "refMateria:refCarrera:nombre" ]
	},

# Resulta en la query

	SELECT TemasMateria.id AS TemasMateria_id, TemasMateria.nombre AS TemasMateria_nombre, temp1_1.* FROM TemasMateria
	JOIN (
		SELECT Materias.id AS Materias_id, Materias.nombre AS Materias_nombre, temp2_1.* FROM Materias
		JOIN (
			SELECT Carreras.id AS Carreras_id, Carreras.nombre AS Carreras_nombre FROM Carreras
		)
		AS temp2_1 ON Materias.refCarrera = temp2_1.Carreras_id
		JOIN (
			SELECT Carreras.id AS Carreras_id, Carreras.nombre AS Carreras_nombre FROM Carreras
		)
		AS temp2_2 ON Materias.refCarrera = temp2_1.Carreras_id
	)
	AS temp1_1 ON TemasMateria.refMateria = temp1_1.Materias_id
	WHERE TemasMateria_id = 1;
*/
func generarSelect(nodo *NodoClave, profundidad int) string {
	nombreTabla := nodo.Tabla.NombreTabla
	claves := make([]string, len(nodo.Claves))
	for i, clave := range nodo.Claves {
		nombreClave := clave.Nombre
		claves[i] = fmt.Sprintf("%s.%s AS %s_%s", nombreTabla, nombreClave, nombreTabla, nombreClave)
	}

	if len(nodo.Referencias) == 0 {
		return fmt.Sprintf("SELECT %s FROM %s", strings.Join(claves, ", "), nombreTabla)
	}

	sentenciasJoin := make([]string, len(nodo.Referencias))
	for i, referencia := range nodo.Referencias {
		sentenciaInterna := generarSelect(referencia, profundidad+1)

		sentenciasJoin[i] = fmt.Sprintf(
			"INNER JOIN (\n\t%s\n) AS temp_%d_%d ON %s.%s  temp_%d_%d.%s_id",
			sentenciaInterna,
			profundidad, i,
			nodo.Tabla.NombreTabla,
			referencia.Nombre,
			profundidad, i,
			referencia.Tabla.NombreTabla,
		)
	}

	return fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(claves, ", "),
		nombreTabla,
		strings.Join(sentenciasJoin, "\n"),
	)
}

func NewQuerySimple(tabla *DescripcionTabla, clavesUsadas []string) (QueryDato, error) {
	indiceId := slices.Index(clavesUsadas, "id")
	if indiceId < 0 {
		indiceId = len(clavesUsadas)
		clavesUsadas = append(clavesUsadas, "id")
	}

	claves := make([]*HojaClave, len(clavesUsadas))

	raiz := NewRaizClave(tabla)
	for i, clave := range clavesUsadas {
		if hoja, err := raiz.Insertar(clave); err != nil {
			return QueryDato{}, fmt.Errorf("no se pudo construir arbol de claves porque %v", err)

		} else {
			claves[i] = hoja
		}
	}

	return QueryDato{
		Select: fmt.Sprintf(
			"%s WHERE %s = ?",
			generarSelect(&raiz, 0),
			claves[indiceId].NombreQuery(),
		),
		Claves: claves,
	}, nil
}

func NewQueryMultiplesCompleto() (QueryDato, error) {
	return QueryDato{}, nil
}

func NewQueryMultiplesParcial() (QueryDato, error) {
	return QueryDato{}, nil
}

/*
type ElementoInformacion struct {
	Tipo           TipoVariable
	Representacion string
}

func NewQueryMultiples(nombreTabla string, clavesRequeridas map[string]TipoVariable, parClaveRepresentacion map[string]ElementoInformacion, clavesOrdenar []string) (FnMultiplesDatos, error) {
	expresionesClaves, representaciones, tiposWhere := informacionWhere(parClaveRepresentacion)
	clavesOrdenadas, datosReferencias, tiposDatos, err := informacionClaves(clavesRequeridas)
	if err != nil {
		return nil, err
	}

	expresionSort := ""
	if len(clavesOrdenar) > 0 {
		expresionSort = fmt.Sprintf("ORDER BY %s", strings.Join(clavesOrdenar, ", "))
	}

	dereferenciarDatos := generarDesreferenciar(tiposDatos, clavesOrdenadas)

	if len(expresionesClaves) > 0 {
		query := strings.TrimSpace(fmt.Sprintf(
			"SELECT %s FROM %s WHERE %s %s",
			strings.Join(clavesOrdenadas, ", "),
			nombreTabla,
			strings.Join(expresionesClaves, " AND "),
			expresionSort,
		))

		valoresWhere := generarValoresWhere(tiposWhere, representaciones)

		return func(bdd *b.Bdd, datosRepresentativos map[string]string) ([]ConjuntoDato, error) {
			conjuntosDeDatos := []ConjuntoDato{}
			datosWhere, err := valoresWhere(datosRepresentativos)
			if err != nil {
				return conjuntosDeDatos, err
			}

			rows, err := bdd.MySQL.Query(query, datosWhere...)
			if err != nil {
				return conjuntosDeDatos, fmt.Errorf("en query mutliple, al hacer query con where se tuvo: %v", err)
			}

			defer rows.Close()
			for rows.Next() {
				if err := rows.Scan(datosReferencias...); err != nil {
					return conjuntosDeDatos, err
				}

				if conjuntoDato, err := dereferenciarDatos(datosReferencias); err != nil {
					return conjuntosDeDatos, err
				} else {
					conjuntosDeDatos = append(conjuntosDeDatos, conjuntoDato)
				}
			}

			return conjuntosDeDatos, nil
		}, nil

	} else {
		query := strings.TrimSpace(fmt.Sprintf(
			"SELECT %s FROM %s %s",
			strings.Join(clavesOrdenadas, ", "),
			nombreTabla,
			expresionSort,
		))

		return func(bdd *b.Bdd, datosRepresentativos map[string]string) ([]ConjuntoDato, error) {
			conjuntosDeDatos := []ConjuntoDato{}
			rows, err := bdd.MySQL.Query(query)
			if err != nil {
				return conjuntosDeDatos, err
			}

			defer rows.Close()
			for rows.Next() {
				if err := rows.Scan(datosReferencias...); err != nil {
					return conjuntosDeDatos, fmt.Errorf("en query mutliple, al hacer query sin where se tuvo: %v", err)
				}

				if conjuntoDato, err := dereferenciarDatos(datosReferencias); err != nil {
					return conjuntosDeDatos, err
				} else {
					conjuntosDeDatos = append(conjuntosDeDatos, conjuntoDato)
				}
			}

			return conjuntosDeDatos, nil
		}, nil
	}
}

func NewQueryFila(nombreTabla string, clavesRequeridas map[string]TipoVariable, parClaveRepresentacion map[string]ElementoInformacion, clavesOrdenar []string) (FnUnDato, error) {
	expresionesClaves, representaciones, tiposWhere := informacionWhere(parClaveRepresentacion)
	clavesOrdenadas, datosReferencias, tiposDatos, err := informacionClaves(clavesRequeridas)
	if err != nil {
		return nil, err
	}

	expresionSort := ""
	if len(clavesOrdenar) > 0 {
		expresionSort = fmt.Sprintf("ORDER BY %s", strings.Join(clavesOrdenar, ", "))
	}

	dereferenciarDatos := generarDesreferenciar(tiposDatos, clavesOrdenadas)
	if len(expresionesClaves) == 0 {
		return nil, fmt.Errorf("no se definió ninguna condicion de unicidad")
	}

	query := strings.TrimSpace(fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s %s",
		strings.Join(clavesOrdenadas, ", "),
		nombreTabla,
		strings.Join(expresionesClaves, " AND "),
		expresionSort,
	))

	valoresWhere := generarValoresWhere(tiposWhere, representaciones)

	return func(bdd *b.Bdd, datosRepresentativos map[string]string) (ConjuntoDato, error) {
		conjuntoDatos := make(ConjuntoDato)
		datosWhere, err := valoresWhere(datosRepresentativos)
		if err != nil {
			return conjuntoDatos, fmt.Errorf("en query fila, se tuvo: %v", err)
		}

		row := bdd.MySQL.QueryRow(query, datosWhere...)
		if err := row.Scan(datosReferencias...); err != nil {
			return conjuntoDatos, err
		}
		return dereferenciarDatos(datosReferencias)
	}, nil
}

func generarDesreferenciar(tiposDatos []TipoVariable, clavesOrdenadas []string) func(datos []any) (ConjuntoDato, error) {
	return func(datos []any) (ConjuntoDato, error) {
		conjuntoDato := make(ConjuntoDato)
		for i, tipo := range tiposDatos {
			clave := clavesOrdenadas[i]
			if valorDesreferenciado, err := tipo.Desreferenciar(datos[i]); err != nil {
				return conjuntoDato, err

			} else {
				conjuntoDato[clave] = valorDesreferenciado
			}
		}
		return conjuntoDato, nil
	}
}

func generarValoresWhere(tiposWhere []TipoVariable, representaciones []string) func(datos map[string]string) ([]any, error) {
	largoTipos := len(tiposWhere)

	return func(datos map[string]string) ([]any, error) {
		valoresReales := make([]any, largoTipos)
		for i, tipo := range tiposWhere {
			claveRepresentacion := representaciones[i]

			if representacion, ok := datos[claveRepresentacion]; !ok {
				return valoresReales, fmt.Errorf("no se paso el dato para la clave %s", claveRepresentacion)

			} else if valorReal, err := tipo.ValorDeString(representacion); err != nil {
				return valoresReales, err

			} else Errorf()
				valoresReales[i] = valorReal
			}
		}
		return valoresReales, nil
	}
}

func informacionWhere(parClaveRepresentacion map[string]ElementoInformacion) ([]string, []string, []TipoVariable) {
	expresionesClaves := make([]string, len(parClaveRepresentacion))
	representaciones := make([]string, len(parClaveRepresentacion))
	tiposWhere := make([]TipoVariable, len(parClaveRepresentacion))
	contador := 0
	for clave := range parClaveRepresentacion {
		elemento := parClaveRepresentacion[clave]

		expresionesClaves[contador] = fmt.Sprintf("%s = ?", clave)
		representaciones[contador] = elemento.Representacion
		tiposWhere[contador] = elemento.Tipo

		contador++
	}

	return expresionesClaves, representaciones, tiposWhere
}

func informacionClaves(clavesRequeridas map[string]TipoVariable) ([]string, []any, []TipoVariable, error) {
	clavesOrdenadas := make([]string, len(clavesRequeridas))
	datosReferencias := make([]any, len(clavesRequeridas))
	tiposDatos := make([]TipoVariable, len(clavesRequeridas))

	contador := 0
	for clave := range clavesRequeridas {
		tipo := clavesRequeridas[clave]

		clavesOrdenadas[contador] = clave
		tiposDatos[contador] = tipo
		if referencia, err := tipo.ReferenciaValor(); err != nil {
			return clavesOrdenadas, datosReferencias, tiposDatos, fmt.Errorf("no se pudo hacer la referencia para la clave %s con: %v", clave, err)

		} else {
			datosReferencias[contador] = referencia
		}

		contador++
	}

	return clavesOrdenadas, datosReferencias, tiposDatos, nil
}

*/
