package vsx

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func (c *Connection) handleVolume(message string) {
	volume, err := strconv.Atoi(message[3:])
	if err != nil {
		log.Println("Atoi failed for message:", message, "err:", err)
		return
	}
	volume = (volume - 1) / 2

	c.volume.last = volume
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

func (c *Connection) handleMuted(message string) {
	muted := message[3] == '0'
	c.reporter.ReportMuted(muted)
}

func (c *Connection) handleSource(message string) {
	sourceID := message[2:]
	source, found := sources[sourceID]
	if !found {
		log.Println("Unsupported source:", source)
		return
	}
	c.reporter.ReportSource(source)
}

func (c *Connection) handler(output <-chan string) {
	log.Println("Handler running...")
	for {
		line := <-output
		if len(line) < 3 {
			continue
		}
		message := line[:len(line)-2]

		if len(message) < 2 {
			log.Println("Skipping too short:", message)
			continue
		}
		go c.handleMessage(message)
	}
}

func (c *Connection) handleMessage(message string) {
	switch message[:2] {
	case "VO":
		c.handleVolume(message)
	case "PW":
		c.handlePower(message)
	case "MU":
		c.handleMuted(message)
	case "FN":
		c.handleSource(message)
	default:
		log.Println("Got unexpected:", message)
	}
}

func (c *Connection) checkVolume() {
	c.input <- "?V"
}

func (c *Connection) CheckVolume() {
	if c.volume.active {
		return
	}
	c.checkVolume()
}

func (c *Connection) CheckPower() {
	c.input <- "?P"
}

func (c *Connection) CheckMuted() {
	c.input <- "?M"
}

func (c *Connection) CheckSource() {
	c.input <- "?F"
}

func (c *Connection) Poll() {
	c.CheckPower()
	c.CheckVolume()
	c.CheckMuted()
	c.CheckSource()
}

func (c *Connection) SetPower(on bool) {
	cmd := "PO"
	if !on {
		cmd = "PF"
	}

	c.input <- cmd
	c.CheckPower()
}

func (c *Connection) SetVolume(volume int) {
	c.volume.wanted = volume

	c.volume.mutex.Lock()
	defer c.volume.mutex.Unlock()

	c.volume.active = true
	defer func() {
		c.volume.active = false
	}()

	if c.volume.wanted != volume {
		return
	}

	start := time.Now()
	last := -1
	for c.volume.wanted != c.volume.last {
		if start.Add(3 * time.Second).Before(time.Now()) {
			log.Println("Failing to reconciliate volume after 2 seconds, aborting")
			return
		}

		if c.volume.last != last {
			last = c.volume.last
			diff := c.volume.wanted - c.volume.last

			cmd := "VU"
			steps := diff
			if steps < 0 {
				cmd = "VD"
				steps = -steps
			}
			if steps > 5 {
				steps = 5
			}

			for i := 0; i < steps; i++ {
				c.input <- cmd
			}
		} else {
			time.Sleep(25 * time.Millisecond)
		}
		c.checkVolume()
	}
}

func (c *Connection) SetMute(muted bool) {
	cmd := "MO"
	if !muted {
		cmd = "MF"
	}

	c.input <- cmd
	c.CheckMuted()
}

func (c *Connection) SetSource(source string) {
	sourceID, found := sourceIDs[source]
	if !found {
		log.Println("Couldn't set to unknown source:", source)
		return
	}

	cmd := fmt.Sprintf("%sFN", sourceID)
	c.input <- cmd
	c.CheckSource()
}
