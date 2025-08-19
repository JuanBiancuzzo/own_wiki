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

type FnDependencias func() []DescripcionTabla
type FnExiste func(bdd *b.Bdd, datosIngresados ConjuntoDato) (bool, error)
type FnInsertar func(bdd *b.Bdd, datosIngresados ConjuntoDato) (int64, error)
type FnFKeys func(hash *Hash, datosIngresados ConjuntoDato) ([]ForeignKey, error)
type FnHash func(hash *Hash, datosIngresados ConjuntoDato) (IntFK, error)
type FnTabla func(bdd *b.Bdd) error

type DescripcionTabla struct {
	NombreTabla   string
	TipoTabla     TipoTabla
	TipoDadoClave map[string]TipoVariable

	// funciones pre computadas
	Existe              FnExiste
	Insertar            FnInsertar
	CrearForeignKey     FnFKeys
	Hash                FnHash
	CrearTablaRelajada  FnTabla
	ObtenerDependencias FnDependencias
}

func ConstruirTabla(nombreTabla string, tipoTabla TipoTabla, elementosRepetidos bool, clavesTipo []ParClaveTipo, referencias []ReferenciaTabla) DescripcionTabla {
	tipoDadoClave := make(map[string]TipoVariable)
	tipoDadoClave["id"] = TV_INT

	for _, claveTipo := range clavesTipo {
		tipoDadoClave[claveTipo.Clave] = claveTipo.tipo
	}

	for _, referencia := range referencias {
		tipoDadoClave[referencia.Clave] = TV_REFERENCIA
		if len(referencia.Tablas) > 1 {
			tipoDadoClave[fmt.Sprintf("tipo%s", referencia.Clave)] = TV_ENUM
		}
	}

	return DescripcionTabla{
		NombreTabla:   nombreTabla,
		TipoTabla:     tipoTabla,
		TipoDadoClave: tipoDadoClave,

		Existe:              generarExiste(nombreTabla, elementosRepetidos, clavesTipo, referencias),
		Insertar:            generarInsertar(nombreTabla, clavesTipo, referencias),
		CrearForeignKey:     generarFKeys(referencias),
		Hash:                generarHash(clavesTipo, referencias),
		CrearTablaRelajada:  generarCrearTabla(nombreTabla, clavesTipo, referencias),
		ObtenerDependencias: generarObtenerDependencias(referencias),
	}
}

func generarExiste(nombreTabla string, elementosRepetidos bool, clavesTipo []ParClaveTipo, referencias []ReferenciaTabla) FnExiste {
	if elementosRepetidos {
		return func(bdd *b.Bdd, datosIngresados ConjuntoDato) (bool, error) {
			return false, nil
		}
	}

	queryParam := []string{}
	claves := []string{} // tiene en cuenta incluso las claves que tienen valores multiples unicamente

	for _, claveTipo := range clavesTipo {
		if claveTipo.Representativa {
			claves = append(claves, claveTipo.Clave)
			queryParam = append(queryParam, fmt.Sprintf("%s = ?", claveTipo.Clave))
		}
	}

	for _, referencia := range referencias {
		if referencia.Representativo && len(referencia.Tablas) > 1 {
			claves = append(claves, referencia.Clave)
			queryParam = append(queryParam, fmt.Sprintf("tipo%s = ?", referencia.Clave))
		}
	}

	query := fmt.Sprintf(
		"SELECT id FROM %s WHERE %s",
		nombreTabla,
		strings.Join(queryParam, " AND "),
	)

	largoDatos := len(claves)

	return func(bdd *b.Bdd, datosIngresados ConjuntoDato) (bool, error) {
		datos := make([]any, largoDatos)
		for _, clave := range claves {
			if dato, ok := datosIngresados[clave]; !ok {
				return false, fmt.Errorf("el usuario no ingreso el dato para %s", clave)

			} else if relacion, ok := dato.(RelacionTabla); ok {
				// podemos hacer esto porque claves solo elige para los que tienen multiples
				datos = append(datos, relacion.Tabla)
			} else {
				datos = append(datos, dato)
			}
		}

		_, err := bdd.Obtener(query, datos...)
		return err == nil, nil
	}
}

func generarInsertar(nombreTabla string, clavesTipo []ParClaveTipo, referencias []ReferenciaTabla) FnInsertar {
	// Las claves insertar no tienen que tener a las claves ref que no tengan multiples
	clavesInsertar := []string{}
	clavesTotales := []string{}
	valores := []string{}

	for _, claveTipo := range clavesTipo {
		clavesTotales = append(clavesTotales, claveTipo.Clave)
		clavesInsertar = append(clavesInsertar, claveTipo.Clave)
		valores = append(valores, "?")
	}

	for _, referencia := range referencias {

		if len(referencia.Tablas) > 1 {
			clavesTotales = append(clavesTotales, fmt.Sprintf("tipo%s", referencia.Clave))
			clavesInsertar = append(clavesInsertar, referencia.Clave)
			valores = append(valores, "?")
		}
		clavesTotales = append(clavesTotales, referencia.Clave)
		valores = append(valores, "0")
	}

	// Este ya tiene que tener los 0 en las referencias, asi no las tenemos q agregar
	insertarQuery := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		nombreTabla,
		strings.Join(clavesTotales, ", "),
		strings.Join(valores, ", "),
	)

	largoDatos := len(clavesInsertar)
	return func(bdd *b.Bdd, datosIngresados ConjuntoDato) (int64, error) {
		datos := make([]any, largoDatos)

		for i, clave := range clavesInsertar {
			if dato, ok := datosIngresados[clave]; !ok {
				return 0, fmt.Errorf("el usuario no ingreso el dato para %s", clave)

			} else if relacion, ok := dato.(RelacionTabla); ok {
				// podemos hacer esto porque claves solo elige para los que tienen multiples
				datos[i] = relacion.Tabla

			} else {
				datos[i] = dato
			}
		}

		return bdd.Insertar(insertarQuery, datos...)
	}
}

func generarFKeys(referencias []ReferenciaTabla) FnFKeys {
	fnHashs := []func(tabla string) FnHash{}
	claves := []string{}

	for _, referencia := range referencias {
		claves = append(claves, referencia.Clave)
		if len(referencia.Tablas) > 1 {
			tablasHash := make(map[string]FnHash)
			for _, tabla := range referencia.Tablas {
				tablasHash[tabla.NombreTabla] = tabla.Hash
			}

			fnHashs = append(fnHashs, func(nombreTabla string) FnHash {
				return tablasHash[nombreTabla]
			})

		} else {
			tabla := referencia.Tablas[0]
			fnHashs = append(fnHashs, func(_ string) FnHash { return tabla.Hash })
		}
	}

	largoClaves := len(claves)
	return func(hash *Hash, datosIngresados ConjuntoDato) ([]ForeignKey, error) {
		fKeys := make([]ForeignKey, largoClaves)

		for i, clave := range claves {
			relacion, ok := datosIngresados[clave].(RelacionTabla)
			if !ok {
				return fKeys, fmt.Errorf("el elemento relacionado a la clave %s no fue pasado", clave)
			}

			fnHash := fnHashs[i](relacion.Tabla)
			if hash, err := fnHash(hash, relacion.Datos); err != nil {
				return fKeys, fmt.Errorf("no se pudo general el hash de los datos, con err: %v", err)
			} else {
				fKeys[i] = NewForeignKey(relacion.Tabla, clave, hash)
			}
		}

		return fKeys, nil
	}
}

func generarHash(clavesTipo []ParClaveTipo, referencias []ReferenciaTabla) FnHash {
	clavesRepresentativas := []string{}
	tieneReferenciasRepresentativos := false
	fnHashs := []func(tabla string) FnHash{}

	for _, claveTipo := range clavesTipo {
		if claveTipo.Representativa {
			clavesRepresentativas = append(clavesRepresentativas, claveTipo.Clave)
			fnHashs = append(fnHashs, nil)
		}
	}

	for _, referencia := range referencias {
		if referencia.Representativo {
			clavesRepresentativas = append(clavesRepresentativas, referencia.Clave)
			tieneReferenciasRepresentativos = true

			if len(referencia.Tablas) > 1 {
				tablasHash := make(map[string]FnHash)
				for _, tabla := range referencia.Tablas {
					tablasHash[tabla.NombreTabla] = tabla.Hash
				}

				fnHashs = append(fnHashs, func(nombreTabla string) FnHash {
					return tablasHash[nombreTabla]
				})

			} else {
				tabla := referencia.Tablas[0]
				fnHashs = append(fnHashs, func(_ string) FnHash { return tabla.Hash })
			}
		}
	}

	largoDatos := len(clavesRepresentativas)
	if !tieneReferenciasRepresentativos {
		return func(hash *Hash, datosIngresados ConjuntoDato) (IntFK, error) {
			datosRepresentativos := make([]any, largoDatos)

			for i, clave := range clavesRepresentativas {
				if dato, ok := datosIngresados[clave]; !ok {
					return 0, fmt.Errorf("no se ingreso el valor para la clave %s", clave)
				} else {
					datosRepresentativos[i] = dato
				}
			}

			return hash.HasearDatos(datosRepresentativos...), nil
		}
	}

	return func(hash *Hash, datosIngresados ConjuntoDato) (IntFK, error) {
		datosRepresentativos := make([]any, largoDatos)

		for i, clave := range clavesRepresentativas {
			if dato, ok := datosIngresados[clave]; !ok {
				return 0, fmt.Errorf("no se ingreso el valor para la clave %s", clave)

			} else if relacion, ok := dato.(RelacionTabla); ok {
				fnHash := fnHashs[i](relacion.Tabla)
				if hash, err := fnHash(hash, relacion.Datos); err != nil {
					return 0, err

				} else {
					datosRepresentativos[i] = hash
				}

			} else {
				datosRepresentativos[i] = dato
			}
		}

		return hash.HasearDatos(datosRepresentativos...), nil
	}
}

func generarObtenerDependencias(referencias []ReferenciaTabla) FnDependencias {
	return func() []DescripcionTabla {
		tablas := []DescripcionTabla{}
		for _, referencia := range referencias {
			for _, tabla := range referencia.Tablas {
				tablas = append(tablas, *tabla)
			}
		}
		return tablas
	}
}

func generarCrearTabla(nombreTabla string, clavesTipo []ParClaveTipo, referencias []ReferenciaTabla) FnTabla {
	parametros := []string{}
	for _, parClaveTipo := range clavesTipo {
		parametros = append(parametros, fmt.Sprintf("%s %s", parClaveTipo.Clave, parClaveTipo.TipoSQL))
	}

	for _, referencia := range referencias {
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
		nombreTabla,
		strings.Join(parametros, ",\n\t"),
	)

	return func(bdd *b.Bdd) error {
		if err := bdd.CrearTabla(tabla); err != nil {
			return fmt.Errorf("no se pudo crear la tabla \n%s\n, con error: %v", tabla, err)
		}
		return nil
	}
}

// TODO
func (dt DescripcionTabla) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func EsTipoDependiente(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == DEPENDIENTE_NO_DEPENDIBLE
}

func EsTipoDependible(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == INDEPENDIENTE_DEPENDIBLE
}
