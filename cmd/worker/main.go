package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		do()
	}

}

func do() {
	fmt.Println("hello world!")
}
