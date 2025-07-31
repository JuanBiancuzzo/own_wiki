package estructura

type Opcional[T any] struct {
	Valor T
	Esta  bool
}

func NewOpcional[T any]() Opcional[T] {
	var valor T
	return Opcional[T]{
		Valor: valor,
		Esta:  false,
	}
}

func (o Opcional[T]) Asignar(valor T) {
	o.Valor = valor
	o.Esta = true
}

func (o Opcional[T]) Obtener() (T, bool) {
	return o.Valor, o.Esta
}
