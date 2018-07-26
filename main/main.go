////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2018 A Bit of Help, Inc. - All Rights Reserved, Worldwide.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Package main is the entry point for the application
// and is responsible for configuring the environment.
package main

import (
	"fmt"
	"github.com/abitofhelp/producerconsumer/consumer"
	"github.com/abitofhelp/producerconsumer/producer"
	"sync"
)

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
		consumer.Consumer(ch)
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
		producer.Producer(ch, []string{"/tmp", "/home/aboh/Downloads", "/home/aboh/go"})
	}()

	// Wait here until the producer and consumer have completed their work,
	// which will be signaled the channel being closed and by wg's internal
	// goroutine counter being zero.
	wg.Wait()

	// Adios!
	fmt.Println("All done!")
}
