package dependencias

import (
	_ "embed"
	"fmt"
	"sync"

	b "own_wiki/system_protocol/base_de_datos"
)

const AUX_DEPENDIBLES = "aux_dependibles"
const TABLA_DEPENDIBLES = `CREATE TABLE IF NOT EXISTS aux_dependibles (
	nombreTabla TEXT CHECK( LENGTH(nombreTabla) <= %d ) NOT NULL,
	hashDatos   BIGINT,
	idDatos     BIGINT
);`

const AUX_INCOMPLETOS = "aux_incompletos"
const TABLA_INCOMPLETOS = `CREATE TABLE IF NOT EXISTS aux_incompletos (
	tablaDependiente 	TEXT CHECK( LENGTH(tablaDependiente) <= %d ) NOT NULL,
	idDependiente   	BIGINT,
	keyAlId 			VARCHAR(255),
	tablaDestino 		TEXT CHECK( LENGTH(tablaDestino) <= %d ) NOT NULL,
	hashDatosDestino   	BIGINT
);`

type SentenciasTracker byte

const (
	ST_INSERTAR_TABLA_DEPENDIENTES = iota
	ST_QUERY_TABLA_DEPENDIENTES
	ST_INSERTAR_TABLA_INCOMPLETOS
	ST_QUERY_TABLA_INCOMPLETOS
	ST_ELIMINAR_TABLA_INCOMPLETOS
	ST_QUERY_TODO_INCOMPLETOS
	ST_ELIMINAR_TODO_INCOMPLETOS

	ST_MAX_SENTENCIA
)

func (st SentenciasTracker) Sentencia() string {
	switch st {
	case ST_INSERTAR_TABLA_DEPENDIENTES:
		return "INSERT INTO aux_dependibles (nombreTabla, idDatos, hashDatos) VALUES (?, ?, ?)"

	case ST_QUERY_TABLA_DEPENDIENTES:
		return "SELECT idDatos FROM aux_dependibles WHERE nombreTabla = ? AND hashDatos = ?"

	case ST_INSERTAR_TABLA_INCOMPLETOS:
		return "INSERT INTO aux_incompletos (tablaDependiente, idDependiente, keyAlId, tablaDestino, hashDatosDestino) VALUES (?, ?, ?, ?, ?)"

	case ST_QUERY_TABLA_INCOMPLETOS:
		return "SELECT tablaDependiente, idDependiente, keyAlId FROM aux_incompletos WHERE tablaDestino = ? AND hashDatosDestino = ?"

	case ST_ELIMINAR_TABLA_INCOMPLETOS:
		return "DELETE FROM aux_incompletos WHERE tablaDestino = ? AND hashDatosDestino = ?"

	case ST_QUERY_TODO_INCOMPLETOS:
		return "SELECT tablaDependiente, idDependiente, keyAlId, tablaDestino, aux_dependibles.idDatos FROM aux_incompletos INNER JOIN aux_dependibles ON aux_dependibles.hashDatos = aux_incompletos.hashDatosDestino AND aux_dependibles.nombreTabla = aux_incompletos.tablaDestino;"

	case ST_ELIMINAR_TODO_INCOMPLETOS:
		return "DELETE FROM aux_incompletos"
	}

	return fmt.Sprintf("ERROR: %v", st)
}

func (st SentenciasTracker) String() string {
	switch st {
	case ST_INSERTAR_TABLA_DEPENDIENTES:
		return "ST_INSERTAR_TABLA_DEPENDIENTES"

	case ST_QUERY_TABLA_DEPENDIENTES:
		return "ST_QUERY_TABLA_DEPENDIENTES"

	case ST_INSERTAR_TABLA_INCOMPLETOS:
		return "ST_INSERTAR_TABLA_INCOMPLETOS"

	case ST_QUERY_TABLA_INCOMPLETOS:
		return "ST_QUERY_TABLA_INCOMPLETOS"

	case ST_ELIMINAR_TABLA_INCOMPLETOS:
		return "ST_ELIMINAR_TABLA_INCOMPLETOS"

	case ST_QUERY_TODO_INCOMPLETOS:
		return "ST_QUERY_TODO_INCOMPLETOS"

	case ST_ELIMINAR_TODO_INCOMPLETOS:
		return "ST_ELIMINAR_TODO_INCOMPLETOS"
	}

	return fmt.Sprintf("no hay una sentencia con el numero: %d", st)
}

type ConjuntoDato map[string]any

type DatosCarga struct {
	Tabla *Tabla
	Datos ConjuntoDato
}

type TrackerDependencias struct {
	Tablas          *b.Bdd
	Temp            *b.Bdd
	RegistrarTablas map[string]Tabla
	Hash            *Hash

	canalCarga chan DatosCarga
	waitCarga  *sync.WaitGroup

	maximoNombreTabla int
	locksTablas       map[string]*sync.Mutex
	lockIncompletos   *sync.Mutex
	lockDependibles   *sync.Mutex

	sentencias []b.Sentencia
}

func NewTrackerDependencias(bdd *b.Bdd, canalMensajes chan string) (*TrackerDependencias, error) {
	var lockIncompletos sync.Mutex
	var lockDependibles sync.Mutex
	var waitCarga sync.WaitGroup

	tempBdd, err := b.NewBdd("/temp", "temp.db", canalMensajes)
	if err != nil {
		return nil, fmt.Errorf("error al crear bdd temp con: %v", err)
	}

	return &TrackerDependencias{
		Tablas:          bdd,
		Temp:            tempBdd,
		RegistrarTablas: make(map[string]Tabla),
		Hash:            NewHash(),

		waitCarga:  &waitCarga,
		canalCarga: make(chan DatosCarga, 100),

		maximoNombreTabla: 0,
		locksTablas:       make(map[string]*sync.Mutex),
		lockIncompletos:   &lockIncompletos,
		lockDependibles:   &lockDependibles,

		sentencias: []b.Sentencia{},
	}, nil
}

func (td *TrackerDependencias) CargarTabla(descripcion Tabla) {
	td.maximoNombreTabla = max(td.maximoNombreTabla, len(descripcion.NombreTabla))

	td.RegistrarTablas[descripcion.NombreTabla] = descripcion
	var lock sync.Mutex
	td.locksTablas[descripcion.NombreTabla] = &lock
}

func (td *TrackerDependencias) EmpezarProcesoInsertarDatos(canalMensajes chan string) error {
	if err := td.Temp.CrearTabla(fmt.Sprintf(TABLA_DEPENDIBLES, td.maximoNombreTabla)); err != nil {
		return fmt.Errorf(
			"creando tabla dependibles (\n%s\n), se tuvo el error: %v",
			fmt.Sprintf(TABLA_DEPENDIBLES, td.maximoNombreTabla), err,
		)

	} else if err := td.Temp.CrearTabla(fmt.Sprintf(TABLA_INCOMPLETOS, td.maximoNombreTabla, td.maximoNombreTabla)); err != nil {
		return fmt.Errorf(
			"creando tabla incompletos (\n%s\n), se tuvo el error: %v",
			fmt.Sprintf(TABLA_INCOMPLETOS, td.maximoNombreTabla, td.maximoNombreTabla), err,
		)
	}

	sentencias := make([]b.Sentencia, ST_MAX_SENTENCIA)
	var err error
	for i := range ST_MAX_SENTENCIA {
		query := SentenciasTracker(i)
		if sentencias[i], err = td.Temp.Preparar(query.Sentencia()); err != nil {
			return fmt.Errorf("preparando la sentencia %v se tuvo el error: %v", query, err)
		}
	}
	td.sentencias = sentencias

	td.waitCarga.Add(1)
	go td.procesarCarga(td.canalCarga, canalMensajes)

	return nil
}

func (td *TrackerDependencias) procesarCarga(canal chan DatosCarga, canalMensajes chan string) {
	capacidadDatos := 20
	cantidadDatos := 0
	leerDatos := make([]DatosCarga, capacidadDatos)

	contadorExtra := 0
	for datoCarga := range canal {
		leerDatos[cantidadDatos] = datoCarga
		cantidadDatos++

		if cantidadDatos >= capacidadDatos {
			contadorExtra++
			cantidadDatos = 0
			if err := td.cargarMultiplesDatos(td.Tablas, leerDatos); err != nil {
				canalMensajes <- fmt.Sprintf("%v", err)

			} else if contadorExtra%10 == 0 {
				td.Tablas.Checkpoint(b.TC_PASSIVE)
			}
		}
	}

	if cantidadDatos > 0 {
		if err := td.cargarMultiplesDatos(td.Tablas, leerDatos[:cantidadDatos]); err != nil {
			canalMensajes <- fmt.Sprintf("%v", err)
		}
	}

	td.waitCarga.Done()
}

func (td *TrackerDependencias) cargarMultiplesDatos(bdd *b.Bdd, datosCarga []DatosCarga) error {
	transaccion, err := bdd.Transaccion()
	if err != nil {
		return fmt.Errorf("error al crear transaccion con error: %v", err)
	}

	for _, datoCarga := range datosCarga {
		if err = td.cargarDato(transaccion, datoCarga.Tabla, datoCarga.Datos); err != nil {
			if errRB := transaccion.RollBack(); errRB != nil {
				return fmt.Errorf("[Rollback] error al cargar datos con error: %v, con rollback: %v", err, errRB)
			} else {
				return fmt.Errorf("[Rollback] error al cargar datos con error: %v", err)
			}
		}
	}

	if err = transaccion.Commit(); err != nil {
		return fmt.Errorf("[Commitear fail] error al commitear con error: %v", err)
	}
	return nil
}

func (td *TrackerDependencias) cargarDato(tx b.Transaccion, tabla *Tabla, datosIngresados ConjuntoDato) error {
	lock := td.locksTablas[tabla.NombreTabla]
	if existe, err := tabla.Existe(tx, datosIngresados, lock); err != nil {
		return err

	} else if existe {
		return nil
	}

	id, err := tabla.Insertar(tx, datosIngresados, lock)
	if err != nil {
		return fmt.Errorf("error insertando en cargar del tracker, con error: %v", err)
	}

	if EsTipoDependiente(tabla.TipoTabla) {
		fKeys, err := tabla.CrearForeignKey(td.Hash, datosIngresados)
		if err != nil {
			return fmt.Errorf("error al crear fkeys: %v", err)
		}

		if err := td.procesoDependiente(tx, tabla, id, fKeys); err != nil {
			return fmt.Errorf("error al verificar o actualizar el elemnto en la tabla %s, con id: %d, con error: %v", tabla.NombreTabla, id, err)
		}
	}

	if EsTipoDependible(tabla.TipoTabla) {
		if hashDatos, err := tabla.Hash(td.Hash, datosIngresados); err != nil {
			return fmt.Errorf("error al calcular el hash: %v", err)

		} else if err := td.procesoDependible(tx, tabla, id, hashDatos); err != nil {
			return fmt.Errorf("error al proceso dependible: %v", err)
		}
	}

	return nil
}

func (td *TrackerDependencias) Cargar(nombreTabla string, datosIngresados ConjuntoDato) error {
	if tabla, ok := td.RegistrarTablas[nombreTabla]; !ok {
		return fmt.Errorf("de alguna forma estas cargando en una tabla no registrada")

	} else {
		td.canalCarga <- DatosCarga{
			Tabla: &tabla,
			Datos: datosIngresados,
		}

		return nil
	}
}

func (td *TrackerDependencias) procesoDependiente(tx b.Transaccion, tabla *Tabla, idInsertado int64, fKeys []ForeignKey) error {
	lock := td.locksTablas[tabla.NombreTabla]

	for _, fKey := range fKeys {
		// Vemos si ya fue insertado la dependencia
		td.lockDependibles.Lock()
		if id, err := td.sentencias[ST_QUERY_TABLA_DEPENDIENTES].Obtener(fKey.TablaDestino, fKey.HashDatosDestino); err == nil {
			td.lockDependibles.Unlock()
			// Si fueron insertados, por lo que actualizamos la tabla
			sentenciaUpdate, err := tabla.ObtenerSentenciaUpdate(fKey.Clave)
			if err != nil {
				return fmt.Errorf("error no existe setencia update en la tabla %s, para la clave %s, dando err: %v", tabla.NombreTabla, fKey.Clave, err)
			}

			sentenciaUpdate = tx.Sentencia(sentenciaUpdate)
			defer sentenciaUpdate.Close()

			lock.Lock()
			if err = sentenciaUpdate.Update(id, idInsertado); err != nil {
				lock.Unlock()
				return fmt.Errorf("error al actualizar %d en tabla %s (proceso dependiente), con error %v", idInsertado, tabla.NombreTabla, err)
			}
			lock.Unlock()

		} else {
			td.lockDependibles.Unlock()

			// Como no fue insertada, tenemos que guardar la información para que se carge correctamente la dependencia
			datos := []any{tabla.NombreTabla, idInsertado, fKey.Clave, fKey.TablaDestino, fKey.HashDatosDestino}
			td.lockIncompletos.Lock()
			if _, err := td.sentencias[ST_INSERTAR_TABLA_INCOMPLETOS].InsertarId(datos...); err != nil {
				td.lockIncompletos.Unlock()
				return fmt.Errorf("error al insertar en la tabla auxiliar de incompletos, con error: %v", err)
			}
			td.lockIncompletos.Unlock()
		}
	}

	return nil
}

type UpdateData struct {
	tablaDependiente, key string
	idDependiente         int64
}

func (td *TrackerDependencias) procesoDependible(tx b.Transaccion, tabla *Tabla, idInsertado int64, hashDatos IntFK) error {
	td.lockDependibles.Lock()
	if _, err := td.sentencias[ST_INSERTAR_TABLA_DEPENDIENTES].InsertarId(tabla.NombreTabla, idInsertado, hashDatos); err != nil {
		td.lockDependibles.Unlock()
		return fmt.Errorf("error al insertar en dependientes: %s, con error: %v", tabla.NombreTabla, err)
	}
	td.lockDependibles.Unlock()

	td.lockIncompletos.Lock()
	if filas, err := td.sentencias[ST_QUERY_TABLA_INCOMPLETOS].Query(tabla.NombreTabla, hashDatos); err != nil {
		td.lockIncompletos.Unlock()
		return fmt.Errorf("error al query cuales son los elementos incompletos con tabla: %s, con error: %v", tabla.NombreTabla, err)

	} else {
		td.lockIncompletos.Unlock()

		updates := []UpdateData{}
		sentenciasPorTabla := make(map[string]b.Sentencia)
		for filas.Next() {
			var u UpdateData

			if err = filas.Scan(&u.tablaDependiente, &u.idDependiente, &u.key); err != nil {
				filas.Close()
				return fmt.Errorf("error al obtener datos de una query de incompletos, con error: %v", err)
			}

			sentenciaUpdate, err := td.RegistrarTablas[u.tablaDependiente].ObtenerSentenciaUpdate(u.key)
			if err != nil {
				filas.Close()
				return fmt.Errorf("error no existe setencia update en la tabla %s, para la clave %s, dando err: %v", u.tablaDependiente, u.key, err)
			}
			sentenciasPorTabla[u.tablaDependiente] = tx.Sentencia(sentenciaUpdate)

			updates = append(updates, u)
		}
		filas.Close()

		for _, u := range updates {
			sentenciaUpdate := sentenciasPorTabla[u.tablaDependiente]

			lock := td.locksTablas[u.tablaDependiente]
			lock.Lock()
			if err = sentenciaUpdate.Update(idInsertado, u.idDependiente); err != nil {
				lock.Unlock()
				return fmt.Errorf("error al actualizar %d en tabla %s (proceso Dependible distinto), con error %v", idInsertado, tabla.NombreTabla, err)
			}
			lock.Unlock()
		}

		for tabla := range sentenciasPorTabla {
			sentenciasPorTabla[tabla].Close()
		}

		if len(updates) > 0 {
			td.lockIncompletos.Lock()
			if err = td.sentencias[ST_ELIMINAR_TABLA_INCOMPLETOS].Eliminar(tabla.NombreTabla, hashDatos); err != nil {
				td.lockIncompletos.Unlock()
				return fmt.Errorf("error al eliminar %d en tabla %s, con error %v", idInsertado, tabla.NombreTabla, err)
			}
			td.lockIncompletos.Unlock()
		}

		return nil
	}
}

func (td *TrackerDependencias) TerminarProcesoInsertarDatos() error {
	close(td.canalCarga)
	td.waitCarga.Wait()

	/*
		transaccion, err := td.Bdd.Transaccion()
		if err != nil {
			return err
		}

		if err := td.procesoUltimasActualizaciones(transaccion); err != nil {
			errRollBack := transaccion.RollBack()
			return fmt.Errorf("error al procesar ultimos: %v, y rollback: %v", err, errRollBack)
		}

		for tabla := range td.RegistrarTablas {
			if err := td.RegistrarTablas[tabla].RestringirTabla(transaccion); err != nil {
				errRollBack := transaccion.RollBack()
				return fmt.Errorf("error al restringir tabla %s: %v, y rollback: %v", tabla, err, errRollBack)
			}
		}

			if err := td.Bdd.EliminarTabla(AUX_DEPENDIBLES); err != nil {
				return fmt.Errorf("error al eliminar tabla auxiliar dependibles, con error: %v", err)

			} else if err = td.Bdd.EliminarTabla(AUX_INCOMPLETOS); err != nil {
				return fmt.Errorf("error al eliminar tabla auxiliar incompletos, con error: %v", err)
			}

		return transaccion.Commit()
	*/

	td.Temp.Close()
	for _, sentencia := range td.sentencias {
		sentencia.Close()
	}
	return nil
}

func (td *TrackerDependencias) procesoUltimasActualizaciones(tx b.Transaccion) error {
	// Hacer un inner join pora obtener ya de por si los id, sin tener que buscarlos
	sentenciaQueryIncompletos := tx.Sentencia(td.sentencias[ST_QUERY_TODO_INCOMPLETOS])
	sentenciaEliminarIncompletos := tx.Sentencia(td.sentencias[ST_ELIMINAR_TODO_INCOMPLETOS])

	td.lockIncompletos.Lock()
	if filas, err := sentenciaQueryIncompletos.Query(); err != nil {
		td.lockIncompletos.Unlock()
		return fmt.Errorf("error al query de todos los elementos de incompletos, con error: %v", err)

	} else {
		td.lockIncompletos.Unlock()

		var tablaDependiente, key, tablaDestino string
		var idDependiente, idInsertado int64

		for filas.Next() {
			if err = filas.Scan(&tablaDependiente, &idDependiente, &key, &tablaDestino, &idInsertado); err != nil {
				filas.Close()
				return fmt.Errorf("error al obtener datos de una query de incompletos, con error: %v", err)
			}

			sentenciaUpdate, err := td.RegistrarTablas[tablaDependiente].ObtenerSentenciaUpdate(key)
			if err != nil {
				filas.Close()
				return fmt.Errorf("error no existe setencia update en la tabla %s, para la clave %s, dando err: %v", tablaDependiente, key, err)

			} else {
				sentenciaUpdate = tx.Sentencia(sentenciaUpdate)
			}

			td.locksTablas[tablaDependiente].Lock()
			if err = sentenciaUpdate.Update(idInsertado, idDependiente); err != nil {
				filas.Close()
				td.locksTablas[tablaDependiente].Unlock()
				return fmt.Errorf("error al actualizar %d en tabla %s, con error %v", idInsertado, tablaDestino, err)
			}
			td.locksTablas[tablaDependiente].Unlock()
		}
		filas.Close()

		if err = sentenciaEliminarIncompletos.Eliminar(); err != nil {
			return fmt.Errorf("error al eliminar el resto de la tabla de incompletos, con error %v", err)
		}

		return nil
	}
}
