package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const QUERY_CARRERA = "SELECT id FROM carreras WHERE nombre = ?"
const INSERTAR_CARRERA = "INSERT INTO carreras (nombre, etapa, tieneCodigoMateria, idArchivo) VALUES (?, ?, ?, ?)"

type Carrera struct {
	Nombre            string
	Etapa             Etapa
	TieneCodigo       bool
	IdArchivo         *Opcional[int64]
	ListaDependencias *l.Lista[Dependencia]
}

func NewCarrera(nombre string, repEtapa string, tieneCodigo string) (*Carrera, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear carrera con error: %v", err)
	} else {
		return &Carrera{
			Nombre:            nombre,
			Etapa:             etapa,
			TieneCodigo:       BooleanoODefault(tieneCodigo, false),
			IdArchivo:         NewOpcional[int64](),
			ListaDependencias: l.NewLista[Dependencia](),
		}, nil
	}
}

func (c *Carrera) CargarDependencia(dependencia Dependencia) {
	c.ListaDependencias.Push(dependencia)
}

func (c *Carrera) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		c.IdArchivo.Asignar(id)
		return c, true
	})
}

func (c *Carrera) Insertar() ([]any, error) {
	if idArchivo, existe := c.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("carrera no tiene todavia el idArchivo")
	} else {
		return []any{c.Nombre, c.Etapa, c.TieneCodigo, idArchivo}, nil
	}
}

func (c *Carrera) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- fmt.Sprintf("Insertar Carrera: %s", c.Nombre)
	if datos, err := c.Insertar(); err != nil {
		return 0, err
	} else {
		return InsertarDirecto(bdd, INSERTAR_CARRERA, datos...)
	}
}

func (c *Carrera) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, c.ListaDependencias.Items())
}
