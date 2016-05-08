package main

import (
	"fmt"
	"sync"
)

var (
	value  int
	change chan int
)

func init() {
	change = make(chan int, 100)
	go writer()
}

func writer() {
	for amount := range change {
		value += amount
	}
}

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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		loop(200000000)
		wg.Done()
	}()
	loop(-100000000)
	wg.Wait()
	fmt.Printf("Resulting value: %d\n", value)
}
