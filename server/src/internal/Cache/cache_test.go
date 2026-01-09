package cache_test

import (
	"fmt"
	"testing"

	cache "github.com/firozt/crawler/src/internal/Cache"
)

func TestSet(t *testing.T) {
	type TestCase struct {
		key   string
		value string
	}
	cachestore := cache.MakeCacheStore[string]()
	tc := []TestCase{
		TestCase{
			key:   "url1",
			value: "value",
		},
		TestCase{
			key:   "url2",
			value: "value2",
		},
		TestCase{
			key:   "url1",
			value: "new value 1",
		},
		TestCase{
			key:   "url3",
			value: "value 3",
		},
	}
	expectedSize := []int{1, 2, 2, 3}

	for idx, test := range tc {
		t.Run(fmt.Sprintf("test-%d", idx), func(t *testing.T) {
			cachestore.Set(test.key, test.value)
			if cachestore.Size() != uint64(expectedSize[idx]) {
				t.Fatal("Set method size incorrect")
			}
		})
	}
}
