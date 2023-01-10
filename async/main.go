package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sum(a int, b int, c chan int) {
	fmt.Println("2")
	time.Sleep(10 * time.Second)
	c <- a + b // send sum to c
	fmt.Println("3")
}

func main2() {
	c := make(chan int)
	go sum(20, 22, c)
	fmt.Println("1")
	x := <-c // receive from c
	fmt.Println("4")

	fmt.Println(x)
}

func main() {
	r := <-longRunningTask()
	fmt.Println(r)
}

func longRunningTask() <-chan int32 {
	r := make(chan int32)

	go func() {
		defer close(r)

		// simulate a workload
		time.Sleep(time.Second * 3)
		r <- rand.Int31n(100)
	}()

	return r
}
