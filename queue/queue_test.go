package queue

import "testing"

func TestEnqueue(t *testing.T) {
	q := new(Queue[int])

	q.Enqueue(21345678432)
	q.Enqueue(1)
	q.Enqueue(4)
	q.Enqueue(26)
	q.Enqueue(86)
	trueVal := q.Peek()
	expVal := 21345678432
	if trueVal != expVal {
		t.Fatalf("got wrong peek value. got: %v, expected: %v", trueVal, expVal)
	}
}
func TestDeque(t *testing.T) {
	q := new(Queue[int])

	q.Enqueue(21345678432)
	q.Enqueue(1)
	q.Enqueue(4)
	q.Enqueue(26)
	q.Enqueue(86)

	trueVal := q.Deque()
	expVal := 21345678432

	trueVal2 := q.Peek()
	expVal2 := 1

	if trueVal != expVal {
		t.Fatalf("got wrong deque value. Got: %v, expected: %v", trueVal, expVal)
	}

	if trueVal2 != expVal2 {
		t.Fatalf("got wrong peek value. got: %v, expected: %v", trueVal2, expVal2)
	}
}
