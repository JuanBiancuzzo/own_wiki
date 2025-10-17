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

func DividirTrabajo[R any](canalInput chan R, cantidadWorkers int, funcionEjecutar func(R), wg *sync.WaitGroup) {
	canalesInput := make([]chan R, cantidadWorkers)
	var waitWorkers sync.WaitGroup

	waitWorkers.Add(cantidadWorkers)
	for i := range cantidadWorkers {
		canalesInput[i] = make(chan R, 5)
		worker := NewWorker(canalesInput[i], funcionEjecutar, &waitWorkers)
		go worker.Ejecutar()
	}

	contador := 0
	for input := range canalInput {
		canalesInput[contador] <- input
		contador = (contador + 1) % cantidadWorkers
	}

	for i := range cantidadWorkers {
		close(canalesInput[i])
	}

	waitWorkers.Wait()
	if wg != nil {
		wg.Done()
	}
}
