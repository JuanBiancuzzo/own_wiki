package datos

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const FORMATO_DIA = "YYYY-MM-DD"
const INSERTAR_NOTA = "INSERT INTO notas (nombre, etapa, dia, idArchivo) VALUES (?, ?, ?, ?)"

type Nota struct {
	IdArchivo         *Opcional[int64]
	Nombre            string
	Etapa             Etapa
	Dia               string
	ListaDependencias *l.Lista[Dependencia]
}

func NewNota(nombre string, repEtapa string, dia string) *Nota {
	return &Nota{
		IdArchivo:         NewOpcional[int64](),
		Nombre:            nombre,
		Etapa:             EtapaODefault(repEtapa, ETAPA_SIN_EMPEZAR),
		Dia:               dia,
		ListaDependencias: l.NewLista[Dependencia](),
	}
}

func (n *Nota) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		n.IdArchivo.Asignar(id)
		return n, true
	})
}

func (n *Nota) CargarDependencia(dependencia Dependencia) {
	n.ListaDependencias.Push(dependencia)
}

func (n *Nota) Insertar() ([]any, error) {
	if idArchivo, existe := n.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("materia no tiene todavia el idArchivo")

	} else {
		return []any{n.Nombre, n.Etapa, n.Dia, idArchivo}, nil
	}
}

func (n *Nota) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	// canal <- fmt.Sprintf("Insertar Nota: %s", n.Nombre)
	if datos, err := n.Insertar(); err != nil {
		return 0, err
	} else {
		return InsertarDirecto(bdd, INSERTAR_NOTA, datos...)
	}
}

func (n *Nota) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, n.ListaDependencias.Items())
}
