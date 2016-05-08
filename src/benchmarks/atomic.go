package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var value int64

func loop(count int) {
	fmt.Printf("Started looping: %d\n", count)
	amount := int64(1)
	if count < 0 {
		count = -count
		amount = -1
	}

	for i := 0; i < count; i++ {
		atomic.AddInt64(&value, amount)
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
