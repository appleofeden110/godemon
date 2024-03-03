package main

type (
	Qnode[T any] struct {
		Value T
		Next  *Qnode[T]
		Prev  *Qnode[T]
	}
	Queue[T any] struct {
		Head   *Qnode[T]
		Tail   *Qnode[T]
		Length int8
	}
	Qinterface[T any] interface {
		enqueue(v T)
		dequeue() *T
	}
)

func newQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) enqueue(v T) {
	newNode := Qnode[T]{Value: v}
	if q.Head == nil {
		q.Head = &newNode
		q.Tail = &newNode
	} else {
		q.Tail = &newNode
		q.Tail.Next = &newNode
	}
	q.Length++
}

func (q *Queue[T]) dequeue() *T {
	if q.Head == nil {
		return nil
	}

	head := q.Head
	q.Head = q.Head.Next

	head.Next = nil
	q.Length--
	return &head.Value
}
