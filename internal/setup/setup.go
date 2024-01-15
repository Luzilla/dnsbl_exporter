package setup

import (
	"github.com/Luzilla/dnsbl_exporter/collector"
	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

func CreateCollector(rbls []string, targets []string, domainBased bool, dnsUtil *dns.DNSUtil, logger *slog.Logger) *collector.RblCollector {
	return collector.NewRblCollector(rbls, targets, domainBased, dnsUtil, logger)
}

func CreateRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}
