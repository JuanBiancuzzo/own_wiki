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
	DEP_NOTA
)

func TipoDependible2String(tipo TipoDependible) string {
	switch tipo {
	case DEP_ARCHIVO:
		return "Dependencia archivo"
	case DEP_CARRERA:
		return "Dependencia carrera"
	case DEP_MATERIA:
		return "Dependencia materia"
	case DEP_MATERIA_EQUIVALENTE:
		return "Dependencia materia equivalente"
	case DEP_TEMA_MATERIA:
		return "Dependencia tema de una materia"
	case DEP_CURSO:
		return "Dependencia curso online"
	case DEP_CURSO_PRESENCIAL:
		return "Dependencia curso presencial"
	case DEP_TEMA_CURSO:
		return "Dependencia tema de un curso"
	case DEP_NOTA:
		return "Dependencia nota"
	}

	return "[[ERROR]]"
}
