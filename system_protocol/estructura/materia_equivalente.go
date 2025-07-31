package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const INSERTAR_MATERIA_EQUIVALENTES = "INSERT INTO materiasEquivalentes (nombre, codigo, idMateria, idArchivo) VALUES (?, ?, ?, ?)"
const QUERY_MATERIA_EQUIVALENTES_PATH = `SELECT res.id FROM (
	SELECT materiasEquivalentes.id, archivos.path FROM archivos INNER JOIN materiasEquivalentes ON archivos.id = materiasEquivalentes.idArchivo
) AS res WHERE res.path = ?`

type ConstructorMateriaEquivalente struct {
	IdArchivo         Opcional[int64]
	IdMateria         Opcional[int64]
	PathMateria       string
	Nombre            string
	Codigo            string
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorMateriaEquivalente(pathMateria string, nombre string, codigo string) *ConstructorMateriaEquivalente {
	return &ConstructorMateriaEquivalente{
		IdArchivo:         NewOpcional[int64](),
		IdMateria:         NewOpcional[int64](),
		PathMateria:       pathMateria,
		Nombre:            nombre,
		Codigo:            codigo,
		ListaDependencias: l.NewLista[Dependencia](),
	}
}
func (cme *ConstructorMateriaEquivalente) CumpleDependencia() (*MateriaEquivalente, bool) {
	if cme.IdArchivo.Esta && cme.IdMateria.Esta {
		return &MateriaEquivalente{
			IdArchivo:         cme.IdArchivo.Valor,
			IdMateria:         cme.IdMateria.Valor,
			Nombre:            cme.Nombre,
			Codigo:            cme.Codigo,
			ListaDependencias: cme.ListaDependencias,
		}, true
	}

	return nil, false
}

func (cme *ConstructorMateriaEquivalente) CumpleDependenciaMateria(id int64) (Cargable, bool) {
	cme.IdMateria.Asignar(id)
	return cme.CumpleDependencia()
}

func (cme *ConstructorMateriaEquivalente) CumpleDependenciaArchivo(id int64) (Cargable, bool) {
	cme.IdArchivo.Asignar(id)
	return cme.CumpleDependencia()
}

func (cme *ConstructorMateriaEquivalente) CargarDependencia(dependencia Dependencia) {
	cme.ListaDependencias.Push(dependencia)
}

type MateriaEquivalente struct {
	IdArchivo         int64
	IdMateria         int64
	Nombre            string
	Codigo            string
	ListaDependencias *l.Lista[Dependencia]
}

func (me *MateriaEquivalente) Insertar() []any {
	return []any{me.Nombre, me.Codigo, me.IdMateria, me.IdArchivo}
}

func (me *MateriaEquivalente) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- fmt.Sprintf("Insertar Materia Equivalentes: %s", me.Nombre)
	return Insertar(
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_MATERIA_EQUIVALENTES, me.Insertar()...) },
	)
}

func (me *MateriaEquivalente) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, me.ListaDependencias)
}
