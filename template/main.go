package main

import "fmt"

func main() {
	fmt.Println("1")
	producer(consumer01)
	fmt.Println("despues")
}

func producer(c func() string) {
	valueFound := c()
	fmt.Println(valueFound)
}

func consumer01() string {
	return "consumer 01"
}
