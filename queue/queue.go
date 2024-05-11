package queue

import "fmt"

type QNode[T any] struct {
	Value T
	Next  *QNode[T]
}

type Queue[T any] struct {
	Head   *QNode[T]
	Tail   *QNode[T]
	Length int
}

// nil -> 1 -> 23 -> 56 -> nil
func (q *Queue[T]) Enqueue(val T) {
	n := new(QNode[T])
	n.Value = val
	if q.Length == 0 {
		q.Head = n
		q.Tail = n
	} else {
		q.Tail.Next = n
		q.Tail = n
	}
	q.Length++
}

func (q *Queue[T]) Deque() T {
	if q.Head == nil {
		return q.Head.Value
	}
	q.Length--
	head := q.Head
	q.Head = q.Head.Next

	head.Next = nil
	return head.Value
}

func (q *Queue[T]) Peek() T {
	if q.Head == nil {
		fmt.Println("EOF: empty")
		return q.Head.Value
	}
	return q.Head.Value
}
