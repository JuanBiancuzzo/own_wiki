package utilidades

import "sync"

type Worker[R any] struct {
	Id              int
	CanalInput      chan R
	CanalLibre      chan int
	FuncionEjecutar func(R)
	WaitGroupt      *sync.WaitGroup
}

func NewWorker[R any](id int, canalInput chan R, canalLibre chan int, funcionEjecutar func(R), wg *sync.WaitGroup) *Worker[R] {
	canalLibre <- id
	return &Worker[R]{
		Id:              id,
		CanalInput:      canalInput,
		CanalLibre:      canalLibre,
		FuncionEjecutar: funcionEjecutar,
		WaitGroupt:      wg,
	}
}

func (w *Worker[R]) Ejecutar() {
	for dato := range w.CanalInput {
		w.FuncionEjecutar(dato)
		w.CanalLibre <- w.Id
	}
	w.WaitGroupt.Done()
}

func DividirTrabajo[R any](canalInput chan R, cantidadWorkers int, funcionEjecutar func(R), wg *sync.WaitGroup) {
	canalesInput := make([]chan R, cantidadWorkers)
	var waitWorkers sync.WaitGroup

	canalLibre := make(chan int, cantidadWorkers)
	waitWorkers.Add(cantidadWorkers)
	for i := range cantidadWorkers {
		canalesInput[i] = make(chan R, 5)
		worker := NewWorker(i, canalesInput[i], canalLibre, funcionEjecutar, &waitWorkers)
		go worker.Ejecutar()
	}

	for input := range canalInput {
		idWorker := <-canalLibre
		canalesInput[idWorker] <- input
	}

	for i := range cantidadWorkers {
		close(canalesInput[i])
	}

	waitWorkers.Wait()
	if wg != nil {
		wg.Done()
	}
}
