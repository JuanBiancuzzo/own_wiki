package query

/*
Deberiamos tener:
 * Where
 * And/Or/Not
 * Join (oculto)
 * Union
 * GroupBy
Estaria bueno tener:
 * OrderBy
 * Distinc
 * Range (Limit)
 * Like/In/Between/Exist/Any/All
*/

/*
Por ahora es la query inyectada a las tablas
Ejemplo:
SELECT * FROM BookEntity WHERE BookEntity.Book.Author == "Jose" OR BookEntity.Book.Title LIKE "The world as %"
*/
type QueryRequest any

type Iterator[T any] struct{}

func NewIterator[T any](elements []T) Iterator[T] {
	return Iterator[T]{}
}

func (r Iterator[T]) Request(amount int) []T {
	return []T{}
}

type Limit[T any] struct {
	request []T
}

func NewLimit[T any](elements []T, amount int) Limit[T] {
	iterator := NewIterator(elements)
	return Limit[T]{
		request: iterator.Request(amount),
	}
}

func (l Limit[T]) Get() []T {
	return l.request
}
