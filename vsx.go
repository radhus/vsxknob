package main

import (
	"log"
)

func handler(output <-chan string) {
	log.Println("Handler running...")
	for {
		line := <-output
		if len(line) < 3 {
			continue
		}
		message := line[:len(line)-2]

		if len(message) < 3 {
			log.Println("Skipping too short:", message)
		}

		switch message[:3] {
		case "VOL":
			log.Println("Got volume:", message)
		case "PWR":
			log.Println("Got power:", message)
		default:
			log.Println("Got unexpected:", message)
		}
	}
}

func volume(input chan<- string) {
	input <- "?V"
}

func power(input chan<- string) {
	input <- "?P"
}
