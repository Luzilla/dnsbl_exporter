package collector

import (
	"sync"
	"time"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/dnsbl_exporter/pkg/rbl"
	x "github.com/miekg/dns"
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
	durationMetric    *prometheus.Desc
	rbls              []string
	resolver          string
	targets           []string
}

func BuildFQName(metric string) string {
	return prometheus.BuildFQName(namespace, subsystem, metric)
}

// NewRblCollector ... creates the collector
func NewRblCollector(rbls []string, targets []string, resolver string) *RblCollector {
	return &RblCollector{
		configuredMetric: prometheus.NewDesc(
			BuildFQName("used"),
			"The number of RBLs to check IPs against (configured via rbls.ini)",
			nil,
			nil,
		),
		blacklistedMetric: prometheus.NewDesc(
			BuildFQName("ips_blacklisted"),
			"Blacklisted IPs",
			[]string{"rbl", "ip", "hostname"},
			nil,
		),
		errorsMetrics: prometheus.NewDesc(
			BuildFQName("errors"),
			"The number of errors which occurred testing the RBLs",
			[]string{"rbl"},
			nil,
		),
		listedMetric: prometheus.NewDesc(
			BuildFQName("listed"),
			"The number of listings in RBLs (this is bad)",
			[]string{"rbl"},
			nil,
		),
		durationMetric: prometheus.NewDesc(
			BuildFQName("duration"),
			"The scrape's duration (in seconds)",
			nil,
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

	start := time.Now()

	// this should be a map of blacklist and a counter (for listings)
	var listed sync.Map

	// iterate over hosts -> resolve to ip, check
	for _, host := range hosts {

		log.Debugln("Checking ...", host)

		r := rbl.New(dns.New(new(x.Client), c.resolver))
		r.Update(host, c.rbls)

		for _, result := range r.Results {
			// this is an "error" from the RBL
			if result.Error {
				log.Errorln(result.Text)
			}

			metricValue := 0

			val, _ := listed.LoadOrStore(result.Rbl, 0)
			if result.Listed {
				metricValue = 1
				listed.Store(result.Rbl, val.(int)+1)
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

	for _, rbl := range c.rbls {
		val, _ := listed.LoadOrStore(rbl, 0)
		ch <- prometheus.MustNewConstMetric(
			c.listedMetric,
			prometheus.GaugeValue,
			float64(val.(int)),
			[]string{rbl}...,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.durationMetric,
		prometheus.GaugeValue,
		time.Since(start).Seconds(),
	)

}
