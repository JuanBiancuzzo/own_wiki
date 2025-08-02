package estructura

type Opcional[T any] struct {
	Valor T
	Esta  bool
}

func NewOpcional[T any]() *Opcional[T] {
	var valorDefault T
	return &Opcional[T]{
		Valor: valorDefault,
		Esta:  false,
	}
}

func (o *Opcional[T]) Asignar(valor T) {
	o.Valor = valor
	o.Esta = true
}

func (o *Opcional[T]) Obtener() (T, bool) {
	return o.Valor, o.Esta
}

func CumpleAll[T any](valores ...*Opcional[T]) bool {
	for _, valor := range valores {
		if !valor.Esta {
			return false
		}
	}
	return true
}
