package collector

import (
	"sync"
	"time"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/dnsbl_exporter/pkg/rbl"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slog"
)

const namespace = "luzilla"
const subsystem = "rbls"

// RblCollector object as a bridge to prometheus
type RblCollector struct {
	configuredMetric  *prometheus.Desc
	blacklistedMetric *prometheus.Desc
	errorsMetrics     *prometheus.Desc
	listedMetric      *prometheus.Desc
	targetsMetric     *prometheus.Desc
	durationMetric    *prometheus.Desc
	rbls              []string
	util              *dns.DNSUtil
	targets           []string
	logger            *slog.Logger
}

func BuildFQName(metric string) string {
	return prometheus.BuildFQName(namespace, subsystem, metric)
}

// NewRblCollector ... creates the collector
func NewRblCollector(rbls []string, targets []string, util *dns.DNSUtil, logger *slog.Logger) *RblCollector {
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
		targetsMetric: prometheus.NewDesc(
			BuildFQName("targets"),
			"The number of targets that are being probed (configured via targets.ini or ?target=)",
			nil,
			nil,
		),
		durationMetric: prometheus.NewDesc(
			BuildFQName("duration"),
			"The scrape's duration (in seconds)",
			nil,
			nil,
		),
		rbls:    rbls,
		util:    util,
		targets: targets,
		logger:  logger,
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

	ch <- prometheus.MustNewConstMetric(
		c.targetsMetric,
		prometheus.GaugeValue,
		float64(len(hosts)),
	)

	start := time.Now()

	// this should be a map of blacklist and a counter (for listings)
	var listed sync.Map

	resolver := rbl.NewRBLResolver(c.logger, c.util)

	// iterate over hosts -> resolve to ip
	targets := make(chan []rbl.Target)
	for _, host := range hosts {
		go resolver.Do(host, targets)
	}

	// run the check
	for _, target := range <-targets {

		results := make([]rbl.Result, 0)

		result := make(chan rbl.Result)
		for _, blocklist := range c.rbls {
			logger := c.logger.With("host", target.Host)

			logger.Debug("starting check")

			r := rbl.New(c.util, logger)
			go r.Update(target, blocklist, result)
			results = append(results, <-result)
		}

		for _, check := range results {
			metricValue := 0

			val, _ := listed.LoadOrStore(check.Rbl, 0)
			if check.Listed {
				metricValue = 1
				listed.Store(check.Rbl, val.(int)+1)
			}

			c.logger.Debug("listed?", slog.Int("v", metricValue), slog.String("rbl", check.Rbl))

			labelValues := []string{check.Rbl, check.Target.IP.String(), check.Target.Host}

			// this is an "error" from the RBL/transport
			if check.Error {
				c.logger.Error(check.ErrorType.Error(), slog.String("text", check.Text))
				ch <- prometheus.MustNewConstMetric(
					c.errorsMetrics,
					prometheus.GaugeValue,
					1,
					[]string{check.Rbl}...,
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

	c.logger.Debug("building listed metric")

	for _, rbl := range c.rbls {
		val, _ := listed.LoadOrStore(rbl, 0)
		ch <- prometheus.MustNewConstMetric(
			c.listedMetric,
			prometheus.GaugeValue,
			float64(val.(int)),
			[]string{rbl}...,
		)
	}

	c.logger.Debug("finished")

	ch <- prometheus.MustNewConstMetric(
		c.durationMetric,
		prometheus.GaugeValue,
		time.Since(start).Seconds(),
	)

}
