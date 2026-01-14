package lists

import (
	"fmt"
)

type List[T any] struct {
	Elements []T
	Capacity uint32
	Length   uint32
}

const DEFAULT_CAPACITY = 2
const CAPACITY_MULTIPLICATION = 2

func NewList[T any]() *List[T] {
	return NewListWithCapacity[T](DEFAULT_CAPACITY)
}

func NewListWithCapacity[T any](capacity uint32) *List[T] {
	elements := make([]T, capacity)
	for i := 0; i < int(capacity); i++ {
		var element T
		elements[i] = element
	}

	return &List[T]{
		Elements: elements,
		Capacity: capacity,
		Length:   0,
	}
}

func (l *List[T]) extend() {
	newCapacity := l.Capacity * CAPACITY_MULTIPLICATION
	newElements := make([]T, newCapacity)

	for i := 0; i < int(l.Length); i++ {
		newElements[i] = l.Elements[i]
	}

	l.Elements = newElements
	l.Capacity = newCapacity
}

func (l *List[T]) Push(element T) {
	if l.Length == l.Capacity {
		l.extend()
	}

	l.Elements[l.Length] = element
	l.Length++
}

func (l *List[T]) AddIn(element T, in uint32) error {
	if l.Length < in {
		return fmt.Errorf("Failed to add element %v in %d, with error: index is greater than list length", element, in)
	}

	if l.Length == l.Capacity {
		l.extend()
	}

	elementToReplace := element
	for i := int(in); i <= int(l.Length); i++ {
		temp := l.Elements[i]
		l.Elements[i] = elementToReplace
		elementToReplace = temp
	}

	l.Length++
	return nil
}

func (l *List[T]) GetAt(at uint32) (element T, err error) {
	if l.Length <= at {
		err = fmt.Errorf("Failed to get at %d, with error: at position is greater than list length of %d", at, l.Length)

	} else {
		element = l.Elements[at]
	}

	return element, nil
}

func (l *List[T]) Pop() (element T, err error) {
	if l.Length == 0 {
		err = fmt.Errorf("Failed to pop element, with error: empty list")

	} else {
		element = l.Elements[l.Length-1]
		l.Length--
	}

	return element, err
}

func (l *List[T]) PopFrom(from uint32) (element T, err error) {
	if l.Length <= from {
		err = fmt.Errorf("Failed to pop from %d, with error: from position is greater than list length of %d", from, l.Length)

	} else {
		element = l.Elements[from]
		for i := int(from); i < int(l.Length)-1; i++ {
			l.Elements[i] = l.Elements[i+1]
		}
		l.Length--
	}

	return element, err

}

func (l *List[T]) Clear() {
	for i := 0; i < int(l.Length); i++ {
		var element T
		l.Elements[i] = element
	}
	l.Length = 0
}

func (l *List[T]) Iterate(yield func(T) bool) {
	length := l.Length
	for i := 0; i < int(length); i++ {
		if !yield(l.Elements[i]) {
			return
		}
	}
}

func (l *List[T]) Items() []T {
	length := l.Length
	slice := make([]T, length)
	for i := 0; i < int(length); i++ {
		slice[i] = l.Elements[i]
	}
	return slice
}

func (l *List[T]) Empty() bool {
	return l.Length == 0
}
