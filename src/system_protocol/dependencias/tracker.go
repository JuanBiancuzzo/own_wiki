package dependencias

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
	"sync"

	b "own_wiki/system_protocol/bass_de_datos"
	u "own_wiki/system_protocol/utilidades"
)

const TABLA_DEPENDIBLES = `CREATE TABLE IF NOT EXISTS aux_dependibles (
	nombreTabla ENUM(%s) NOT NULL,
	hashDatos   BIGINT,
	idDatos     INT
);`
const TABLA_INCOMPLETOS = `CREATE TABLE IF NOT EXISTS aux_incompletos (
	tablaDependiente 	ENUM(%s) NOT NULL,
	idDependiente   	INT,
	keyAlId 			VARCHAR(255),
	tablaDestino 		ENUM(%s) NOT NULL,
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
	RegistrarTablas map[string]DescripcionTabla
	Hash            *Hash

	lockIncompletos *sync.Mutex
	lockTablas      map[string]*sync.Mutex
}

func NewTrackerDependencias(bdd *b.Bdd, tablas []DescripcionTabla, canalMensajes chan string) (*TrackerDependencias, error) {
	nombreTablas := []string{}
	tablasProcesar := u.NewCola[DescripcionTabla]()
	registrarTablas := make(map[string]DescripcionTabla)
	lockTablas := make(map[string]*sync.Mutex)

	for _, descripcion := range tablas {
		nombreTablas = append(nombreTablas, fmt.Sprintf("\"%s\"", descripcion.NombreTabla))

		registrarTablas[descripcion.NombreTabla] = descripcion
		var lock sync.Mutex
		lockTablas[descripcion.NombreTabla] = &lock

		if !EsTipoDependible(descripcion.TipoTabla) {
			tablasProcesar.Encolar(descripcion)
		}
	}
	enumTablas := strings.Join(nombreTablas, ", ")

	if err := bdd.CrearTabla(fmt.Sprintf(TABLA_DEPENDIBLES, enumTablas)); err != nil {
		return nil, err

	} else if err := bdd.CrearTabla(fmt.Sprintf(TABLA_INCOMPLETOS, enumTablas, enumTablas)); err != nil {
		return nil, err

	} else if tablasOrdenadas, err := crearTablas(tablasProcesar); err != nil {
		return nil, err

	} else {
		for _, tabla := range slices.Backward(tablasOrdenadas) {
			if err := bdd.EliminarTabla(tabla.NombreTabla); err != nil {
				canalMensajes <- fmt.Sprintf("error al eliminar tabla %s con error: %v", tabla.NombreTabla, err)
				continue
			}
		}

		canalMensajes <- "Orden final de cargado:"
		for _, tabla := range tablasOrdenadas {
			canalMensajes <- "Tabla: " + tabla.NombreTabla

			if err := tabla.CrearTablaRelajada(bdd); err != nil {
				return nil, fmt.Errorf("error al crear tablas relajadas, especificamente en %s, con error: %v", tabla.NombreTabla, err)
			}
		}
	}

	var lock sync.Mutex
	return &TrackerDependencias{
		BasesDeDatos:    bdd,
		RegistrarTablas: registrarTablas,
		Hash:            NewHash(),

		lockIncompletos: &lock,
		lockTablas:      lockTablas,
	}, nil
}

func crearTablas(tablasProcesar *u.Cola[DescripcionTabla]) ([]DescripcionTabla, error) {
	// Creando las tablas relajadas
	var tablasOrdenadas []DescripcionTabla = []DescripcionTabla{}
	for tabla := range tablasProcesar.DesencolarIterativamente {
		nombreTabla := tabla.NombreTabla

		for i, tablaExistente := range tablasOrdenadas {
			if tablaExistente.NombreTabla == nombreTabla {
				tablasOrdenadas = append(tablasOrdenadas[:i], tablasOrdenadas[i+1:]...)
				break
			}
		}

		tablasOrdenadas = append([]DescripcionTabla{tabla}, tablasOrdenadas...)
		for _, tablaDependible := range tabla.ObtenerDependencias {
			tablasProcesar.Encolar(tablaDependible)
		}
	}

	return tablasOrdenadas, nil
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

	if err := td.BasesDeDatos.EliminarTabla("aux_dependibles"); err != nil {
		return fmt.Errorf("error al eliminar tabla auxiliar dependibles, con error: %v", err)

	} else if err = td.BasesDeDatos.EliminarTabla("aux_incompletos"); err != nil {
		return fmt.Errorf("error al eliminar tabla auxiliar incompletos, con error: %v", err)
	}

	return nil
}

func (td *TrackerDependencias) Cargar(nombreTabla string, datosIngresados ConjuntoDato) error {
	if _, ok := td.RegistrarTablas[nombreTabla]; !ok {
		return fmt.Errorf("de alguna forma estas cargando en una tabla no registrada")
	}
	tabla := td.RegistrarTablas[nombreTabla]

	td.lockTablas[nombreTabla].Lock()
	if existe, err := tabla.Existe(td.BasesDeDatos, datosIngresados); err != nil {
		td.lockTablas[nombreTabla].Unlock()
		return err

	} else if existe {
		td.lockTablas[nombreTabla].Unlock()
		return nil
	}
	id, err := tabla.Insertar(td.BasesDeDatos, datosIngresados)
	td.lockTablas[nombreTabla].Unlock()

	if err != nil {
		return err
	}

	if EsTipoDependiente(tabla.TipoTabla) {
		fKeys, err := tabla.CrearForeignKey(td.Hash, datosIngresados)
		if err != nil {
			return err
		}

		if err := td.procesoDependiente(tabla, id, fKeys); err != nil {
			return fmt.Errorf("error al verificar o actualizar el elemnto en la tabla tabla %s, con id: %d, con error: %v", tabla.NombreTabla, id, err)
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

func (td *TrackerDependencias) procesoDependiente(tabla DescripcionTabla, idInsertado int64, fKeys []ForeignKey) error {
	for _, fKey := range fKeys {
		// Vemos si ya fue insertado la dependencia
		if id, err := td.BasesDeDatos.Obtener(QUERY_TABLA_DEPENDIENTES, fKey.TablaDestino, fKey.HashDatosDestino); err == nil {
			// Si fueron insertados, por lo que actualizamos la tabla
			query := fmt.Sprintf("UPDATE %s SET %s = %d WHERE id = %d", tabla.NombreTabla, fKey.Clave, id, idInsertado)
			if _, err = td.BasesDeDatos.MySQL.Exec(query); err != nil {
				return fmt.Errorf("error al actualizar %d en tabla %s, con error %v", idInsertado, tabla.NombreTabla, err)
			}

		} else {
			td.lockIncompletos.Lock()

			// Como no fue insertada, tenemos que guardar la información para que se carge correctamente la dependencia
			datos := []any{tabla.NombreTabla, idInsertado, fKey.Clave, fKey.TablaDestino, fKey.HashDatosDestino}
			if _, err := td.BasesDeDatos.Insertar(INSERTAR_TABLA_INCOMPLETOS, datos...); err != nil {
				td.lockIncompletos.Unlock()
				return fmt.Errorf("error al insertar en la tabla auxiliar de incompletos, con error: %v", err)

			} else {
				td.lockIncompletos.Unlock()
			}
		}
	}

	return nil
}

func (td *TrackerDependencias) procesoDependible(tabla DescripcionTabla, idInsertado int64, hashDatos IntFK) error {
	if _, err := td.BasesDeDatos.Insertar(INSERTAR_TABLA_DEPENDIENTES, tabla.NombreTabla, idInsertado, hashDatos); err != nil {
		return fmt.Errorf("error al insertar en dependientes: %s, con error: %v", tabla.NombreTabla, err)
	}

	td.lockIncompletos.Lock()
	if filas, err := td.BasesDeDatos.MySQL.Query(QUERY_TABLA_INCOMPLETOS, tabla.NombreTabla, hashDatos); err != nil {
		td.lockIncompletos.Unlock()
		return fmt.Errorf("error al query cuales son los elementos incompletos con tabla: %s, con error: %v", tabla.NombreTabla, err)

	} else {
		defer filas.Close()

		hayFilasAfectadas := false
		for filas.Next() {
			hayFilasAfectadas = true
			var tablaDependiente string
			var idDependiente int64
			var key string

			if err = filas.Scan(&tablaDependiente, &idDependiente, &key); err != nil {
				td.lockIncompletos.Unlock()
				return fmt.Errorf("error al obtener datos de una query de incompletos, con error: %v", err)
			}

			query := fmt.Sprintf("UPDATE %s SET %s = %d WHERE id = %d", tablaDependiente, key, idInsertado, idDependiente)
			if _, err = td.BasesDeDatos.MySQL.Exec(query); err != nil {
				td.lockIncompletos.Unlock()
				return fmt.Errorf("error al actualizar %d en tabla %s, con error %v", idInsertado, tabla.NombreTabla, err)
			}
		}

		if hayFilasAfectadas {
			if _, err = td.BasesDeDatos.MySQL.Exec(ELIMINAR_TABLA_INCOMPLETOS, tabla.NombreTabla, hashDatos); err != nil {
				td.lockIncompletos.Unlock()
				return fmt.Errorf("error al eliminar %d en tabla %s, con error %v", idInsertado, tabla.NombreTabla, err)
			}
		}

		td.lockIncompletos.Unlock()
		return nil
	}
}

func (td *TrackerDependencias) procesoUltimasActualizaciones() error {
	// Hacer un inner join pora obtener ya de por si los id, sin tener que buscarlos
	if filas, err := td.BasesDeDatos.MySQL.Query(QUERY_TODO_INCOMPLETOS); err != nil {
		return fmt.Errorf("error al query de todos los elementos de incompletos, con error: %v", err)

	} else {
		defer filas.Close()

		var tablaDependiente, key, tablaDestino string
		var idDependiente, idInsertado int64

		for filas.Next() {
			if err = filas.Scan(&tablaDependiente, &idDependiente, &key, &tablaDestino, &idInsertado); err != nil {
				return fmt.Errorf("error al obtener datos de una query de incompletos, con error: %v", err)
			}

			query := fmt.Sprintf("UPDATE %s SET %s = %d WHERE id = %d", tablaDependiente, key, idInsertado, idDependiente)
			if _, err = td.BasesDeDatos.MySQL.Exec(query); err != nil {
				return fmt.Errorf("error al actualizar %d en tabla %s, con error %v", idInsertado, tablaDestino, err)
			}
		}

		if _, err = td.BasesDeDatos.MySQL.Exec(ELIMINAR_TODO_INCOMPLETOS); err != nil {
			return fmt.Errorf("error al eliminar el resto de la tabla de incompletos, con error %v", err)
		}

		return nil
	}
}
