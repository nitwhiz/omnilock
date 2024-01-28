package prometheus

import (
	"fmt"
	"github.com/nitwhiz/omnilock/pkg/server"
	"io"
	"net/http"
	"strings"
	"time"
)

func getMetrics(s *server.Server, w http.ResponseWriter, r *http.Request) {
	now := time.Now().UnixMilli()

	var m []string

	m = append(m, "# HELP omnilock_current_lock_count The current lock count")
	m = append(m, "# TYPE gauge")
	m = append(m, fmt.Sprintf("omnilock_current_lock_count %d %d", s.GetCurrentLockCount(), now))

	_, err := io.WriteString(w, strings.Join(m, "\n"))

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
