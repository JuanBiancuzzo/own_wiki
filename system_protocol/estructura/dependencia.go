package estructura

// Es una estructura que es capaz de tener dependencias
type Dependibles interface {
	CargarDependencia(dependencia Dependencia)
}

type Dependencia interface {
	CumpleDependencia(id int64) (Cargable, bool)
}
