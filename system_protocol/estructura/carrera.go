package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const QUERY_CARRERA = "SELECT id FROM carreras WHERE nombre = ?"
const QUERY_CARRERA_PATH = `SELECT res.id FROM (
	SELECT carreras.id, archivos.path FROM archivos INNER JOIN carreras ON archivos.id = carreras.idArchivo
) AS res WHERE res.path = ?`
const INSERTAR_CARRERA = "INSERT INTO carreras (nombre, etapa, tieneCodigoMateria, idArchivo) VALUES (?, ?, ?, ?)"

type ConstructorCarrera struct {
	Nombre            string
	Etapa             Etapa
	TieneCodigo       bool
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorCarrera(nombre string, repEtapa string, tieneCodigo string) (*ConstructorCarrera, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear carrera con error: %v", err)
	} else {
		return &ConstructorCarrera{
			Nombre:            nombre,
			Etapa:             etapa,
			TieneCodigo:       BooleanoODefault(tieneCodigo, false),
			ListaDependencias: l.NewLista[Dependencia](),
		}, nil
	}
}

func (cd *ConstructorCarrera) CargarDependencia(dependencia Dependencia) {
	cd.ListaDependencias.Push(dependencia)
}

func (cd *ConstructorCarrera) CumpleDependencia(id int64) (Cargable, bool) {
	return &Carrera{
		Nombre:            cd.Nombre,
		Etapa:             cd.Etapa,
		TieneCodigo:       cd.TieneCodigo,
		IdArchivo:         id,
		ListaDependencias: cd.ListaDependencias,
	}, true
}

type Carrera struct {
	Nombre            string
	Etapa             Etapa
	TieneCodigo       bool
	IdArchivo         int64
	ListaDependencias *l.Lista[Dependencia]
}

func (c *Carrera) Insertar() []any {
	return []any{
		c.Nombre,
		c.Etapa,
		c.TieneCodigo,
		c.IdArchivo,
	}
}

func (c *Carrera) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- fmt.Sprintf("Insertar Carrera: %s", c.Nombre)
	return Insertar(
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_CARRERA, c.Insertar()...) },
	)
}

func (c *Carrera) ResolverDependencias(id int64) []Cargable {
	cantidadCumple := 0
	cargables := make([]Cargable, c.ListaDependencias.Largo)

	for cumpleDependencia := range c.ListaDependencias.Iterar {
		if cargable, cumple := cumpleDependencia(id); cumple {
			cargables[cantidadCumple] = cargable
			cantidadCumple++
		}
	}

	return cargables[:cantidadCumple]
}
