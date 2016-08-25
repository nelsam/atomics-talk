package main

import (
	"math"
	"os"
	"strconv"
	"sync"
)

const depth = 10000000

var (
	in    = make(chan float64, 64)
	out   = make(chan float64, 64)
	done  = make(chan struct{})
	wg    sync.WaitGroup
	width int
)

func init() {
	var err error
	width, err = strconv.Atoi(os.Getenv("GOROUTINES"))
	if err != nil {
		panic(err)
	}

	wg.Add(width)
}

func producer() {
	for i := 1; i <= depth; i++ {
		in <- float64(i)
	}
	close(in)
}

func consumer() {
	defer wg.Done()
	for x := range in {
		out <- math.Sqrt(math.Sqrt(math.Sqrt(math.Sqrt(math.Sqrt(x)))))
	}
}

func reporter() {
	var total float64
	defer func() {
		println("total:", total)
		close(done)
	}()
	for x := range out {
		total += x
	}
}

func main() {
	go producer()
	for i := 0; i < width; i++ {
		go consumer()
	}
	go reporter()
	wg.Wait()
	close(out)
	<-done
}
