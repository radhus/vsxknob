package main

import (
	"log"
	"os"
	"time"

	"github.com/radhus/vsxknob/handler"
	"github.com/radhus/vsxknob/mqtt"
	"github.com/radhus/vsxknob/prometheus"
	"github.com/radhus/vsxknob/vsx"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Usage: vsx-addr:8102 mqtt-addr:1883")
	}
	addr := os.Args[1]
	mqttAddr := os.Args[2]

	prom := prometheus.New(addr)
	go prom.Start(":8080")

	mqtt, err := mqtt.New(mqttAddr)
	if err != nil {
		log.Fatalln("MQTT failed: ", err)
	}

	reporterMux := handler.Multiplex(prom, mqtt)
	connection, err := vsx.Connect(addr, reporterMux)
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
