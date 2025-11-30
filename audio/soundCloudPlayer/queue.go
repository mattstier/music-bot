package soundCloudPlayer

type Node[T any] struct {
	value T
	next  *Node[T]
}

type Queue[T any] struct {
	head *Node[T] // points to the first element
	tail *Node[T] // points to the last element
	size int
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Enqueue(value T) {
	node := &Node[T]{value: value}
	if q.tail == nil {
		// Queue is empty
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.size++
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if q.head == nil {
		var zero T
		return zero, false
	}
	node := q.head
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return node.value, true
}

func (q *Queue[T]) Peek() (T, bool) {
	if q.head == nil {
		var zero T
		return zero, false
	}
	return q.head.value, true
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}
