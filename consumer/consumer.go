////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package consumer implements a consumer of file system directory paths.
package consumer

import (
	"fmt"
	"sync"
)

// Function Consumer receives file system directory paths from the channel.
// Parameter ch is a unidirectional channel from which directory paths are
// retrieved.
// It returns when all of the goroutines have completed processing the channel.
func Consumer(ch <-chan string) {
	// Variable wg is the consumer's WaitGroup, which detects when all of the
	// goroutines that were launched have completed.
	var wg sync.WaitGroup

	// The consumer's wait group will block at wg.Wait(), which will be invoked just
	// before exiting from the function.  It will block until wg's internal counter is zero,
	// which happens when all of the goroutines that were launched have completed.
	defer wg.Wait()

	for path := range ch {
		// Increment the consumer's WaitGroup counter for each goroutine that is launched.
		wg.Add(1)
		go func(path string) {
			// Decrement the consumers's WaitGroup counter after each goroutine completes its work.
			defer wg.Done()

			// The consumer's work is pretty simple...  Write the directory path that was retrieved
			// from the channel to stdout.
			fmt.Println("Received: " + path)
		}(path)
		// The closure is only bound to that one variable, 'path'.  There is a very good
		// chance that not adding 'path' as a parameter to the closure, will result in seeing
		// the last element printed for every iteration instead of each value in sequence.
		// This is due to the high probability that goroutines will execute after the loop.
		//
		// By adding 'path' as a parameter to the closure, 'path' is evaluated at each iteration
		// and placed on the stack for the goroutine, so each slice element is available
		// to the goroutine when it is eventually executed.
	}
}
