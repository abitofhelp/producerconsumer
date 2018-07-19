////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package producer implements a producer of file system directory paths.
package producer

import (
	"fmt"
	"sync"
)

// Function Producer sends file system directory paths to the channel.
// Parameter ch is a unidirectional channel where directory paths will be sent.
// Parameter wg is a WaitGroup from the caller, and it is only used
// to decrement the goroutine counter in wg when production has completed.
func Producer(ch chan<- string, wg *sync.WaitGroup) {
	// Decrement the goroutine counter just before returning to the caller.
	// The caller used wg.Add(1) before its "go producer(ch, &wg) invocation,
	// so wg's goroutine counter was increased by one.  wg.Done() will decrease
	// the counter by one, so they cancel each other out.
	defer wg.Done()

	// Variable pwg is the producer's WaitGroup that is used to be able
	// to detect when all of the goroutines that are launched by the producer
	// have completed.
	var pwg sync.WaitGroup

	// Load the source directories into the channel.  In reality, this would be
	// some kind of a recursive file system walk.
	dirs := []string{"/tmp", "/home/aboh/Downloads", "/home/aboh/go"}

	for _, dir := range dirs {
		// Increment the producer's WaitGroup counter for each goroutine that is launched.
		pwg.Add(1)
		go func(d string) {
			// Decrement the producer's WaitGroup counter after each goroutine completes its work.
			// In this example, each goroutine will send a file system directory path to the channel.
			defer pwg.Done()

			// The producer's work is pretty simple...  Write the directory path to stdout and
			// to the channel.
			fmt.Println("Loading: " + d)
			ch <- d
		}(dir)
	}

	// The producer's wait group will block here until pwg's internal counter is zero.
	// In other words, after each goroutine that the producer has launched has completed.
	pwg.Wait()

	// The consumer is running in parallel to the producer.  The caller's code will wait for
	// all of the goroutines to complete at its wg.Wait() invocation.  It is important to
	// know that if the producer does not close the channel, then the consumer will block
	// and wg.Wait() will never be reached in Main().  Closing the channel signals the consumer
	// that production has completed.  When the consumer empties the channel, wg.Wait() in
	// Main() will complete.  So, now that all of the producer's goroutines have completed
	// and production has completed, we will close the channel.
	close(ch)

	return
}