package setup

import (
	"github.com/Luzilla/dnsbl_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

func CreateCollector(rbls []string, targets []string, resolver string) *collector.RblCollector {
	return collector.NewRblCollector(rbls, targets, resolver)
}

func CreateRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}
