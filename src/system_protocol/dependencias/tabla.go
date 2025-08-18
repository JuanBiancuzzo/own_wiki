package dependencias

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	"reflect"
	"strings"
)

type TipoTabla byte

const (
	DEPENDIENTE_NO_DEPENDIBLE   = 0b00
	DEPENDIENTE_DEPENDIBLE      = 0b10
	INDEPENDIENTE_NO_DEPENDIBLE = 0b01
	INDEPENDIENTE_DEPENDIBLE    = 0b11
)

type DescripcionTabla struct {
	NombreTabla string
	TipoTabla   TipoTabla
	ClavesTipo  []ParClaveTipo
	Referencias []ReferenciaTabla

	necesarioQuery            bool
	clavesRepresentativas     []string
	clavesGenerales           []string
	clavesTotal               []string
	referenciaMultiplesTablas map[string]bool
	tablasReferencias         map[string]map[string]*DescripcionTabla
	tipoDadoClave             map[string]TipoVariable
}

func ConstruirTabla(nombreTabla string, tipoTabla TipoTabla, elementosRepetidos bool, clavesTipo []ParClaveTipo, referencias []ReferenciaTabla) DescripcionTabla {
	clavesGenerales := []string{}
	clavesRepresentativas := []string{}
	clavesTotal := []string{"id"}

	referenciaMultiplesTablas := make(map[string]bool)
	tipoDadoClave := make(map[string]TipoVariable)
	tipoDadoClave["id"] = TV_INT

	tablasReferencias := make(map[string]map[string]*DescripcionTabla)
	for _, claveTipo := range clavesTipo {
		clavesGenerales = append(clavesGenerales, claveTipo.Clave)
		clavesTotal = append(clavesTotal, claveTipo.Clave)
		if claveTipo.Representativa {
			clavesRepresentativas = append(clavesRepresentativas, claveTipo.Clave)
		}
		tipoDadoClave[claveTipo.Clave] = claveTipo.tipo
	}
	for _, referencia := range referencias {
		clavesGenerales = append(clavesGenerales, referencia.Clave)
		referenciaMultiplesTablas[referencia.Clave] = len(referencia.Tablas) > 1

		if referencia.Representativo {
			clavesRepresentativas = append(clavesRepresentativas, referencia.Clave)
		}

		mapaTablas := make(map[string]*DescripcionTabla)
		for _, tabla := range referencia.Tablas {
			mapaTablas[tabla.NombreTabla] = tabla
		}
		tablasReferencias[referencia.Clave] = mapaTablas
		tipoDadoClave[referencia.Clave] = TV_REFERENCIA
		if len(referencia.Tablas) > 1 {
			tipoDadoClave[fmt.Sprintf("tipo%s", referencia.Clave)] = TV_ENUM
			clavesTotal = append(clavesTotal, fmt.Sprintf("tipo%s", referencia.Clave))
		}
		clavesTotal = append(clavesTotal, referencia.Clave)
	}

	return DescripcionTabla{
		NombreTabla: nombreTabla,
		TipoTabla:   tipoTabla,
		ClavesTipo:  clavesTipo,
		Referencias: referencias,

		necesarioQuery:            elementosRepetidos,
		clavesGenerales:           clavesGenerales,
		clavesRepresentativas:     clavesRepresentativas,
		referenciaMultiplesTablas: referenciaMultiplesTablas,
		tipoDadoClave:             tipoDadoClave,
		clavesTotal:               clavesTotal,

		tablasReferencias: tablasReferencias,
	}
}

func (dt DescripcionTabla) CrearTablaRelajada(bdd *b.Bdd) error {
	parametros := []string{}
	for _, parClaveTipo := range dt.ClavesTipo {
		parametros = append(parametros, fmt.Sprintf("%s %s", parClaveTipo.Clave, parClaveTipo.TipoSQL))
	}

	for _, referencia := range dt.Referencias {
		if len(referencia.Tablas) > 1 {
			tablasReferenciadas := make([]string, len(referencia.Tablas))
			for i, tablaReferenciada := range referencia.Tablas {
				tablasReferenciadas[i] = tablaReferenciada.NombreTabla
			}
			parClaveTipo := NewClaveEnum(false, fmt.Sprintf("tipo%s", referencia.Clave), tablasReferenciadas, false)
			parametros = append(parametros, fmt.Sprintf("%s %s", parClaveTipo.Clave, parClaveTipo.TipoSQL))
		}

		parametros = append(parametros, fmt.Sprintf("%s INT", referencia.Clave))
	}

	tabla := fmt.Sprintf(
		"CREATE TABLE %s (\n\tid INT AUTO_INCREMENT PRIMARY KEY,\n\t%s\n);",
		dt.NombreTabla,
		strings.Join(parametros, ",\n\t"),
	)

	if err := bdd.CrearTabla(tabla); err != nil {
		return fmt.Errorf("no se pudo crear la tabla \n%s\n, con error: %v", tabla, err)
	}
	return nil
}

func (dt DescripcionTabla) Existe(bdd *b.Bdd, datosIngresados ConjuntoDato) (bool, error) {
	if !dt.necesarioQuery {
		return false, nil
	}

	datos := []any{}
	queryParam := make([]string, len(dt.clavesRepresentativas))
	for i, clave := range dt.clavesRepresentativas {
		if dato, ok := datosIngresados[clave]; !ok {
			return false, fmt.Errorf("el usuario no ingreso el dato para %s", clave)

		} else if relacion, ok := dato.(RelacionTabla); ok {
			if dt.referenciaMultiplesTablas[clave] {
				datos = append(datos, relacion.Tabla)
			}

		} else {
			datos = append(datos, dato)
		}

		queryParam[i] = fmt.Sprintf("%s = ?", clave)
	}

	query := fmt.Sprintf(
		"SELECT id FROM %s WHERE %s",
		dt.NombreTabla,
		strings.Join(queryParam, " AND "),
	)

	_, err := bdd.Obtener(query, datos...)
	return err == nil, nil
}

func (dt DescripcionTabla) Insertar(bdd *b.Bdd, datosIngresados ConjuntoDato) (int64, error) {
	datos := []any{}

	insertarParam := []string{}
	for _, clave := range dt.clavesGenerales {
		if dato, ok := datosIngresados[clave]; !ok {
			return 0, fmt.Errorf("el usuario no ingreso el dato para %s", clave)

		} else if relacion, ok := dato.(RelacionTabla); ok {
			if dt.referenciaMultiplesTablas[clave] {
				insertarParam = append(insertarParam, fmt.Sprintf("tipo%s", clave))
				datos = append(datos, relacion.Tabla)
			}
			insertarParam = append(insertarParam, clave)
			datos = append(datos, 0)
		} else {
			insertarParam = append(insertarParam, clave)
			datos = append(datos, dato)
		}
	}

	insertar := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		dt.NombreTabla,
		strings.Join(insertarParam, ", "),
		strings.Join(strings.Split(strings.Repeat("?", len(insertarParam)), ""), ", "),
	)
	return bdd.Insertar(insertar, datos...)
}

func (dt DescripcionTabla) QueryAll(bdd *b.Bdd, condicion Condicion, datosCondicion map[string]string) (func(yield func(ConjuntoDato) bool), error) {
	funcionNula := func(yield func(ConjuntoDato) bool) {}

	expresionCondicion, datos, err := condicion.Expresion(datosCondicion, dt.tipoDadoClave)
	if err != nil {
		return funcionNula, fmt.Errorf("se tuvo un error: %v en la expresion", err)
	}
	query := fmt.Sprintf(
		"SELECT * FROM %s %s",
		dt.NombreTabla,
		expresionCondicion,
	)

	if rows, err := bdd.MySQL.Query(query, datos...); err != nil {
		return funcionNula, fmt.Errorf("se tuvo un error: %v la hacer una query dada por %s y datos %+v", err, query, datos)

	} else {
		return func(yield func(ConjuntoDato) bool) {
			defer rows.Close()
			for rows.Next() {
				datosTabla := make([]any, len(dt.clavesTotal))
				for i, clave := range dt.clavesTotal {
					if referencia, err := dt.tipoDadoClave[clave].ReferenciaValor(); err != nil {
						fmt.Printf("Dejando de iterar, hay un error en el valor de referencia: %v\n", err)
						return

					} else {
						datosTabla[i] = referencia
					}
				}

				if err := rows.Scan(datosTabla...); err != nil {
					fmt.Printf("Dejando de iterar, hay un error en el scan %v\n", err)
					return
				}

				conjuntoDato := make(ConjuntoDato)
				for i, clave := range dt.clavesTotal {
					valor := reflect.ValueOf(datosTabla[i])
					if valor.Kind() == reflect.Ptr {
						conjuntoDato[clave] = valor.Elem()

					} else {
						fmt.Println("Por alguna razon no es un puntero")
						return
					}
				}

				if !yield(conjuntoDato) {
					return
				}
			}
		}, nil
	}
}

func (dt DescripcionTabla) QueryElemento(bdd *b.Bdd, condicion Condicion, datosCondicion map[string]string) (ConjuntoDato, error) {
	expresionCondicion, datos, err := condicion.Expresion(datosCondicion, dt.tipoDadoClave)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(
		"SELECT * FROM %s %s",
		dt.NombreTabla,
		expresionCondicion,
	)

	row := bdd.MySQL.QueryRow(query, datos...)

	datosTabla := make([]any, len(dt.clavesTotal))
	for i, clave := range dt.clavesTotal {
		if referencia, err := dt.tipoDadoClave[clave].ReferenciaValor(); err != nil {
			return nil, fmt.Errorf("dejando de iterar, hay un error en el valor de referencia: %v", err)

		} else {
			datosTabla[i] = referencia
		}
	}

	if err := row.Scan(datosTabla...); err != nil {
		return nil, fmt.Errorf("dejando de iterar, hay un error en el scan %v", err)
	}

	conjuntoDato := make(ConjuntoDato)
	for i, clave := range dt.clavesTotal {
		valor := reflect.ValueOf(datosTabla[i])
		if valor.Kind() == reflect.Ptr {
			conjuntoDato[clave] = valor.Elem()

		} else {
			return nil, fmt.Errorf("por alguna razon no es un puntero")
		}
	}

	return conjuntoDato, nil
}

// TODO
func (dt DescripcionTabla) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (dt DescripcionTabla) CrearForeignKey(hash *Hash, datosIngresados ConjuntoDato) ([]ForeignKey, error) {
	fKeys := []ForeignKey{}

	for clave := range datosIngresados {
		if relacion, ok := datosIngresados[clave].(RelacionTabla); !ok {
			continue

		} else if mapaTablas, ok := dt.tablasReferencias[clave]; !ok {
			return fKeys, fmt.Errorf("no hay tablas con la clave %s", clave)

		} else if tabla, ok := mapaTablas[relacion.Tabla]; !ok {
			return fKeys, fmt.Errorf("no hay tabla (%s) para esa relacion", relacion.Tabla)

		} else if hash, err := tabla.Hash(hash, relacion.Datos); err != nil {
			return fKeys, fmt.Errorf("no se pudo general el hash de los datos, con err: %v", err)

		} else {
			fKeys = append(fKeys, NewForeignKey(relacion.Tabla, clave, hash))
		}
	}

	return fKeys, nil
}

func (dt DescripcionTabla) Hash(hash *Hash, datosIngresados ConjuntoDato) (IntFK, error) {
	datosRepresentativos := []any{}

	for _, clave := range dt.clavesRepresentativas {
		if dato, ok := datosIngresados[clave]; !ok {
			return 0, fmt.Errorf("no se ingreso el valor para la clave %s", clave)

		} else if relacion, ok := dato.(RelacionTabla); ok {
			if mapaTablas, ok := dt.tablasReferencias[clave]; !ok {
				return 0, fmt.Errorf("no hay tablas con la clave %s", clave)

			} else if tabla, ok := mapaTablas[relacion.Tabla]; !ok {
				return 0, fmt.Errorf("no existe relación de %s a %s con la clave %s", dt.NombreTabla, relacion.Tabla, clave)

			} else if hash, err := tabla.Hash(hash, relacion.Datos); err != nil {
				return 0, err

			} else {
				datosRepresentativos = append(datosRepresentativos, hash)
			}

		} else {
			datosRepresentativos = append(datosRepresentativos, dato)
		}
	}

	return hash.HasearDatos(datosRepresentativos...), nil
}

func (dt DescripcionTabla) ObtenerDependencias() []DescripcionTabla {
	tablas := []DescripcionTabla{}

	for _, referencia := range dt.Referencias {
		for _, tabla := range referencia.Tablas {
			tablas = append(tablas, *tabla)
		}
	}

	return tablas
}
func EsTipoDependiente(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == DEPENDIENTE_NO_DEPENDIBLE
}

func EsTipoDependible(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == INDEPENDIENTE_DEPENDIBLE
}
