package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	volumeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "volume",
			Help: "current volume level",
		},
		[]string{"addr"},
	)

	powerGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "up",
			Help: "power level",
		},
		[]string{"addr"},
	)

	addrLabel string
)

func reportVolume(volume int) {
	volumeGauge.With(prometheus.Labels{"addr": addrLabel}).Set(float64(volume))
}

func reportPower(on bool) {
	var value float64 = 1
	if !on {
		value = 0
		reportVolume(0)
	}
	powerGauge.With(prometheus.Labels{"addr": addrLabel}).Set(value)
}

func webserver(addr string) {
	addrLabel = addr
	prometheus.MustRegister(volumeGauge, powerGauge)

	http.Handle("/metrics", prometheus.Handler())
	err := http.ListenAndServe(":8080", nil)
	log.Fatalln("ListenAndServe failed:", err)
}
