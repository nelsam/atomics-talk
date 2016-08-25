package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

type lock uint32

func (l *lock) Lock() {
	for {
		if atomic.CompareAndSwapUint32((*uint32)(l), 0, 1) {
			return
		}
		runtime.Gosched()
	}
}

func (l *lock) Unlock() {
	atomic.StoreUint32((*uint32)(l), 0)
}

var (
	l     lock
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
		l.Lock()
		value += amount
		l.Unlock()
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
