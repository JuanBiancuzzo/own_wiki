package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const INSERTAR_TEMA_MATERIA = "INSERT INTO temasMateria (nombre, capitulo, parte, idMateria, idArchivo) VALUES (?, ?, ?, ?, ?)"

type TemaMateria struct {
	IdArchivo         *Opcional[int64]
	IdMateria         *Opcional[int64]
	Nombre            string
	Capitulo          int
	Parte             int
	ListaDependencias *l.Lista[Dependencia]
}

func NewTemaMateria(nombre string, capitulo string, parte string) *TemaMateria {
	return &TemaMateria{
		IdArchivo:         NewOpcional[int64](),
		IdMateria:         NewOpcional[int64](),
		Nombre:            nombre,
		Capitulo:          NumeroODefault(capitulo, 1),
		Parte:             NumeroODefault(parte, 0),
		ListaDependencias: l.NewLista[Dependencia](),
	}
}

func (tm *TemaMateria) CrearDependenciaMateria(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		tm.IdMateria.Asignar(id)
		return tm, CumpleAll(tm.IdArchivo, tm.IdMateria)
	})
}

func (tm *TemaMateria) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		tm.IdArchivo.Asignar(id)
		return tm, CumpleAll(tm.IdArchivo, tm.IdMateria)
	})
}

func (tm *TemaMateria) CargarDependencia(dependencia Dependencia) {
	tm.ListaDependencias.Push(dependencia)
}

func (tm *TemaMateria) Insertar() ([]any, error) {
	if idArchivo, existe := tm.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("tema materia no tiene todavia el idArchivo")

	} else if idMateria, existe := tm.IdMateria.Obtener(); !existe {
		return []any{}, fmt.Errorf("tema materia no tiene todavia el idMateria")

	} else {
		return []any{tm.Nombre, tm.Capitulo, tm.Parte, idMateria, idArchivo}, nil
	}
}

func (tm *TemaMateria) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	// canal <- fmt.Sprintf("Insertar Tema Materia: %s", tm.Nombre)
	if datos, err := tm.Insertar(); err != nil {
		return 0, err
	} else {
		return InsertarDirecto(bdd, INSERTAR_TEMA_MATERIA, datos...)
	}
}

func (tm *TemaMateria) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, tm.ListaDependencias.Items())
}
