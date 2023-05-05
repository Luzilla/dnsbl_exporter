package prober

import (
	"net/http"

	"github.com/Luzilla/dnsbl_exporter/internal/setup"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ProberHandler struct {
	Resolver string
	Rbls     []string
}

func (p ProberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("target") {
		http.Error(w, "missing ?target parameter", http.StatusBadRequest)
		return
	}

	targets := make([]string, 0, 1)
	targets = append(targets, r.URL.Query().Get("target"))

	registry := setup.CreateRegistry()
	collector := setup.CreateCollector(p.Rbls, targets, p.Resolver)
	registry.MustRegister(collector)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	})
	h.ServeHTTP(w, r)
}
