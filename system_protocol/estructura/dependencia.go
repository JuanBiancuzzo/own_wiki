package estructura

// Es una estructura que es capaz de tener dependencias
type Dependible interface {
	CargarDependencia(dependencia Dependencia)
}

type Dependencia func(id int64) (Cargable, bool)

type FnVincular func(dependible Dependible)

type TipoDependible byte

const (
	DEP_ARCHIVO = iota
	DEP_CARRERA
	DEP_MATERIA
	DEP_MATERIA_EQUIVALENTE
	DEP_TEMA_MATERIA
	DEP_CURSO
	DEP_CURSO_PRESENCIAL
	DEP_TEMA_CURSO
)
