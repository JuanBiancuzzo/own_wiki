package estructura

import (
	"database/sql"
	"fmt"
)

type Cargable interface {
	CargarDatos(bdd *sql.DB, canal chan string) bool
}

func ObtenerOInsertar(query func() *sql.Row, insert func() (sql.Result, error)) (int64, error) {
	if id, seObtuvo := Obtener(query); seObtuvo {
		return id, nil
	}

	return Insertar(insert)
}

func Obtener(query func() *sql.Row) (int64, bool) {
	var id int64
	row := query()
	if err := row.Scan(&id); err != nil {
		return 0, false
	}
	return id, true
}

func Insertar(insert func() (sql.Result, error)) (int64, error) {
	if filaAfectada, err := insert(); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}
