package estructura

// Es una estructura que es capaz de tener dependencias
type Dependible interface {
	CargarDependencia(dependencia Dependencia)
}

type Dependencia func(id int64) (Cargable, bool)
