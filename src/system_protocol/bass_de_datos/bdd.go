package bass_de_datos

import (
	"database/sql"
	"fmt"
)

const NOMBRE_BDD = "baseDeDatos.db"
const MAX_CONN = 20
const TABLA_GENERAL = "General"

type Bdd struct {
	archivoBdd string
	conexiones map[string]*conexion
	contador   int
}

func NewBdd(carpetaOutput string, canalMensajes chan string) (*Bdd, error) {
	archivoBdd := fmt.Sprintf("%s/%s", carpetaOutput, NOMBRE_BDD)
	conexiones := make(map[string]*conexion)

	if conexionGeneral, err := newConexion(archivoBdd); err != nil {
		return nil, err

	} else {
		canalMensajes <- "Se conecto correctamente a SQLite"
		conexiones[TABLA_GENERAL] = conexionGeneral

		return &Bdd{
			archivoBdd: archivoBdd,
			conexiones: conexiones,
			contador:   0,
		}, nil
	}

}

func (bdd *Bdd) conexion(tabla string) (*conexion, bool) {
	conn, ok := bdd.conexiones[tabla]
	return conn, ok
}

func (bdd *Bdd) Close() {
	for tabla := range bdd.conexiones {
		bdd.conexiones[tabla].sql.Close()
	}
}

func (bdd *Bdd) CrearTabla(nombreTabla, query string, datos ...any) error {
	if conn, ok := bdd.conexion(TABLA_GENERAL); !ok {
		return fmt.Errorf("no se creo la conexion general")

	} else if _, err := conn.Exec(query, datos...); err != nil {
		return err

	} else {
		bdd.CrearConexionATabla(nombreTabla)
		return nil
	}
}

func (bdd *Bdd) EliminarTabla(nombreTabla string) error {
	if conn, ok := bdd.conexion(TABLA_GENERAL); !ok {
		return fmt.Errorf("no se creo la conexion general")
	} else {
		_, err := conn.Exec(fmt.Sprintf("DROP TABLE %s", nombreTabla))
		return err
	}
}

func (bdd *Bdd) CrearConexionATabla(tabla string) error {
	if _, ok := bdd.conexiones[tabla]; ok {
		return fmt.Errorf("ya existe la conexion a esa tabla")
	}

	if conexionATabla, err := newConexion(bdd.archivoBdd); err != nil {
		return err

	} else {
		bdd.conexiones[tabla] = conexionATabla
		return nil
	}
}

func (bdd *Bdd) Existe(tabla, query string, datos ...any) (bool, error) {
	lectura := make([]any, len(datos))
	fila := bdd.QueryRow(tabla, query, datos...)

	if err := fila.Scan(lectura...); err != nil {
		return false, nil
	}

	return true, nil
}

func (bdd *Bdd) Obtener(tabla, query string, datos ...any) (int64, error) {
	var id int64
	fila := bdd.QueryRow(tabla, query, datos...)
	if fila == nil {
		return id, fmt.Errorf("error al obtener query con tabla '%s'", tabla)
	}

	if err := fila.Scan(&id); err != nil {
		return id, fmt.Errorf("error al intentar query la bdd, con error: %v", err)
	}

	return id, nil
}

func (bdd *Bdd) Insertar(tabla, query string, datos ...any) (int64, error) {
	if filaAfectada, err := bdd.Exec(tabla, query, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}

func (bdd *Bdd) ObtenerOInsertar(tabla, queryObtener, queryInsertar string, datos ...any) (int64, error) {
	if id, err := bdd.Obtener(tabla, queryObtener, datos...); err == nil {
		return id, nil
	}
	return bdd.Insertar(tabla, queryInsertar, datos...)
}

func (bdd *Bdd) Exec(tabla string, query string, datos ...any) (sql.Result, error) {
	if conexionTabla, ok := bdd.conexion(tabla); !ok {
		var res sql.Result
		return res, fmt.Errorf("no hay tabla %s", tabla)

	} else {
		return conexionTabla.Exec(query, datos...)
	}
}

func (bdd *Bdd) QueryRow(tabla string, query string, datos ...any) *sql.Row {
	if conexionTabla, ok := bdd.conexion(tabla); !ok {
		return nil

	} else {
		return conexionTabla.QueryRow(query, datos...)
	}
}

func (bdd *Bdd) Query(tabla string, query string, datos ...any) (filasSQL, error) {
	if conexionTabla, ok := bdd.conexion(tabla); !ok {
		return filasSQL{}, fmt.Errorf("no hay tabla %s", tabla)

	} else {
		return conexionTabla.Query(query, datos...)
	}
}
