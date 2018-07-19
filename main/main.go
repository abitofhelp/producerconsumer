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
	. "github.com/abitofhelp/producerconsumer/consumer"
	. "github.com/abitofhelp/producerconsumer/producer"
)

// Function main is the entry point for the application and is responsible
// for configuring its environment.
func main() {

	// Variable wg is main's WaitGroup that is used to be able
	// to detect when all of the goroutines that are launched
	// have completed.
	var wg sync.WaitGroup

	// Variable ch is the channel into which the producer sends file
	// system directory paths, and from which the consumer will
	// retrieve the paths.
	ch := make(chan string)

	// Increment the producer's WaitGroup counter for the consumer goroutine.
	wg.Add(1)
	// Launch the consumer goroutine, which will block until the producer
	// sends file system directory paths into the channel.
	go Consumer(ch, &wg)

	// Increment the producer's WaitGroup counter for the consumer goroutine.
	wg.Add(1)
	// Launch the consumer goroutine, which will send file system directory
	// paths into the channel.
	go Producer(ch, &wg)

	// Wait here until the producer and consumer have completed their work,
	// which will be signaled by wg's internal goroutine counter being zero.
	wg.Wait()
	//close(ch)

	// Adios!
	fmt.Println("All done!")
}