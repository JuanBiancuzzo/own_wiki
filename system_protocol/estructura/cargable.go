package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

type Cargable interface {
	CargarDatos(bdd *sql.DB, canal chan string) (int64, error)

	ResolverDependencias(id int64) []Cargable
}

func CargableDefault() Cargable {
	var cargable Cargable
	return cargable
}

func ResolverDependencias(id int64, dependencias *l.Lista[Dependencia]) []Cargable {
	cantidadCumple := 0
	cargables := make([]Cargable, dependencias.Largo)

	for cumpleDependencia := range dependencias.Iterar {
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

func Insertar(insert func() (sql.Result, error)) (int64, error) {
	if filaAfectada, err := insert(); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}
