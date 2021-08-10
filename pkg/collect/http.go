package collect

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/ianunruh/shelly-exporter/pkg/shelly"
)

func ProbeHandler(log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get("target")
		if target == "" {
			http.Error(w, "target is required", http.StatusBadRequest)
			return
		}

		if !strings.Contains(target, "://") {
			target = "http://" + target
		}

		_, err := url.Parse(target)
		if err != nil {
			http.Error(w, "target is invalid URL", http.StatusBadRequest)
			return
		}

		c := NewCollector(newClient(target), log)

		registry := prometheus.NewRegistry()
		registry.MustRegister(c)

		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}
}

func newClient(target string) shelly.Client {
	username := os.Getenv("SHELLY_USERNAME")
	password := os.Getenv("SHELLY_PASSWORD")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	return shelly.NewClient(target, username, password, httpClient)
}
