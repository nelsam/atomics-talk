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

type Channel struct {
	rcond    *sync.Cond
	wcond    *sync.Cond
	elements []float64
	closed   uint32
}

func New(cap int) *Channel {
	var r, w lock
	return &Channel{
		rcond:    sync.NewCond(&r),
		wcond:    sync.NewCond(&w),
		elements: make([]float64, 0, cap),
	}
}

func (c *Channel) Send(v float64) {
	if atomic.LoadUint32(&c.closed) == 1 {
		panic("can not send on a closed channel")
	}
	c.wcond.L.Lock()
	defer c.wcond.L.Unlock()
	c.rcond.L.Lock()
	defer c.rcond.L.Unlock()
	for len(c.elements) == cap(c.elements) {
		c.rcond.L.Unlock()
		c.wcond.Wait()
		c.rcond.L.Lock()
		if atomic.LoadUint32(&c.closed) == 1 {
			panic("can not send on a closed channel")
		}
	}
	c.elements = append(c.elements, v)
	c.rcond.Signal()
}

func (c *Channel) Receive() float64 {
	c.rcond.L.Lock()
	defer c.rcond.L.Unlock()
	for len(c.elements) == 0 {
		if atomic.LoadUint32(&c.closed) == 1 {
			return 0
		}
		c.rcond.Wait()
	}
	v := c.elements[0]
	copy(c.elements, c.elements[1:])
	c.elements = c.elements[:len(c.elements)-1]
	c.wcond.Signal()
	return v
}

func (c *Channel) Close() {
	atomic.StoreUint32(&c.closed, 1)
	c.rcond.Broadcast()
	c.wcond.Broadcast()
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
