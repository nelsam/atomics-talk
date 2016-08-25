package main

import (
	"math"
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

type channel struct {
	l        lock
	elements []float64
	closed   bool
	rHead    int
	wHead    int
}

func New(cap int) *channel {
	return &channel{
		elements: make([]float64, 0, cap),
	}
}

func (c *channel) Send(v float64) {
	for {
		c.l.Lock()
		if c.closed {
			defer c.l.Unlock()
			panic("can not send on a closed channel")
		}
		if len(c.elements) == cap(c.elements) {
			c.l.Unlock()
			continue
		}
		// add element
		c.elements = append(c.elements, v)
		c.l.Unlock()
		return
	}
}

func (c *channel) Receive() float64 {
	for {
		c.l.Lock()
		if c.closed && len(c.elements) == 0 {
			defer c.l.Unlock()
			return 0
		}
		if len(c.elements) == 0 {
			c.l.Unlock()
			continue
		}
		// remove element
		v := c.elements[0]
		copy(c.elements, c.elements[1:])
		c.elements = c.elements[:len(c.elements)-1]
		c.l.Unlock()
		return v
	}
}

func (c *channel) Close() {
	c.l.Lock()
	defer c.l.Unlock()
	c.closed = true
}

const depth = 10000000

var (
	in    = New(64)
	out   = New(64)
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
		in.Send(float64(i))
	}
	in.Close()
}

func consumer() {
	defer wg.Done()
	for {
		x := in.Receive()
		if x == 0 {
			return
		}
		y := math.Sqrt(math.Sqrt(math.Sqrt(math.Sqrt(math.Sqrt(x)))))
		out.Send(y)
	}
}

func reporter() {
	var total float64
	defer func() {
		println("total:", total)
		close(done)
	}()
	for {
		x := out.Receive()
		if x == 0 {
			return
		}
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
	out.Close()
	<-done
	println("done")
}
