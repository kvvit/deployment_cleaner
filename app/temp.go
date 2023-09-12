package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second * 5)
	for ; true; <-ticker.C {
		fmt.Println("hi")
	}
}
