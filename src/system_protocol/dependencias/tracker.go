package dependencias

import (
	_ "embed"
	"fmt"
	"slices"
	"sync"

	b "own_wiki/system_protocol/bass_de_datos"
	u "own_wiki/system_protocol/utilidades"
)

const AUX_DEPENDIBLES = "aux_dependibles"
const TABLA_DEPENDIBLES = `CREATE TABLE IF NOT EXISTS aux_dependibles (
	nombreTabla TEXT CHECK( LENGTH(nombreTabla) <= %d ) NOT NULL,
	hashDatos   BIGINT,
	idDatos     INT
);`

const AUX_INCOMPLETOS = "aux_incompletos"
const TABLA_INCOMPLETOS = `CREATE TABLE IF NOT EXISTS aux_incompletos (
	tablaDependiente 	TEXT CHECK( LENGTH(tablaDependiente) <= %d ) NOT NULL,
	idDependiente   	INT,
	keyAlId 			VARCHAR(255),
	tablaDestino 		TEXT CHECK( LENGTH(tablaDestino) <= %d ) NOT NULL,
	hashDatosDestino   	BIGINT
);`

const (
	INSERTAR_TABLA_DEPENDIENTES = "INSERT INTO aux_dependibles (nombreTabla, idDatos, hashDatos) VALUES (?, ?, ?)"
	QUERY_TABLA_DEPENDIENTES    = "SELECT idDatos FROM aux_dependibles WHERE nombreTabla = ? AND hashDatos = ?"
)

const (
	INSERTAR_TABLA_INCOMPLETOS = "INSERT INTO aux_incompletos (tablaDependiente, idDependiente, keyAlId, tablaDestino, hashDatosDestino) VALUES (?, ?, ?, ?, ?)"
	QUERY_TABLA_INCOMPLETOS    = "SELECT tablaDependiente, idDependiente, keyAlId FROM aux_incompletos WHERE tablaDestino = ? AND hashDatosDestino = ?"
	ELIMINAR_TABLA_INCOMPLETOS = "DELETE FROM aux_incompletos WHERE tablaDestino = ? AND hashDatosDestino = ?"
)

const (
	QUERY_TODO_INCOMPLETOS    = "SELECT tablaDependiente, idDependiente, keyAlId, tablaDestino, aux_dependibles.idDatos FROM aux_incompletos INNER JOIN aux_dependibles ON aux_dependibles.hashDatos = aux_incompletos.hashDatosDestino AND aux_dependibles.nombreTabla = aux_incompletos.tablaDestino;"
	ELIMINAR_TODO_INCOMPLETOS = "DELETE FROM aux_incompletos"
)

type ConjuntoDato map[string]any

type TrackerDependencias struct {
	BasesDeDatos    *b.Bdd
	RegistrarTablas map[string]Tabla
	Hash            *Hash

	tablasProcesar    *u.Cola[Tabla]
	maximoNombreTabla int

	locksTablas     map[string]*sync.Mutex
	lockIncompletos *sync.Mutex
	lockDependibles *sync.Mutex
}

func NewTrackerDependencias(bdd *b.Bdd) (*TrackerDependencias, error) {
	var lockIncompletos sync.Mutex
	var lockDependibles sync.Mutex

	return &TrackerDependencias{
		BasesDeDatos:    bdd,
		RegistrarTablas: make(map[string]Tabla),
		Hash:            NewHash(),

		tablasProcesar:    u.NewCola[Tabla](),
		maximoNombreTabla: 0,

		locksTablas:     make(map[string]*sync.Mutex),
		lockIncompletos: &lockIncompletos,
		lockDependibles: &lockDependibles,
	}, nil
}

func crearTablas(tablasProcesar *u.Cola[Tabla]) ([]Tabla, error) {
	// Creando las tablas relajadas
	var tablasOrdenadas []Tabla = []Tabla{}
	for tabla := range tablasProcesar.DesencolarIterativamente {
		nombreTabla := tabla.NombreTabla

		for i, tablaExistente := range tablasOrdenadas {
			if tablaExistente.NombreTabla == nombreTabla {
				tablasOrdenadas = append(tablasOrdenadas[:i], tablasOrdenadas[i+1:]...)
				break
			}
		}

		tablasOrdenadas = append([]Tabla{tabla}, tablasOrdenadas...)
		for _, tablaDependible := range tabla.ObtenerDependencias {
			tablasProcesar.Encolar(tablaDependible)
		}
	}

	return tablasOrdenadas, nil
}

func (td *TrackerDependencias) CargarTabla(descripcion Tabla) {
	td.maximoNombreTabla = max(td.maximoNombreTabla, len(descripcion.NombreTabla))

	td.RegistrarTablas[descripcion.NombreTabla] = descripcion
	var lock sync.Mutex
	td.locksTablas[descripcion.NombreTabla] = &lock

	if !EsTipoDependible(descripcion.TipoTabla) {
		td.tablasProcesar.Encolar(descripcion)
	}
}

func (td *TrackerDependencias) EmpezarProcesoInsertarDatos(canalMensajes chan string) error {
	if err := td.BasesDeDatos.CrearTabla(fmt.Sprintf(TABLA_DEPENDIBLES, td.maximoNombreTabla)); err != nil {
		return fmt.Errorf("creando tabla dependibles (\n%s\n), se tuvo el error: %v", fmt.Sprintf(TABLA_DEPENDIBLES, td.maximoNombreTabla), err)

	} else if err := td.BasesDeDatos.CrearTabla(fmt.Sprintf(TABLA_INCOMPLETOS, td.maximoNombreTabla, td.maximoNombreTabla)); err != nil {
		return fmt.Errorf("creando tabla incompletos (\n%s\n), se tuvo el error: %v", fmt.Sprintf(TABLA_INCOMPLETOS, td.maximoNombreTabla, td.maximoNombreTabla), err)

	} else if tablasOrdenadas, err := crearTablas(td.tablasProcesar); err != nil {
		return fmt.Errorf("creando tabla generales, se tuvo el error: %v", err)

	} else {
		for _, tabla := range slices.Backward(tablasOrdenadas) {
			if err := td.BasesDeDatos.EliminarTabla(tabla.NombreTabla); err != nil {
				canalMensajes <- fmt.Sprintf("error al eliminar tabla %s con error: %v", tabla.NombreTabla, err)
				continue
			}
		}

		canalMensajes <- "Orden final de cargado:"
		for _, tabla := range tablasOrdenadas {
			canalMensajes <- "Tabla: " + tabla.NombreTabla

			if err := tabla.CrearTablaRelajada(td.BasesDeDatos); err != nil {
				return fmt.Errorf("error al crear tablas relajadas, especificamente en %s, con error: %v", tabla.NombreTabla, err)
			}
		}
	}

	return nil
}

func (td *TrackerDependencias) TerminarProcesoInsertarDatos() error {
	if err := td.procesoUltimasActualizaciones(); err != nil {
		return err
	}

	for tabla := range td.RegistrarTablas {
		if err := td.RegistrarTablas[tabla].RestringirTabla(td.BasesDeDatos); err != nil {
			return err
		}
	}

	if err := td.BasesDeDatos.EliminarTabla(AUX_DEPENDIBLES); err != nil {
		return fmt.Errorf("error al eliminar tabla auxiliar dependibles, con error: %v", err)

	} else if err = td.BasesDeDatos.EliminarTabla(AUX_INCOMPLETOS); err != nil {
		return fmt.Errorf("error al eliminar tabla auxiliar incompletos, con error: %v", err)
	}

	return nil
}

func (td *TrackerDependencias) Cargar(nombreTabla string, datosIngresados ConjuntoDato) error {
	tabla, ok := td.RegistrarTablas[nombreTabla]
	if !ok {
		return fmt.Errorf("de alguna forma estas cargando en una tabla no registrada")
	}

	lock := td.locksTablas[nombreTabla]
	if existe, err := tabla.Existe(td.BasesDeDatos, datosIngresados, lock); err != nil {
		return err

	} else if existe {
		return nil
	}

	id, err := tabla.Insertar(td.BasesDeDatos, datosIngresados, lock)
	if err != nil {
		return fmt.Errorf("error insertando en cargar del tracker, con error: %v", err)
	}

	if EsTipoDependiente(tabla.TipoTabla) {
		fKeys, err := tabla.CrearForeignKey(td.Hash, datosIngresados)
		if err != nil {
			return err
		}

		if err := td.procesoDependiente(tabla, id, fKeys); err != nil {
			return fmt.Errorf("error al verificar o actualizar el elemnto en la tabla %s, con id: %d, con error: %v", tabla.NombreTabla, id, err)
		}
	}

	if EsTipoDependible(tabla.TipoTabla) {
		if hashDatos, err := tabla.Hash(td.Hash, datosIngresados); err != nil {
			return err
		} else {
			return td.procesoDependible(tabla, id, hashDatos)
		}
	}

	return nil
}

func (td *TrackerDependencias) procesoDependiente(tabla Tabla, idInsertado int64, fKeys []ForeignKey) error {
	lock := td.locksTablas[tabla.NombreTabla]

	for _, fKey := range fKeys {
		// Vemos si ya fue insertado la dependencia
		td.lockDependibles.Lock()
		if id, err := td.BasesDeDatos.Obtener(QUERY_TABLA_DEPENDIENTES, fKey.TablaDestino, fKey.HashDatosDestino); err == nil {
			td.lockDependibles.Unlock()
			// Si fueron insertados, por lo que actualizamos la tabla
			query := fmt.Sprintf("UPDATE %s SET %s = %d WHERE id = %d", tabla.NombreTabla, fKey.Clave, id, idInsertado)
			lock.Lock()
			if err = td.BasesDeDatos.Update(query); err != nil {
				lock.Unlock()
				return fmt.Errorf("error al actualizar %d en tabla %s (proceso dependiente), con error %v", idInsertado, tabla.NombreTabla, err)
			}
			lock.Unlock()

		} else {
			td.lockDependibles.Unlock()

			// Como no fue insertada, tenemos que guardar la información para que se carge correctamente la dependencia
			datos := []any{tabla.NombreTabla, idInsertado, fKey.Clave, fKey.TablaDestino, fKey.HashDatosDestino}
			td.lockIncompletos.Lock()
			if _, err := td.BasesDeDatos.InsertarId(INSERTAR_TABLA_INCOMPLETOS, datos...); err != nil {
				td.lockIncompletos.Unlock()
				return fmt.Errorf("error al insertar en la tabla auxiliar de incompletos, con error: %v", err)
			}
			td.lockIncompletos.Unlock()
		}
	}

	return nil
}

type updateDependible struct {
	tablaDependiente string
	idDependiente    int64
	key              string
}

func (td *TrackerDependencias) procesoDependible(tabla Tabla, idInsertado int64, hashDatos IntFK) error {
	td.lockDependibles.Lock()
	if _, err := td.BasesDeDatos.InsertarId(INSERTAR_TABLA_DEPENDIENTES, tabla.NombreTabla, idInsertado, hashDatos); err != nil {
		td.lockDependibles.Unlock()
		return fmt.Errorf("error al insertar en dependientes: %s, con error: %v", tabla.NombreTabla, err)
	}
	td.lockDependibles.Unlock()

	td.lockIncompletos.Lock()
	if filas, err := td.BasesDeDatos.Query(QUERY_TABLA_INCOMPLETOS, tabla.NombreTabla, hashDatos); err != nil {
		td.lockIncompletos.Unlock()
		return fmt.Errorf("error al query cuales son los elementos incompletos con tabla: %s, con error: %v", tabla.NombreTabla, err)

	} else {
		td.lockIncompletos.Unlock()

		hayUpdates := false
		for filas.Next() {
			hayUpdates = true

			var update updateDependible
			if err = filas.Scan(&update.tablaDependiente, &update.idDependiente, &update.key); err != nil {
				filas.Close()
				return fmt.Errorf("error al obtener datos de una query de incompletos, con error: %v", err)
			}

			query := fmt.Sprintf("UPDATE %s SET %s = %d WHERE id = %d", update.tablaDependiente, update.key, idInsertado, update.idDependiente)
			lock := td.locksTablas[update.tablaDependiente]
			lock.Lock()
			if err = td.BasesDeDatos.Update(query); err != nil {
				lock.Unlock()
				return fmt.Errorf("error al actualizar %d en tabla %s (proceso Dependible distinto), con error %v", idInsertado, tabla.NombreTabla, err)
			}
			lock.Unlock()

		}
		filas.Close()

		if hayUpdates {
			td.lockIncompletos.Lock()
			if err = td.BasesDeDatos.Eliminar(ELIMINAR_TABLA_INCOMPLETOS, tabla.NombreTabla, hashDatos); err != nil {
				td.lockIncompletos.Unlock()
				return fmt.Errorf("error al eliminar %d en tabla %s, con error %v", idInsertado, tabla.NombreTabla, err)
			}
			td.lockIncompletos.Unlock()
		}

		return nil
	}
}

type updateTodos struct {
	tablaDependiente, key, tablaDestino string
	idDependiente, idInsertado          int64
}

func (td *TrackerDependencias) procesoUltimasActualizaciones() error {
	// Hacer un inner join pora obtener ya de por si los id, sin tener que buscarlos
	td.lockIncompletos.Lock()
	if filas, err := td.BasesDeDatos.Query(QUERY_TODO_INCOMPLETOS); err != nil {
		td.lockIncompletos.Unlock()
		return fmt.Errorf("error al query de todos los elementos de incompletos, con error: %v", err)

	} else {
		td.lockIncompletos.Unlock()
		updates := []updateTodos{}
		for filas.Next() {
			var u updateTodos
			if err = filas.Scan(&u.tablaDependiente, &u.idDependiente, &u.key, &u.tablaDestino, &u.idInsertado); err != nil {
				filas.Close()
				return fmt.Errorf("error al obtener datos de una query de incompletos, con error: %v", err)
			}
			updates = append(updates, u)

		}
		filas.Close()

		for _, u := range updates {
			query := fmt.Sprintf("UPDATE %s SET %s = %d WHERE id = %d", u.tablaDependiente, u.key, u.idInsertado, u.idDependiente)
			td.locksTablas[u.tablaDependiente].Lock()
			if err = td.BasesDeDatos.Update(query); err != nil {
				td.locksTablas[u.tablaDependiente].Unlock()
				return fmt.Errorf("error al actualizar %d en tabla %s, con error %v", u.idInsertado, u.tablaDestino, err)
			}
			td.locksTablas[u.tablaDependiente].Unlock()
		}

		if err = td.BasesDeDatos.Eliminar(ELIMINAR_TODO_INCOMPLETOS); err != nil {
			return fmt.Errorf("error al eliminar el resto de la tabla de incompletos, con error %v", err)
		}

		return nil
	}
}
