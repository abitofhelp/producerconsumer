////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package main is the entry point for the application
// and is responsible for configuring the environment.
package main

import (
	"fmt"
	"sync"
)

// Function Consumer receives file system directory paths from the channel.
// Parameter ch is a unidirectional channel from which directory paths are
// retrieved.
// It returns after the Producer has closed the channel to signal no more data, and
// the Consumer is no longer getting data from the channel.
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
			fmt.Println("Received:\t" + path)
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

// Function Producer sends file system directory paths to the channel for consumption.
// Parameter ch is a unidirectional channel where directory paths will be sent.
// Parameter dirs is an array of file system directory paths.
// It returns when all of the goroutines have sent their paths in the channel.
func Producer(ch chan<- string, dirs []string) {
	// Variable wg is the producer's WaitGroup, which detects when all of the
	// goroutines that were launched have completed.
	var wg sync.WaitGroup

	// Loop through the array of directory paths...
	for _, dir := range dirs {
		// Increment the producer's WaitGroup counter for each goroutine that is launched.
		wg.Add(1)
		go func(dir string) {
			// Decrement the producer's WaitGroup counter after each goroutine completes its work.
			defer wg.Done()

			// The producer's work is pretty simple...  Write the directory path to stdout and
			// to the channel.
			fmt.Println("Sending:\t" + dir)
			ch <- dir
		}(dir)
		// The closure is only bound to that one variable, 'dir'.  There is a very good
		// chance that not adding 'dir' as a parameter to the closure, will result in seeing
		// the last element printed for every iteration instead of each value in sequence.
		// This is due to the high probability that goroutines will execute after the loop.
		//
		// By adding 'dir' as a parameter to the closure, 'dir' is evaluated at each iteration
		// and placed on the stack for the goroutine, so each slice element is available
		// to the goroutine when it is eventually executed.
	}

	// The producer's wait group will block here until wg's internal counter is zero, which
	// happens when all of the goroutines that were launched have completed.
	wg.Wait()
}

// Function main is the entry point for the application and is responsible
// for configuring its environment.
func main() {
	// Variable wg is main's WaitGroup, which detects when all of the
	// goroutines that were launched have completed.
	var wg sync.WaitGroup

	// Variable ch is the channel into which the producer sends file
	// system directory paths, and from which the consumer will
	// retrieve the paths.
	ch := make(chan string, 2)

	// Increment the producer's WaitGroup counter for the goroutine
	// launching the Consumer().
	wg.Add(1)
	// Launch the goroutine, which will block until the producer
	// sends file system directory paths into the channel.
	go func() {
		// Decrement the consumer's WaitGroup counter just before the goroutine exits.
		defer wg.Done()
		Consumer(ch)
	}()

	// Increment the producer's WaitGroup counter for the goroutine
	// launching the Producer().
	wg.Add(1)
	// Launch the goroutine, which will send file system directory
	// paths into the channel.
	go func() {
		// Decrement the producer's WaitGroup counter just before the goroutine exits.
		defer wg.Done()

		// The consumer is running in parallel to the producer.  Function main() will wait for
		// the consumer and producer goroutines to complete at its wg.Wait().  It is important to
		// know that if the producer does not close the channel, then the consumer will block
		// and wg.Wait() will never be reached in Main().  Closing the channel signals the consumer
		// that production has completed.  When the consumer empties the channel, wg.Wait() in
		// Main() will complete.  The channel will be closed just before the goroutine exits.
		defer close(ch)

		// Load the source directories into the channel so they can be consumed.
		Producer(ch, []string{"/tmp", "/home/aboh/Downloads", "/home/aboh/go"})
	}()

	// Wait here until the producer and consumer have completed their work,
	// which will be signaled the channel being closed and by wg's internal
	// goroutine counter being zero.
	wg.Wait()

	// Adios!
	fmt.Println("All done!")
}
