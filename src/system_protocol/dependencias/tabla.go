package dependencias

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
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

	necesarioQuery    bool
	query             string
	insertar          string
	tablasReferencias map[string]*DescripcionTabla
}

func ConstruirTabla(nombreTabla string, tipoTabla TipoTabla, elementosRepetidos bool, clavesTipo []ParClaveTipo, referencias []ReferenciaTabla) DescripcionTabla {
	insertarParam := []string{}
	insertarValues := []string{}
	queryParam := []string{}
	tablasReferencias := make(map[string]*DescripcionTabla)
	for _, claveTipo := range clavesTipo {
		insertarParam = append(insertarParam, claveTipo.Clave)
		insertarValues = append(insertarValues, "?")
		queryParam = append(queryParam, fmt.Sprintf("%s = ?", claveTipo.Clave))
	}
	for _, referencia := range referencias {
		insertarParam = append(insertarParam, referencia.Clave)
		insertarValues = append(insertarValues, "0")
		for _, tabla := range referencia.Tablas {
			tablasReferencias[tabla.NombreTabla+referencia.Clave] = tabla
		}
	}

	return DescripcionTabla{
		NombreTabla: nombreTabla,
		TipoTabla:   tipoTabla,
		ClavesTipo:  clavesTipo,
		Referencias: referencias,

		necesarioQuery: elementosRepetidos,
		insertar: fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			nombreTabla,
			strings.Join(insertarParam, ", "),
			strings.Join(insertarValues, ", "),
		),
		query: fmt.Sprintf(
			"SELECT id FROM %s WHERE %s",
			nombreTabla,
			strings.Join(queryParam, " AND "),
		),
		tablasReferencias: tablasReferencias,
	}
}

func (dt DescripcionTabla) CrearTablaRelajada(bdd *b.Bdd) error {
	parametros := []string{}
	for i, parClaveTipo := range dt.ClavesTipo {
		extra := ","
		if i+1 == len(dt.ClavesTipo) && len(dt.Referencias) == 0 {
			extra = ""
		}
		parametros = append(parametros, fmt.Sprintf("%s %s%s", parClaveTipo.Clave, parClaveTipo.Tipo, extra))
	}

	for i, referencia := range dt.Referencias {
		extra := ","
		if i+1 == len(dt.Referencias) {
			extra = ""
		}
		parametros = append(parametros, fmt.Sprintf("%s INT%s", referencia.Clave, extra))
	}

	tabla := fmt.Sprintf(
		"CREATE TABLE %s (\nid INT AUTO_INCREMENT PRIMARY KEY,\n\t%s\n);",
		dt.NombreTabla,
		strings.Join(parametros, "\n\t"),
	)

	if err := bdd.CrearTabla(tabla); err != nil {
		return fmt.Errorf("no se pudo crear la tabla \n%s\n, con error: %v", tabla, err)
	}
	return nil
}

// TODO
func (dt DescripcionTabla) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (dt DescripcionTabla) CrearForeignKey(hash *Hash, relaciones []RelacionTabla) ([]ForeignKey, error) {
	fKeys := make([]ForeignKey, len(relaciones))

	for i, relacion := range relaciones {
		datos := relacion.Datos
		tabla, ok := dt.tablasReferencias[relacion.Tabla+relacion.Clave]
		if !ok {
			return fKeys, fmt.Errorf("no hay tabla para esa relacion")
		}

		if len(relacion.InfoRelacionada) > 0 {
			if fKeysRelacionados, err := tabla.CrearForeignKey(hash, relacion.InfoRelacionada); err != nil {
				return fKeys, fmt.Errorf("error info relacionada con error: %v", err)
			} else {
				for _, fKey := range fKeysRelacionados {
					datos = append(datos, fKey.HashDatosDestino)
				}
			}
		}

		fKeys[i] = NewForeignKey(hash, relacion.Tabla, relacion.Clave, datos...)
	}

	return fKeys, nil
}

func (dt DescripcionTabla) Hash(hash *Hash, fKeys []ForeignKey, datos ...any) (IntFK, error) {
	if len(datos) != len(dt.ClavesTipo) {
		return 0, fmt.Errorf("en la tabla %s, al hashear %T, no tenia la misma estructura que la esperada", dt.NombreTabla, datos)
	}

	datosRepresentativos := []any{}
	for i, claveTipo := range dt.ClavesTipo {
		if claveTipo.Representativa {
			datosRepresentativos = append(datosRepresentativos, datos[i])
		}
	}

	for _, referencia := range dt.Referencias {
		if !referencia.Representativo {
			continue
		}

		encontrado := false
		for _, fKey := range fKeys {
			for _, tabla := range referencia.Tablas {
				if referencia.Clave == fKey.Clave && tabla.NombreTabla == fKey.TablaDestino {
					encontrado = true
					datosRepresentativos = append(datosRepresentativos, fKey.HashDatosDestino)
					break
				}
			}
		}

		if !encontrado {
			return 0, fmt.Errorf("no tiene la foreign key necesaria para hacer su hash")
		}
	}

	return hash.HasearDatos(datosRepresentativos...), nil
}

func (dt DescripcionTabla) Existe(bdd *b.Bdd, datos ...any) (bool, error) {
	if !dt.necesarioQuery {
		return false, nil
	}

	_, err := bdd.Obtener(dt.query, datos...)
	return err == nil, nil
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
