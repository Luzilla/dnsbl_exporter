package setup

import (
	"github.com/Luzilla/dnsbl_exporter/collector"
	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

func CreateCollector(rblsIP []string, rblsDomain []string, targetsIP []string, targetsDomain []string, dnsUtil *dns.DNSUtil, logger *slog.Logger) *collector.RblCollector {
	return collector.NewRblCollector(rblsIP, rblsDomain, targetsIP, targetsDomain, dnsUtil, logger)
}

func CreateRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}
