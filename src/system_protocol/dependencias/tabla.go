package dependencias

import (
	"fmt"
	b "own_wiki/system_protocol/base_de_datos"
	"strings"
	"sync"
)

type TipoTabla byte

const (
	DEPENDIENTE_NO_DEPENDIBLE   = 0b00
	DEPENDIENTE_DEPENDIBLE      = 0b10
	INDEPENDIENTE_NO_DEPENDIBLE = 0b01
	INDEPENDIENTE_DEPENDIBLE    = 0b11
)

type FnExiste func(datosIngresados ConjuntoDato, lock *sync.Mutex) (bool, error)
type FnInsertar func(datosIngresados ConjuntoDato, lock *sync.Mutex) (int64, error)

// type FnActualizar func(bdd *b.Bdd, datosIngresados ConjuntoDato) error
// type FnEliminar func(bdd *b.Bdd, id int64) error

type FnFKeys func(hash *Hash, datosIngresados ConjuntoDato) ([]ForeignKey, error)
type FnHash func(hash *Hash, datosIngresados ConjuntoDato) (IntFK, error)

type Tabla struct {
	NombreTabla         string
	TipoTabla           TipoTabla
	Variables           map[string]Variable
	ObtenerDependencias []Tabla

	// funciones pre computadas
	Existe          FnExiste
	Insertar        FnInsertar
	CrearForeignKey FnFKeys
	Hash            FnHash
}

func ConstruirTabla(tracker *TrackerDependencias, nombreTabla string, tipoTabla TipoTabla, elementosRepetidos bool, variables []Variable) (Tabla, error) {
	var err error
	if err = crearTabla(tracker.Bdd, nombreTabla, variables); err != nil {
		return Tabla{}, err
	}

	var existe FnExiste
	if elementosRepetidos {
		existe = func(datosIngresados ConjuntoDato, lock *sync.Mutex) (bool, error) { return false, nil }

	} else if existe, err = generarExiste(tracker.Bdd, nombreTabla, variables); err != nil {
		return Tabla{}, err
	}

	var insertar FnInsertar
	if insertar, err = generarInsertar(nombreTabla, tracker, variables); err != nil {
		return Tabla{}, err
	}

	variablesPorNombre := make(map[string]Variable)
	for _, variable := range variables {
		variablesPorNombre[variable.Clave] = variable
	}

	return Tabla{
		NombreTabla: nombreTabla,
		Variables:   variablesPorNombre,
		TipoTabla:   tipoTabla,

		Existe:              existe,
		Insertar:            insertar,
		CrearForeignKey:     generarFKeys(variables),
		Hash:                generarHash(variables),
		ObtenerDependencias: describirDependencias(variables),
	}, nil
}

func crearTabla(bdd *b.Bdd, nombreTabla string, variables []Variable) error {
	parametros := []string{}

	for _, variable := range variables {
		parametros = append(parametros, variable.ObtenerParametroSQL()...)
	}

	tabla := fmt.Sprintf(
		"CREATE TABLE %s (\n\tid INT AUTO_INCREMENT PRIMARY KEY,\n\t%s\n);",
		nombreTabla,
		strings.Join(parametros, ",\n\t"),
	)

	if err := bdd.EliminarTabla(nombreTabla); err != nil {
		return fmt.Errorf("no se pudo crear la tabla \n%s\n, con error: %v", tabla, err)
	}

	fmt.Printf("Tabla para %s:\n%s\n", nombreTabla, tabla)
	if err := bdd.CrearTabla(tabla); err != nil {
		return fmt.Errorf("no se pudo crear la tabla \n%s\n, con error: %v", tabla, err)
	}

	return nil
}

func generarExiste(bdd *b.Bdd, nombreTabla string, variables []Variable) (FnExiste, error) {
	queryParam := []string{}
	claves := []string{} // tiene en cuenta incluso las claves que tienen valores multiples unicamente

	for _, variable := range variables {
		switch variable.Informacion.(type) {
		case VariableSimple:
			queryParam = append(queryParam, fmt.Sprintf("%s = ?", variable.Clave))
			claves = append(claves, variable.Clave)

		case VariableString:
			queryParam = append(queryParam, fmt.Sprintf("%s = ?", variable.Clave))
			claves = append(claves, variable.Clave)
		case VariableEnum:
			queryParam = append(queryParam, fmt.Sprintf("%s = ?", variable.Clave))
			claves = append(claves, variable.Clave)

		case VariableReferencia:
			// queryParam = append(queryParam, fmt.Sprintf("tipo%s = ?", variable.Clave))
			// claves = append(claves, variable.Clave)

		case VariableArrayReferencia:
			// si la variable es esta, no deberia hacer nada porque no es un valor posible para buscar si existe
		}
	}

	query := fmt.Sprintf(
		"SELECT id FROM %s WHERE %s",
		nombreTabla,
		strings.Join(queryParam, " AND "),
	)

	largoDatos := len(claves)
	if largoDatos == 0 {
		return func(datosIngresados ConjuntoDato, lock *sync.Mutex) (bool, error) {
			return false, nil
		}, nil
	}

	sentencia, err := bdd.Preparar(query)
	if err != nil {
		return nil, fmt.Errorf("al preparar la sentencia '%s' se tuvo el error: %v", query, err)
	}

	return func(datosIngresados ConjuntoDato, lock *sync.Mutex) (bool, error) {
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

		lock.Lock()
		_, err := sentencia.Query(datos...)
		lock.Unlock()
		return err == nil, nil
	}, nil
}

func generarInsertar(nombreTabla string, tracker *TrackerDependencias, variables []Variable) (FnInsertar, error) {
	// Las claves insertar no tienen que tener a las claves ref que no tengan multiples
	clavesInsertar := []string{}
	clavesTotales := []string{}
	valores := []string{}

	// Esto es para manejar los arrays
	clavesExternas := []string{}
	tablasExterna := []string{}
	selfClaves := []string{}

	for _, variable := range variables {
		switch informacion := variable.Informacion.(type) {
		case VariableSimple:
			clavesTotales = append(clavesTotales, variable.Clave)
			clavesInsertar = append(clavesInsertar, variable.Clave)
			valores = append(valores, "?")
		case VariableString:
			clavesTotales = append(clavesTotales, variable.Clave)
			clavesInsertar = append(clavesInsertar, variable.Clave)
			valores = append(valores, "?")
		case VariableEnum:
			clavesTotales = append(clavesTotales, variable.Clave)
			clavesInsertar = append(clavesInsertar, variable.Clave)
			valores = append(valores, "?")

		case VariableReferencia:
			if len(informacion.Tablas) > 1 {
				clavesTotales = append(clavesTotales, fmt.Sprintf("tipo%s", variable.Clave))
				clavesInsertar = append(clavesInsertar, variable.Clave)
				valores = append(valores, "?")
			}
			clavesTotales = append(clavesTotales, variable.Clave)
			valores = append(valores, "0")

		case VariableArrayReferencia:
			clavesExternas = append(clavesExternas, variable.Clave)
			tablasExterna = append(tablasExterna, informacion.TablaCreada)
			selfClaves = append(selfClaves, informacion.ClaveSelf)
		}
	}

	// Este ya tiene que tener los 0 en las referencias, asi no las tenemos q agregar
	insertarQuery := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		nombreTabla,
		strings.Join(clavesTotales, ", "),
		strings.Join(valores, ", "),
	)
	sentencia, err := tracker.Bdd.Preparar(insertarQuery)
	if err != nil {
		return nil, fmt.Errorf("al preparar la sentencia '%s' se tuvo el error: %v", insertarQuery, err)
	}

	largoDatos := len(clavesInsertar)
	return func(datosIngresados ConjuntoDato, lock *sync.Mutex) (int64, error) {
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

		lock.Lock()
		id, err := sentencia.InsertarId(datos...)
		lock.Unlock()
		if err != nil {
			return 0, err
		}

		for i, clave := range clavesExternas {
			if dato, ok := datosIngresados[clave]; !ok {
				continue

			} else if datosRelacion, ok := dato.([]ConjuntoDato); ok {
				for _, datoRelacion := range datosRelacion {
					datoRelacion[selfClaves[i]] = NewRelacion(nombreTabla, datosIngresados)
					tracker.Cargar(tablasExterna[i], datoRelacion)
				}
			}
		}

		return id, nil
	}, nil
}

func generarFKeys(variables []Variable) FnFKeys {
	fnHashs := []func(tabla string) FnHash{}
	claves := []string{}

	for _, variable := range variables {
		if informacion, ok := variable.Informacion.(VariableReferencia); ok {
			claves = append(claves, variable.Clave)
			if len(informacion.Tablas) > 1 {
				tablasHash := make(map[string]FnHash)
				for _, tabla := range informacion.Tablas {
					tablasHash[tabla.NombreTabla] = tabla.Hash
				}

				fnHashs = append(fnHashs, func(nombreTabla string) FnHash {
					return tablasHash[nombreTabla]
				})

			} else {
				tabla := informacion.Tablas[0]
				fnHashs = append(fnHashs, func(_ string) FnHash { return tabla.Hash })
			}
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

func generarHash(variables []Variable) FnHash {
	clavesRepresentativas := []string{}
	tieneReferenciasRepresentativos := false
	fnHashs := []func(tabla string) FnHash{}

	for _, variable := range variables {
		switch informacion := variable.Informacion.(type) {
		case VariableSimple:
			if informacion.Representativo {
				clavesRepresentativas = append(clavesRepresentativas, variable.Clave)
				fnHashs = append(fnHashs, nil)
			}
		case VariableString:
			if informacion.Representativo {
				clavesRepresentativas = append(clavesRepresentativas, variable.Clave)
				fnHashs = append(fnHashs, nil)
			}
		case VariableEnum:
			if informacion.Representativo {
				clavesRepresentativas = append(clavesRepresentativas, variable.Clave)
				fnHashs = append(fnHashs, nil)
			}

		case VariableReferencia:
			if informacion.Representativo {
				clavesRepresentativas = append(clavesRepresentativas, variable.Clave)
				tieneReferenciasRepresentativos = true

				if len(informacion.Tablas) > 1 {
					tablasHash := make(map[string]FnHash)
					for _, tabla := range informacion.Tablas {
						tablasHash[tabla.NombreTabla] = tabla.Hash
					}

					fnHashs = append(fnHashs, func(nombreTabla string) FnHash {
						return tablasHash[nombreTabla]
					})

				} else {
					tabla := informacion.Tablas[0]
					fnHashs = append(fnHashs, func(_ string) FnHash { return tabla.Hash })
				}
			}
		case VariableArrayReferencia:
			// Para este en particular, no la necestio, por lo tanto es irrelevante
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

func describirDependencias(variables []Variable) []Tabla {
	tablas := []Tabla{}
	for _, variable := range variables {
		if informacion, ok := variable.Informacion.(VariableReferencia); ok {
			for _, tabla := range informacion.Tablas {
				tablas = append(tablas, *tabla)
			}
		}
	}

	return tablas
}

// TODO
func (dt Tabla) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func EsTipoDependiente(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == DEPENDIENTE_NO_DEPENDIBLE
}

func EsTipoDependible(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == INDEPENDIENTE_DEPENDIBLE
}
