
# Ordsync

Package ordsync allows to gather results (or do something else)
from concurrent goroutines in particular order.

## Example

```go
package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/josephbuchma/ordsync"
)

func main() {
	jobs := []time.Duration{4, 15, 3, 7, 1, 3, 23, 10, 5}
	results := []time.Duration{}

	dg := ordsync.DeferGroup{}
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

	dg.Wait()                                     // Wait until last deferred function is done
	fmt.Println(reflect.DeepEqual(jobs, results)) // true
}

```
