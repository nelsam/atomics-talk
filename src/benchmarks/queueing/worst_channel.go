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
	l        lock
	elements []float64
	closed   uint32
}

func New(cap int) *Channel {
	return &Channel{
		elements: make([]float64, 0, cap),
	}
}

func (c *Channel) Send(v float64) {
	for !c.trySend(v) {
		runtime.Gosched()
	}
}

func (c *Channel) trySend(v float64) bool {
	if atomic.LoadUint32(&c.closed) == 1 {
		panic("can not send on a closed channel")
	}
	c.l.Lock()
	defer c.l.Unlock()
	if len(c.elements) == cap(c.elements) {
		return false
	}
	c.elements = append(c.elements, v)
	return true
}

func (c *Channel) Receive() float64 {
	v, found, closed := c.tryReceive()
	for !found && !closed {
		runtime.Gosched()
		v, found, closed = c.tryReceive()
	}
	return v
}

func (c *Channel) tryReceive() (v float64, found, closed bool) {
	c.l.Lock()
	defer c.l.Unlock()
	if len(c.elements) == 0 {
		return 0, false, atomic.LoadUint32(&c.closed) == 1
	}
	v = c.elements[0]
	copy(c.elements, c.elements[1:])
	c.elements = c.elements[:len(c.elements)-1]
	return v, true, false
}

func (c *Channel) Close() {
	atomic.StoreUint32(&c.closed, 1)
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
