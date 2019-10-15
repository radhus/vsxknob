package vsx

import (
	"log"
	"strconv"
)

func (c *Connection) handleVolume(message string) {
	volume, err := strconv.Atoi(message[3:])
	if err != nil {
		log.Println("Atoi failed for message:", message, "err:", err)
		return
	}
	c.reporter.ReportVolume(volume)
}

func (c *Connection) handlePower(message string) {
	switch message[3] {
	case '0':
		c.reporter.ReportPower(true)
	case '2':
		c.reporter.ReportPower(false)
	default:
		log.Println("Unexpected power:", message)
	}
}

func (c *Connection) handler(output <-chan string) {
	log.Println("Handler running...")
	for {
		line := <-output
		if len(line) < 3 {
			continue
		}
		message := line[:len(line)-2]

		if len(message) < 3 {
			log.Println("Skipping too short:", message)
			return
		}

		switch message[:3] {
		case "VOL":
			c.handleVolume(message)
		case "PWR":
			c.handlePower(message)
		default:
			log.Println("Got unexpected:", message)
		}
	}
}

func (c *Connection) CheckVolume() {
	c.input <- "?V"
}

func (c *Connection) CheckPower() {
	c.input <- "?P"
}
