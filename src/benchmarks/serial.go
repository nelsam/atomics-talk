package main

import "fmt"

var value int

func loop(count int) {
	fmt.Printf("Started looping: %d\n", count)
	amount := 1
	if count < 0 {
		count = -count
		amount = -1
	}

	for i := 0; i < count; i++ {
		value += amount
	}
}

func main() {
	loop(200000000)
	loop(-100000000)
	fmt.Printf("Resulting value: %d\n", value)
}
