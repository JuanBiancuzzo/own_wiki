package dependencias

import (
	"fmt"
	"math/rand"
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

type FnExiste func(tx b.Transaccion, datosIngresados ConjuntoDato, lock *sync.Mutex) (bool, error)
type FnInsertar func(tx b.Transaccion, datosIngresados ConjuntoDato, lock *sync.Mutex) (int64, error)
type FnUpdateClave func(clave string) (b.Sentencia, error)

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
	Existe                 FnExiste
	Insertar               FnInsertar
	CrearForeignKey        FnFKeys
	Hash                   FnHash
	ObtenerSentenciaUpdate FnUpdateClave
}

func ConstruirTabla(tracker *TrackerDependencias, nombreTabla string, tipoTabla TipoTabla, elementosRepetidos bool, variables []Variable) (Tabla, error) {
	var err error
	if err = crearTabla(tracker.Bdd, nombreTabla, variables); err != nil {
		return Tabla{}, err
	}

	var existe FnExiste
	if elementosRepetidos {
		existe = func(tx b.Transaccion, datosIngresados ConjuntoDato, lock *sync.Mutex) (bool, error) {
			return false, nil
		}

	} else if existe, err = generarExiste(tracker.Bdd, nombreTabla, variables); err != nil {
		return Tabla{}, err
	}

	var insertar FnInsertar
	if insertar, err = generarInsertar(nombreTabla, tracker, variables); err != nil {
		return Tabla{}, err
	}

	var sentenciasUpdate FnUpdateClave
	if sentenciasUpdate, err = generarUpdate(tracker.Bdd, nombreTabla, variables); err != nil {
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

		Existe:                 existe,
		Insertar:               insertar,
		CrearForeignKey:        generarFKeys(variables),
		Hash:                   generarHash(variables),
		ObtenerDependencias:    describirDependencias(variables),
		ObtenerSentenciaUpdate: sentenciasUpdate,
	}, nil
}

func crearTabla(bdd *b.Bdd, nombreTabla string, variables []Variable) error {
	parametros := []string{}

	for _, variable := range variables {
		parametros = append(parametros, variable.ObtenerParametroSQL()...)
	}

	tabla := fmt.Sprintf(
		"CREATE TABLE %s (\n\tid BIGINT PRIMARY KEY,\n\t%s\n);",
		nombreTabla,
		strings.Join(parametros, ",\n\t"),
	)

	if err := bdd.EliminarTabla(nombreTabla); err != nil {
		return fmt.Errorf("no se pudo crear la tabla \n%s\n, con error: %v", tabla, err)
	}

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
		return func(tx b.Transaccion, datosIngresados ConjuntoDato, lock *sync.Mutex) (bool, error) {
			return false, nil
		}, nil
	}

	sentenciaQuery, err := bdd.Preparar(query)
	if err != nil {
		return nil, fmt.Errorf("al preparar la sentencia '%s' se tuvo el error: %v", query, err)
	}

	return func(tx b.Transaccion, datosIngresados ConjuntoDato, lock *sync.Mutex) (bool, error) {
		sentencia := tx.Sentencia(sentenciaQuery)
		defer sentencia.Close()

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

func generarUpdate(bdd *b.Bdd, nombreTabla string, variables []Variable) (FnUpdateClave, error) {
	mapaSentencias := make(map[string]b.Sentencia)

	for _, variable := range variables {
		if _, ok := variable.Informacion.(VariableReferencia); !ok {
			continue
		}

		clave := variable.Clave

		query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", nombreTabla, clave)
		if sentencia, err := bdd.Preparar(query); err != nil {
			return nil, err

		} else {
			mapaSentencias[clave] = sentencia
		}
	}

	return func(clave string) (b.Sentencia, error) {
		if sentencia, ok := mapaSentencias[clave]; ok {
			return sentencia, nil
		}

		return b.Sentencia{}, fmt.Errorf("no se reconoce la clave")
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
		"INSERT INTO %s (id, %s) VALUES (?, %s)",
		nombreTabla,
		strings.Join(clavesTotales, ", "),
		strings.Join(valores, ", "),
	)
	sentenciaInsertar, err := tracker.Bdd.Preparar(insertarQuery)
	if err != nil {
		return nil, fmt.Errorf("al preparar la sentencia '%s' se tuvo el error: %v", insertarQuery, err)
	}

	cargarExtra := func(datosIngresados ConjuntoDato) {
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
	}

	r := rand.New(rand.NewSource(0))
	largoDatos := len(clavesInsertar)
	return func(tx b.Transaccion, datosIngresados ConjuntoDato, lock *sync.Mutex) (int64, error) {
		sentencia := tx.Sentencia(sentenciaInsertar)
		defer sentencia.Close()

		datos := make([]any, largoDatos+1)
		id := int64(r.Uint32())
		datos[0] = id

		if nombreTabla == "Carreras" {
			fmt.Printf("Datos: %+v\n", datosIngresados)
		}

		for i, clave := range clavesInsertar {
			if dato, ok := datosIngresados[clave]; !ok {
				return id, fmt.Errorf("el usuario no ingreso el dato para %s", clave)

			} else if relacion, ok := dato.(RelacionTabla); ok {
				// podemos hacer esto porque claves solo elige para los que tienen multiples
				datos[i+1] = relacion.Tabla

			} else {
				datos[i+1] = dato
			}
		}

		lock.Lock()
		_, err := sentencia.InsertarId(datos...)
		lock.Unlock()
		if err != nil {
			return id, err
		}

		if len(clavesExternas) > 0 {
			go cargarExtra(datosIngresados)
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
func (dt Tabla) RestringirTabla(tx b.Transaccion) error {
	return nil
}

func EsTipoDependiente(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == DEPENDIENTE_NO_DEPENDIBLE
}

func EsTipoDependible(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == INDEPENDIENTE_DEPENDIBLE
}
