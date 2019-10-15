package main

import (
	"log"
	"os"
	"time"

	"github.com/radhus/vsxknob/prometheus"
	"github.com/radhus/vsxknob/vsx"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: vsx-addr:8102")
	}
	addr := os.Args[1]

	prom := prometheus.New(addr)
	go prom.Start(":8080")

	connection, err := vsx.Connect(addr, prom)
	if err != nil {
		log.Fatalln("Couldn't connect to VSX:", err)
	}

	for {
		connection.CheckVolume()
		time.Sleep(250 * time.Millisecond)
		connection.CheckPower()
		time.Sleep(1 * time.Second)
	}
}
