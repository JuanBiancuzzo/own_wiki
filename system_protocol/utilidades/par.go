package utilidades

type Par[P any, S any] struct {
	Primero P
	Segundo S
}

func NewPar[P any, S any](primero P, segundo S) *Par[P, S] {
	return &Par[P, S]{
		Primero: primero,
		Segundo: segundo,
	}
}
