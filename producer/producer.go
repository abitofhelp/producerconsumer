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
