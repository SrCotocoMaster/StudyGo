package main

import "fmt"

func main() {
	ch := make(chan string, 2)
	ch <- "hello"
	ch <- "mundo"

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	//fmt.Println(<-ch)
}
