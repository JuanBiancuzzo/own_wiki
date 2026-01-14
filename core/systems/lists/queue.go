package lists

type Queue[T any] struct {
	List *List[T]
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		List: NewList[T](),
	}
}

func NewColaConCapacidad[T any](capacity uint32) *Queue[T] {
	return &Queue[T]{
		List: NewListWithCapacity[T](capacity),
	}
}

func (q *Queue[T]) Queue(elemento T) {
	q.List.Push(elemento)
}

func (q *Queue[T]) Dequeue() (T, error) {
	return q.List.PopFrom(0)
}

func (q *Queue[T]) Empty() bool {
	return q.List.Empty()
}

func (q *Queue[T]) Iterate(yield func(T) bool) {
	for !q.Empty() {
		if element, err := q.Dequeue(); err != nil {
			return

		} else if !yield(element) {
			return
		}
	}
}
