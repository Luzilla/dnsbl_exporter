package setup

import (
	"github.com/Luzilla/dnsbl_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

func CreateCollector(rbls []string, targets []string, resolver string, logger *slog.Logger) *collector.RblCollector {
	return collector.NewRblCollector(rbls, targets, resolver, logger)
}

func CreateRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}
