package main

import (
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: vsx-addr:8102")
	}
	addr := os.Args[1]

	input := make(chan string)
	go connection(addr, input)

	for {
		volume(input)
		time.Sleep(250 * time.Millisecond)
		power(input)
		time.Sleep(1 * time.Second)
	}
}
