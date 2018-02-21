package main

import (
	"log"
	"strconv"
)

func handleVolume(message string) {
	volume, err := strconv.Atoi(message[3:])
	if err != nil {
		log.Println("Atoi failed for message:", message, "err:", err)
		return
	}

	log.Println("Got volume:", volume)
}

func handlePower(message string) {
	switch message[3] {
	case '0':
		log.Println("Power is on")
	case '2':
		log.Println("Power is standby")
	default:
		log.Println("Unexpected power:", message)
	}
}

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
			handleVolume(message)
		case "PWR":
			handlePower(message)
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
