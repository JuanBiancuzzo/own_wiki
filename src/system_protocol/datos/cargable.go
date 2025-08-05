package datos

import (
	"database/sql"
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
)

type Cargable interface {
	CargarDatos(bdd *b.Bdd, canal chan string) (int64, error)

	ResolverDependencias(id int64) []Cargable
}

func CargableDefault() Cargable {
	var cargable Cargable
	return cargable
}

func ResolverDependencias(id int64, dependencias []Dependencia) []Cargable {
	cantidadCumple := 0
	cargables := make([]Cargable, len(dependencias))

	for _, cumpleDependencia := range dependencias {
		if cargable, cumple := cumpleDependencia(id); cumple {
			cargables[cantidadCumple] = cargable
			cantidadCumple++
		}
	}

	return cargables[:cantidadCumple]
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

func ObtenerDirecto(bdd *b.Bdd, query string, datos ...any) (int64, bool) {
	var id int64
	row := bdd.MySQL.QueryRow(QUERY_PERSONAS, datos...)
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

func InsertarDirecto(bdd *b.Bdd, query string, datos ...any) (int64, error) {
	if filaAfectada, err := bdd.MySQL.Exec(query, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}
