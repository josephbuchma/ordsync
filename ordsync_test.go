package ordsync

import (
	"reflect"
	"testing"
	"time"
)

func TestDeferGroup(t *testing.T) {
	for i := 0; i < 20; i++ {
		jobs := []time.Duration{4, 15, 3, 7, 1, 3, 23, 10, 5}
		results := []time.Duration{}

		dg := DeferGroup{}
		for _, j := range jobs {
			j := j
			deferred := dg.Defer()
			go func() {
				time.Sleep(j * time.Millisecond)
				deferred.Do(func() {
					results = append(results, j) // no data race here
				})
			}()
		}

		dg.Wait() // Wait until last deferred function is done
		if !reflect.DeepEqual(jobs, results) {
			t.Errorf("Invalid results, expected %#v, got %#v", jobs, results)
		}
	}

	// Test panic on double Deferred.Do call
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()
	dfr := DeferGroup{}
	deferred := dfr.Defer()
	deferred.Do(func() {})
	deferred.Do(func() {})
	dfr.Wait()
}
