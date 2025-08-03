package utilidades

import "sync"

type Worker[R any] struct {
	CanalInput      chan R
	FuncionEjecutar func(R)
	WaitGroupt      *sync.WaitGroup
}

func NewWorker[R any](canalInput chan R, funcionEjecutar func(R), wg *sync.WaitGroup) *Worker[R] {
	return &Worker[R]{
		CanalInput:      canalInput,
		FuncionEjecutar: funcionEjecutar,
		WaitGroupt:      wg,
	}
}

func (w *Worker[R]) Ejecutar() {
	for dato := range w.CanalInput {
		w.FuncionEjecutar(dato)
	}
	w.WaitGroupt.Done()
}
