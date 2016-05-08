package main

import (
	"fmt"
	"os"
	"strconv"
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

	goroutines, err := strconv.Atoi(os.Getenv("GOROUTINES"))
	if err != nil {
		panic(err)
	}

	writerWg.Add(1)
	go func() {
		for amount := range change {
			value += amount
		}
		writerWg.Done()
	}()

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
	close(change)
	writerWg.Wait()

	fmt.Printf("Resulting value: %d\n", value)
}
