package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const INSERTAR_TEMA_MATERIA = "INSERT INTO temasMateria (nombre, capitulo, parte, idMateria, idArchivo) VALUES (?, ?, ?, ?, ?)"

// Corregir antes de usar
/* const QUERY_RESUMEN_MATERIA_PATH = `SELECT res.id FROM (
	SELECT materiasEquivalentes.id, archivos.path FROM archivos INNER JOIN materiasEquivalentes ON archivos.id = materiasEquivalentes.idArchivo
) AS res WHERE res.path = ?`*/

type ConstructorTemaMateria struct {
	IdArchivo         *Opcional[int64]
	IdMateria         *Opcional[int64]
	Nombre            string
	Capitulo          int
	Parte             int
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorTemaMateria(nombre string, capitulo string, parte string) *ConstructorTemaMateria {
	return &ConstructorTemaMateria{
		IdArchivo:         NewOpcional[int64](),
		IdMateria:         NewOpcional[int64](),
		Nombre:            nombre,
		Capitulo:          NumeroODefault(capitulo, 1),
		Parte:             NumeroODefault(parte, 0),
		ListaDependencias: l.NewLista[Dependencia](),
	}
}

func (ctm *ConstructorTemaMateria) CumpleDependencia() (*TemaMateria, bool) {
	if idArchivo, existe := ctm.IdArchivo.Obtener(); !existe {
		return nil, false

	} else if idMateria, existe := ctm.IdMateria.Obtener(); !existe {
		return nil, false

	} else {
		return &TemaMateria{
			IdArchivo:         idArchivo,
			IdMateria:         idMateria,
			Nombre:            ctm.Nombre,
			Capitulo:          ctm.Capitulo,
			Parte:             ctm.Parte,
			ListaDependencias: ctm.ListaDependencias,
		}, true
	}
}

func (ctm *ConstructorTemaMateria) CrearDependenciaMateria(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		ctm.IdMateria.Asignar(id)
		return ctm.CumpleDependencia()
	})
}

func (ctm *ConstructorTemaMateria) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		ctm.IdArchivo.Asignar(id)
		return ctm.CumpleDependencia()
	})
}

func (ctm *ConstructorTemaMateria) CargarDependencia(dependencia Dependencia) {
	ctm.ListaDependencias.Push(dependencia)
}

type TemaMateria struct {
	IdArchivo         int64
	IdMateria         int64
	Nombre            string
	Capitulo          int
	Parte             int
	ListaDependencias *l.Lista[Dependencia]
}

func (tm *TemaMateria) Insertar() []any {
	return []any{tm.Nombre, tm.Capitulo, tm.Parte, tm.IdMateria, tm.IdArchivo}
}

func (tm *TemaMateria) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- fmt.Sprintf("Insertar Resumen Materia: %s", tm.Nombre)
	return Insertar(
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_TEMA_MATERIA, tm.Insertar()...) },
	)
}

func (tm *TemaMateria) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, tm.ListaDependencias.Items())
}
