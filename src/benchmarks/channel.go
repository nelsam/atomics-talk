package main

import (
	"fmt"
	"sync"
)

var (
	value  int
	change chan int
)

func loop(count int) {
	fmt.Printf("Started looping: %d\n", count)
	amount := 1
	if count < 0 {
		count = -count
		amount = -1
	}

	for i := 0; i < count; i++ {
		change <- amount
	}
}

func main() {
	var wg, writerWg sync.WaitGroup
	change = make(chan int, 100)
	wg.Add(1)
	writerWg.Add(1)

	go func() {
		for amount := range change {
			value += amount
		}
		writerWg.Done()
	}()
	go func() {
		loop(200000000)
		wg.Done()
	}()
	loop(-100000000)

	wg.Wait()
	close(change)
	writerWg.Wait()

	fmt.Printf("Resulting value: %d\n", value)
}
