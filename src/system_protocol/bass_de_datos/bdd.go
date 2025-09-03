package bass_de_datos

import (
	"fmt"
)

const NOMBRE_BDD = "baseDeDatos.db"

type Bdd struct {
	archivoBdd string
	pool       *poolConexiones
	contador   int
}

func NewBdd(carpetaOutput string, canalMensajes chan string) (*Bdd, error) {
	archivoBdd := fmt.Sprintf("%s/%s", carpetaOutput, NOMBRE_BDD)

	if pool, err := newPoolConexiones(archivoBdd); err != nil {
		return nil, err

	} else {
		canalMensajes <- "Se conecto correctamente a SQLite"

		return &Bdd{
			archivoBdd: archivoBdd,
			pool:       pool,
			contador:   0,
		}, nil
	}

}

func (bdd *Bdd) Close() {
	bdd.pool.Close()
}

func (bdd *Bdd) CrearTabla(query string, datos ...any) error {
	if conn, err := bdd.pool.Conexion(); err != nil {
		return err

	} else {
		_, err := conn.Exec(query, datos...)
		return err
	}
}

func (bdd *Bdd) EliminarTabla(nombreTabla string) error {
	if conn, err := bdd.pool.Conexion(); err != nil {
		return err

	} else {
		_, err := conn.Exec(fmt.Sprintf("DROP TABLE %s", nombreTabla))
		return err
	}
}

func (bdd *Bdd) Existe(query string, datos ...any) (bool, error) {
	lectura := make([]any, len(datos))
	fila := bdd.QueryRow(query, datos...)

	if err := fila.Scan(lectura...); err != nil {
		return false, nil
	}

	return true, nil
}

func (bdd *Bdd) Obtener(query string, datos ...any) (int64, error) {
	var id int64
	fila := bdd.QueryRow(query, datos...)
	if fila == nil {
		return id, fmt.Errorf("error al obtener query")
	}

	if err := fila.Scan(&id); err != nil {
		return id, fmt.Errorf("error al intentar query la bdd, con error: %v", err)
	}

	return id, nil
}

func (bdd *Bdd) Insertar(query string, datos ...any) (int64, error) {
	if filaAfectada, err := bdd.Exec(query, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query (ejecutando exec), con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}

func (bdd *Bdd) ObtenerOInsertar(queryObtener, queryInsertar string, datos ...any) (int64, error) {
	if id, err := bdd.Obtener(queryObtener, datos...); err == nil {
		return id, nil
	}
	return bdd.Insertar(queryInsertar, datos...)
}

func (bdd *Bdd) Exec(query string, datos ...any) (resultadoSQL, error) {
	if conn, err := bdd.pool.Conexion(); err != nil {
		return resultadoSQL{}, err

	} else {
		return conn.Exec(query, datos...)
	}
}

func (bdd *Bdd) QueryRow(query string, datos ...any) *filaSQL {
	if conn, err := bdd.pool.Conexion(); err != nil {
		return nil

	} else {
		return conn.QueryRow(query, datos...)
	}
}

func (bdd *Bdd) Query(query string, datos ...any) (filasSQL, error) {
	if conn, err := bdd.pool.Conexion(); err != nil {
		return filasSQL{}, nil

	} else {
		return conn.Query(query, datos...)
	}
}
