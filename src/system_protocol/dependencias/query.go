package dependencias

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	"strings"
)

type FnMultiplesDatos func(bdd *b.Bdd, datosRepresentativos map[string]string) ([]ConjuntoDato, error)
type FnUnDato func(bdd *b.Bdd, datosRepresentativos map[string]string) (ConjuntoDato, error)

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

			} else {
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
