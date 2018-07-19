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
// Parameter wg is a WaitGroup from the caller, and it is only used
// to decrement the goroutine counter in wg when consumption has completed.
func Consumer(ch <-chan string, wg *sync.WaitGroup) {
	// Decrement the goroutine counter just before returning to the caller.
	// The caller used wg.Add(1) before its "go consumer(ch, &wg) invocation,
	// so wg's goroutine counter was increased by one.  wg.Done() will decrease
	// the counter by one, so they cancel each other out.
	defer wg.Done()

	// Variable cwg is the consumer's WaitGroup that is used to be able
	// to detect when all of the goroutines that are launched by the consumer
	// have completed.
	var cwg sync.WaitGroup

	for path := range ch {
		// Increment the consumer's WaitGroup counter for each goroutine that is launched.
		cwg.Add(1)
		go func(p string) {
			// Decrement the consumers's WaitGroup counter after each goroutine completes its work.
			// In this example, each goroutine will send a file system directory path to the channel.
			defer cwg.Done()

			// The consumer's work is pretty simple...  Write the directory path that was retrieved
			// to stdout.
			fmt.Println("Received: " + p)

			return
		}(path)
	}

	// The consumer's wait group will block here until cwg's internal counter is zero.
	// In other words, after each goroutine that the consumer has launched has completed.
	cwg.Wait()

	return
}