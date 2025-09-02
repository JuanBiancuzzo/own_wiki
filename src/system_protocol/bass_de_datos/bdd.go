package bass_de_datos

import (
	"database/sql"
	"fmt"
)

const NOMBRE_BDD = "baseDeDatos.db"
const MAX_CONN = 20

type Bdd struct {
	conexiones []*conexion
	contador   int
}

func NewBdd(carpetaOutput string, canalMensajes chan string) (*Bdd, error) {
	archivoBdd := fmt.Sprintf("%s/%s", carpetaOutput, NOMBRE_BDD)
	conexiones := make([]*conexion, MAX_CONN)
	var err error

	for i := range MAX_CONN {
		if conexiones[i], err = newConexion(archivoBdd); err != nil {
			return nil, fmt.Errorf("en la conexion n°%d se tuvo: %v", i+1, err)
		}
	}

	canalMensajes <- "Se conecto correctamente a SQLite"

	return &Bdd{
		conexiones: conexiones,
		contador:   0,
	}, nil
}

func (bdd *Bdd) conexion() *conexion {
	conexion := bdd.conexiones[bdd.contador]
	bdd.contador = (bdd.contador + 1) % MAX_CONN
	return conexion
}

func (bdd *Bdd) Close() {
	for _, conexion := range bdd.conexiones {
		conexion.sql.Close()
	}
}

func (bdd *Bdd) CrearTabla(query string, datos ...any) error {
	_, err := bdd.Exec(query, datos...)
	return err
}

func (bdd *Bdd) EliminarTabla(nombreTabla string) error {
	_, err := bdd.Exec(fmt.Sprintf("DROP TABLE %s", nombreTabla))
	return err
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

	if err := fila.Scan(&id); err != nil {
		return id, fmt.Errorf("error al intentar query la bdd, con error: %v", err)
	}

	return id, nil
}

func (bdd *Bdd) Insertar(query string, datos ...any) (int64, error) {
	if filaAfectada, err := bdd.Exec(query, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

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

func (bdd *Bdd) Exec(query string, datos ...any) (sql.Result, error) {
	return bdd.conexion().Exec(query, datos...)
}

func (bdd *Bdd) QueryRow(query string, datos ...any) *sql.Row {
	return bdd.conexion().QueryRow(query, datos...)
}

func (bdd *Bdd) Query(query string, datos ...any) (filasSQL, error) {
	return bdd.conexion().Query(query, datos...)
}
