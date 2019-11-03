package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const namespace = "luzilla"
const subsystem = "rbls"

// RblCollector object as a bridge to prometheus
type RblCollector struct {
	configuredMetric  *prometheus.Desc
	blacklistedMetric *prometheus.Desc
	errorsMetrics     *prometheus.Desc
	listedMetric      *prometheus.Desc
	rbls              []string
	resolver          string
	targets           []string
}

func buildFQName(metric string) string {
	return prometheus.BuildFQName(namespace, subsystem, metric)
}

// NewRblCollector ... creates the collector
func NewRblCollector(rbls []string, targets []string, resolver string) *RblCollector {
	return &RblCollector{
		configuredMetric: prometheus.NewDesc(
			buildFQName("used"),
			"The number of RBLs to check IPs against (configured via rbls.ini)",
			nil,
			nil,
		),
		blacklistedMetric: prometheus.NewDesc(
			buildFQName("ips_blacklisted"),
			"Blacklisted IPs",
			[]string{"rbl", "ip", "hostname"},
			nil,
		),
		errorsMetrics: prometheus.NewDesc(
			buildFQName("errors"),
			"The number of errors which occurred testing the RBLs",
			[]string{"rbl"},
			nil,
		),
		listedMetric: prometheus.NewDesc(
			buildFQName("listed"),
			"The number of listings in RBLs (this is bad)",
			[]string{"rbl"},
			nil,
		),
		rbls:     rbls,
		resolver: resolver,
		targets:  targets,
	}
}

// Describe ...
func (c *RblCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.configuredMetric
	ch <- c.blacklistedMetric
}

// Collect ...
func (c *RblCollector) Collect(ch chan<- prometheus.Metric) {
	// these are our targets to check
	hosts := c.targets

	ch <- prometheus.MustNewConstMetric(
		c.configuredMetric,
		prometheus.GaugeValue,
		float64(len(c.rbls)),
	)

	// this should be a map of blacklist and a counter
	listed := 0

	for _, host := range hosts {

		log.Debugln("Checking ...", host)

		rbl := NewRbl(c.resolver)
		rbl.Update(host, c.rbls)

		for _, result := range rbl.Results {
			// this is an "error" from the RBL
			if result.Error {
				log.Errorln(result.Text)
			}

			metricValue := 0

			if result.Listed {
				metricValue = 1
				listed = +1
			}

			labelValues := []string{result.Rbl, result.Address, host}

			if result.Error {
				ch <- prometheus.MustNewConstMetric(
					c.errorsMetrics,
					prometheus.GaugeValue,
					1,
					[]string{result.Rbl}...,
				)
			}

			ch <- prometheus.MustNewConstMetric(
				c.blacklistedMetric,
				prometheus.GaugeValue,
				float64(metricValue),
				labelValues...,
			)
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.listedMetric,
		prometheus.GaugeValue,
		float64(listed),
		[]string{"foo"}...,
	)
}
