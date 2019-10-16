package vsx

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/radhus/vsxknob/handler"
)

type Connection struct {
	conn     *net.TCPConn
	reporter handler.Reporter

	input chan string
}

func (c *Connection) reader(output chan<- string) {
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			log.Fatalln("Read failed:", err)
		}
		output <- string(buf[:n])
	}
}

func Connect(addr string, reporter handler.Reporter) (*Connection, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("ResolveTCPAddr failed: %w", err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("DialTCP failed: %w", err)
	}

	connection := &Connection{
		conn:     conn,
		reporter: reporter,

		input: make(chan string),
	}

	output := make(chan string)
	go connection.handler(output)
	go connection.reader(output)
	go connection.sender()

	connection.Poll()

	return connection, nil
}

func (c *Connection) sender() {
	for {
		message := <-c.input
		_, err := c.conn.Write([]byte(message + "\r\n"))
		if err != nil {
			log.Fatalln("Write failed:", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
