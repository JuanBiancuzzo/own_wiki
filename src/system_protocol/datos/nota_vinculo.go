package datos

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	u "own_wiki/system_protocol/utilidades"
)

type TipoNota string

const (
	TN_FACULTAD      = "Facultad"
	TN_COLECCION     = "Coleccion"
	TN_CURSO         = "Curso"
	TN_INVESTIGACION = "Investigacion"
	TN_PROYECTO      = "Proyecto"
)

const INSERTAR_NOTA_VINCULO = "INSERT INTO notasVinculo (idNota, idVinculo, tipoVinculo) VALUES (?, ?, ?)"

type NotaVinculo struct {
	IdNota    *u.Opcional[int64]
	IdVinculo *u.Opcional[int64]
	Tipo      TipoNota
}

func NewNotaVinculo(tipo TipoNota) *NotaVinculo {
	return &NotaVinculo{
		IdNota:    u.NewOpcional[int64](),
		IdVinculo: u.NewOpcional[int64](),
		Tipo:      tipo,
	}
}

func (nv *NotaVinculo) CrearDependenciaVinculo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		nv.IdVinculo.Asignar(id)
		return nv, u.CumpleAll(nv.IdNota, nv.IdVinculo)
	})
}

func (nv *NotaVinculo) CrearDependenciaNota(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		nv.IdNota.Asignar(id)
		return nv, u.CumpleAll(nv.IdNota, nv.IdVinculo)
	})
}

func (nv *NotaVinculo) Insertar() ([]any, error) {
	if idNota, existe := nv.IdNota.Obtener(); !existe {
		return []any{}, fmt.Errorf("nota vinculo no tiene todavia el idNota")

	} else if idVinculo, existe := nv.IdVinculo.Obtener(); !existe {
		return []any{}, fmt.Errorf("nota vinculo no tiene todavia el idVinculo")

	} else {
		return []any{idNota, idVinculo, nv.Tipo}, nil
	}
}

func (nv *NotaVinculo) CargarDatos(bdd *b.Bdd, canal chan string) (int64, error) {
	// canal <- "Insertar Nota vinculante"
	if datos, err := nv.Insertar(); err != nil {
		return 0, err
	} else {
		return InsertarDirecto(bdd, INSERTAR_NOTA_VINCULO, datos...)
	}
}

func (nv *NotaVinculo) ResolverDependencias(id int64) []Cargable {
	return []Cargable{}
}
