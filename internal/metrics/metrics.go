package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	Registry         *prometheus.Registry
	RegistryExporter *prometheus.Registry
}

func (m MetricsHandler) Handler() http.Handler {
	return promhttp.HandlerFor(
		prometheus.Gatherers{
			m.Registry,
			m.RegistryExporter,
		},
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
			Registry:      m.Registry,
		},
	)
}
