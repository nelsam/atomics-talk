package main

import (
	"fmt"
	"sync"
)

var (
	mu    sync.Mutex
	value int
)

func loop(count int) {
	fmt.Printf("Started looping: %d\n", count)
	amount := 1
	if count < 0 {
		count = -count
		amount = -1
	}

	for i := 0; i < count; i++ {
		mu.Lock()
		value += amount
		mu.Unlock()
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
