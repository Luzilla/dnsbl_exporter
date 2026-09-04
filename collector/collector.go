package collector

import (
	"sync"
	"time"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/dnsbl_exporter/pkg/ip"
	"github.com/Luzilla/dnsbl_exporter/pkg/rbl"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
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
	rblsIP            []string
	rblsDomain        []string
	util              *dns.DNSUtil
	targetsIP         []string
	targetsDomain     []string
	logger            *slog.Logger
}

func BuildFQName(metric string) string {
	return prometheus.BuildFQName(namespace, subsystem, metric)
}

// NewRblCollector ... creates the collector. rblsIP is checked against
// targetsIP resolved to IPs, rblsDomain is checked against targetsDomain
// used as-is. Either pair may be empty to run in a single-mode configuration.
func NewRblCollector(rblsIP []string, rblsDomain []string, targetsIP []string, targetsDomain []string, util *dns.DNSUtil, logger *slog.Logger) *RblCollector {
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
			"Whether an error occurred while testing this target against the RBL (1) or not (0)",
			[]string{"rbl", "ip", "hostname"},
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
		rblsIP:        rblsIP,
		rblsDomain:    rblsDomain,
		util:          util,
		targetsIP:     targetsIP,
		targetsDomain: targetsDomain,
		logger:        logger,
	}
}

// Describe ...
func (c *RblCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.configuredMetric
	ch <- c.blacklistedMetric
	ch <- c.errorsMetrics
	ch <- c.listedMetric
	ch <- c.targetsMetric
	ch <- c.durationMetric
}

// Collect ...
func (c *RblCollector) Collect(ch chan<- prometheus.Metric) {
	// these are our targets to check
	hostsIP := ip.ExpandCIDRs(c.targetsIP)
	hostsDomain := ip.ExpandCIDRs(c.targetsDomain)

	ch <- prometheus.MustNewConstMetric(
		c.configuredMetric,
		prometheus.GaugeValue,
		float64(len(c.rblsIP)+len(c.rblsDomain)),
	)

	ch <- prometheus.MustNewConstMetric(
		c.targetsMetric,
		prometheus.GaugeValue,
		float64(len(hostsIP)+len(hostsDomain)),
	)

	start := time.Now()

	// this should be a map of blacklist and a counter (for listings)
	var listed sync.Map

	resolver := rbl.NewRBLResolver(c.logger, c.util)

	// domain based RBLs: targets are used as-is, no resolution needed
	if len(c.rblsDomain) > 0 {
		targets := make(chan rbl.Target)
		wg := sync.WaitGroup{}
		wg.Add(len(hostsDomain))
		go func() {
			wg.Wait()
			close(targets)
		}()
		for _, host := range hostsDomain {
			go func(hostname string) {
				targets <- rbl.Target{Host: hostname}
				wg.Done()
			}(host)
		}
		c.check(targets, c.rblsDomain, &listed, ch)
	}

	// IP based RBLs: targets are resolved to IPs first
	if len(c.rblsIP) > 0 {
		targets := make(chan rbl.Target)
		wg := sync.WaitGroup{}
		wg.Add(len(hostsIP))
		go func() {
			wg.Wait()
			close(targets)
		}()
		for _, host := range hostsIP {
			go resolver.Do(host, targets, wg.Done)
		}
		c.check(targets, c.rblsIP, &listed, ch)
	}

	c.logger.Debug("building listed metric")

	for _, rbl := range append(append([]string{}, c.rblsIP...), c.rblsDomain...) {
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

// check runs the given RBLs against every target on the channel, emitting
// the errors/blacklisted metrics and tallying listings into listed.
func (c *RblCollector) check(targets <-chan rbl.Target, rbls []string, listed *sync.Map, ch chan<- prometheus.Metric) {
	for target := range targets {

		results := make([]rbl.Result, 0)

		result := make(chan rbl.Result)
		for _, blocklist := range rbls {
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

			c.logger.Debug("listed?", slog.Int("v", metricValue), slog.String("rbl", check.Rbl), slog.String("reason", check.Text))
			ip := ""
			if len(check.Target.IP) > 0 {
				ip = check.Target.IP.String()
			}
			labelValues := []string{check.Rbl, ip, check.Target.Host}

			errorValue := 0.0
			if check.Error {
				c.logger.Error(check.ErrorType.Error(), slog.String("text", check.Text))
				errorValue = 1
			}
			ch <- prometheus.MustNewConstMetric(
				c.errorsMetrics,
				prometheus.GaugeValue,
				errorValue,
				labelValues...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.blacklistedMetric,
				prometheus.GaugeValue,
				float64(metricValue),
				labelValues...,
			)
		}
	}
}
