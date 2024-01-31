package prom

import (
	"github.com/nitwhiz/omnilock/pkg/server"
	"io"
	"net/http"
	"strings"
)

var currentLockCount = Metric{
	Name: "current_lock_count",
	Type: "gauge",
	Help: "The current lock count",
	Sampler: func(s *server.Server) *Sample {
		return NewSample(float64(s.GetCurrentLockCount()))
	},
}

func getMetrics(s *server.Server, w http.ResponseWriter, r *http.Request) {
	str := strings.Builder{}

	str.WriteString(currentLockCount.string(s))

	_, err := io.WriteString(w, str.String())

	if err != nil {
		return
	}
}

func Listen(s *server.Server) error {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		getMetrics(s, w, r)
	})

	return http.ListenAndServe("0.0.0.0:7195", nil)
}
