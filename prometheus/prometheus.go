package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	volumeGauge *prometheus.GaugeVec
	powerGauge  *prometheus.GaugeVec
	labels      prometheus.Labels
}

func New(addressLabel string) *Server {
	server := &Server{
		volumeGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "volume",
				Help: "current volume level",
			},
			[]string{"addr"},
		),
		powerGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "up",
				Help: "power level",
			},
			[]string{"addr"},
		),
		labels: prometheus.Labels{"addr": addressLabel},
	}

	prometheus.MustRegister(server.volumeGauge, server.powerGauge)
	return server
}

func (s *Server) ReportVolume(volume int) {
	s.volumeGauge.With(s.labels).Set(float64(volume))
}

func (s *Server) ReportPower(on bool) {
	var value float64 = 1
	if !on {
		value = 0
	}
	s.powerGauge.With(s.labels).Set(value)
}

func (s *Server) Start(listenAddr string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listenAddr, nil)
	log.Fatalln("ListenAndServe failed:", err)
}
