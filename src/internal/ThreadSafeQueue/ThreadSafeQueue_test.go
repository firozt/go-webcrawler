package ThreadSafeQueue

import (
	"fmt"
	"sync"
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
	initial := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	q := NewThreadSafeQueueFromList[int](initial)
	NUM_DEQUES := len(initial) * 2 // last half should be failing
	var wg sync.WaitGroup
	wg.Add(NUM_DEQUES)

	falseCount, trueCount := 0, 0

	for i := 0; i < NUM_DEQUES; i++ {
		go func(i int) {
			defer wg.Done()
			_, isDequeud := q.Dequeue()

			if isDequeud {
				trueCount++
			} else {
				falseCount++
			}
		}(i)
	}

	wg.Wait()

	if trueCount != falseCount {
		t.Errorf("wanted %v got %v", trueCount, falseCount)
	}

}

// TestWasSeen() {}
