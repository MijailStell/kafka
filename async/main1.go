package main

import (
	"fmt"
	"time"
)

func sum(a int, b int, c chan int) {
	fmt.Println("2")
	time.Sleep(10 * time.Second)
	c <- a + b // send sum to c
	fmt.Println("3")
}

func main() {
	c := make(chan int)
	go sum(20, 22, c)
	fmt.Println("1")
	x := <-c // receive from c
	fmt.Println("4")

	fmt.Println(x)
}
