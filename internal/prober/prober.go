package prober

import (
	"net/http"

	"github.com/Luzilla/dnsbl_exporter/internal/setup"
	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
)

type ProberHandler struct {
	DNS         *dns.DNSUtil
	Rbls        []string
	DomainBased bool
	Logger      *slog.Logger
}

func (p ProberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("target") {
		p.Logger.Error("missing ?target parameter")
		http.Error(w, "missing ?target parameter", http.StatusBadRequest)
		return
	}

	targets := make([]string, 0, 1)
	targets = append(targets, r.URL.Query().Get("target"))

	registry := setup.CreateRegistry()
	collector := setup.CreateCollector(p.Rbls, targets, p.DomainBased, p.DNS, p.Logger)
	registry.MustRegister(collector)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	})
	h.ServeHTTP(w, r)
}
