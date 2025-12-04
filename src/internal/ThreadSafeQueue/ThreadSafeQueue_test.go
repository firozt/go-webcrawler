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

// func TestDequeue() {}

// TestWasSeen() {}
