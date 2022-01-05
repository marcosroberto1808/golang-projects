package main

import (
	"fmt"
	"time"
)

func test(){
	fmt.Println("Testing")
}

func main() {
	condition := false
	for {
		test()
		time.Sleep(time.Second * 5)
		if condition {
			break
		}
	}
}