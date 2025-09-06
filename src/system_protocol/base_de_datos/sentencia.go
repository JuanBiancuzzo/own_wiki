package base_de_datos

import (
	"database/sql"
	"fmt"
)

type Sentencia struct {
	sentencia *sql.Stmt
}

func NewSentencia(bdd *sql.DB, query string) (Sentencia, error) {
	if sentencia, err := bdd.Prepare(query); err != nil {
		return Sentencia{}, err
	} else {
		return Sentencia{
			sentencia: sentencia,
		}, nil
	}
}

func (s Sentencia) Obtener(datos ...any) (int64, error) {
	var id int64
	fila := s.QueryRow(datos...)
	if fila == nil {
		return id, fmt.Errorf("error al obtener query")
	}

	if err := fila.Scan(&id); err != nil {
		return id, fmt.Errorf("error al intentar query la bdd, con error: %v", err)
	}

	return id, nil
}

func (s Sentencia) InsertarId(datos ...any) (int64, error) {
	if filaAfectada, err := s.exec(datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query (ejecutando exec), con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}

func (s Sentencia) Update(datos ...any) error {
	_, err := s.exec(datos...)
	return err
}

func (s Sentencia) Eliminar(datos ...any) error {
	_, err := s.exec(datos...)
	return err
}

func (s Sentencia) QueryRow(datos ...any) *sql.Row {
	return s.sentencia.QueryRow(datos...)
}

func (s Sentencia) Query(datos ...any) (*sql.Rows, error) {
	return s.sentencia.Query(datos...)
}

func (s Sentencia) exec(datos ...any) (sql.Result, error) {
	return s.sentencia.Exec(datos...)
}
