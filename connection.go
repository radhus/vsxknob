package main

import (
	"io"
	"log"
	"net"
)

func closer(conn io.Closer) {
	if err := conn.Close(); err != nil {
		log.Println("Close failed:", err)
	}
}

func reader(conn *net.TCPConn, output chan<- string) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalln("Read failed:", err)
		}
		output <- string(buf[:n])
	}
}

func connection(addr string, input <-chan string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalln("ResolveTCPAddr failed: err")
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		log.Fatalln("DialTCP failed:", err)
	}
	defer closer(conn)

	output := make(chan string)
	go handler(output)
	go reader(conn, output)
	for {
		message := <-input
		_, err = conn.Write([]byte(message + "\r\n"))
		if err != nil {
			log.Fatalln("Write failed:", err)
		}
	}
}
