package main

import (
	"fmt"
	"os"
	"strconv"
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

	goroutines, err := strconv.Atoi(os.Getenv("GOROUTINES"))
	if err != nil {
		panic(err)
	}

	adders := goroutines - 1
	wg.Add(adders)
	delta := 200000000 / adders
	remaining := 200000000 % adders
	for i := 0; i < adders; i++ {
		go func() {
			loop(delta)
			wg.Done()
		}()
	}
	loop(-100000000 + remaining)
	wg.Wait()
	fmt.Printf("Resulting value: %d\n", value)
}
