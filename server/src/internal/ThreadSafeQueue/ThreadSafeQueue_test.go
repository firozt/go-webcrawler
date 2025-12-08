package ThreadSafeQueue

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestEnqueue(t *testing.T) {
	q := NewThreadSafeQueue[string]()
	NUM_ENQUEUES := 1000
	var wg sync.WaitGroup
	wg.Add(NUM_ENQUEUES)

	// start NUM_ENQUEUES number of goroutines
	for i := 0; i < NUM_ENQUEUES; i++ {
		go func(i int) {
			defer wg.Done()
			q.Enqueue(fmt.Sprint(i))
		}(i)
	}

	wg.Wait()

	// check all were added
	if q.Len() != NUM_ENQUEUES {
		t.Errorf("wanted %v got %v", NUM_ENQUEUES, q.Len())
	}

	// check every item was enqueued
	for i := 0; i < NUM_ENQUEUES; i++ {
		val := fmt.Sprint(i)
		if !q.WasSeen(val) {
			t.Errorf("value %v was not seen", val)
		}
	}
}

func TestDequeue(t *testing.T) {
	initial := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	q := NewThreadSafeQueueFromList[int](initial)

	NUM_DEQUES := len(initial) * 2 // 40 dequeues
	var wg sync.WaitGroup
	wg.Add(NUM_DEQUES)

	var successCount int64
	var failCount int64

	for i := 0; i < NUM_DEQUES; i++ {
		go func() {
			defer wg.Done()
			_, ok := q.Dequeue()
			if ok {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}
		}()
	}

	wg.Wait()

	if successCount != int64(len(initial)) {
		t.Errorf("wanted %v successful dequeues, got %v", len(initial), successCount)
	}
	if failCount != int64(len(initial)) {
		t.Errorf("wanted %v failed dequeues, got %v", len(initial), failCount)
	}
}
