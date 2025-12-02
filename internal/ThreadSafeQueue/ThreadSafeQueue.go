package ThreadSafeQueue

import (
	"fmt"
	"sync"
)

type ThreadSafeQueue[T comparable] struct {
	elements []T
	seen     map[T]struct{}
	lock     sync.Mutex
}

// creates a new queue
func NewThreadSafeQueue[T comparable]() *ThreadSafeQueue[T] {
	return &ThreadSafeQueue[T]{
		elements: make([]T, 0),
		seen:     make(map[T]struct{}),
	}
}

// adds element to back of the queue
func (q *ThreadSafeQueue[T]) Enqueue(elem T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.elements = append(q.elements, elem)
	q.seen[elem] = struct{}{}
}

// removes element from the front of the queue
func (q *ThreadSafeQueue[T]) Dequeue() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.elements) == 0 {
		var zero T
		return zero, false
	}

	elem := q.elements[0]
	q.elements = q.elements[1:]
	return elem, true
}

// checks if a value has been enqueued before
func (q *ThreadSafeQueue[T]) WasSeen(elem T) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	_, ok := q.seen[elem]
	return ok
}

// returns the current queue length
func (q *ThreadSafeQueue[T]) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.elements)
}

func (q *ThreadSafeQueue[t]) All() {
	fmt.Println(len(q.seen))
}
