package rbl

import (
	"fmt"
	"net"
	"strings"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/godnsbl"
	"golang.org/x/exp/slog"
)

// RblResult extends godnsbl and adds RBL name
type Result struct {
	Target    Target
	Listed    bool
	Text      string
	Error     bool
	ErrorType error
	Rbl       string
}

// Rbl ... object
type RBL struct {
	Results []Result
	util    *dns.DNSUtil
	logger  *slog.Logger
}

// NewRbl ... factory
func New(util *dns.DNSUtil, logger *slog.Logger) *RBL {
	return &RBL{
		logger:  logger,
		util:    util,
		Results: make([]Result, 0),
	}
}

// Update runs the checks for an IP against against all "rbls"
func (rbl *RBL) Update(target Target, blocklist string, c chan<- Result) {
	go rbl.lookup(blocklist, target, c, rbl.logger.With(
		slog.Group("unit",
			slog.String("target", target.Host),
			slog.String("rbl", blocklist))))
}

func (r *RBL) lookup(blocklist string, ip Target, c chan<- Result, logger *slog.Logger) {
	logger.Debug("next up")

	result := Result{
		Target: ip,
		Listed: false,
		Rbl:    blocklist,
	}

	logger.Debug("about to query RBL")

	lookup := godnsbl.Reverse(ip.IP) + "." + result.Rbl
	logger.Debug("built lookup", slog.String("lookup", lookup))

	res, err := r.util.GetARecords(lookup)
	if err != nil {
		logger.Error("error occurred fetching A record", slog.String("msg", err.Error()))

		result.Error = true
		result.ErrorType = err
		c <- result
		return
	}

	if len(res) == 0 {
		// ip is not listed
		c <- result
		return
	}

	logger.Debug("ip is listed")

	result.Listed = true

	reason := net.ParseIP(res[0])
	if reason == nil {
		logger.Error("error getting (first) reason: %s", strings.Join(res, ", "))
		result.Error = true
		result.ErrorType = fmt.Errorf("error getting the (first) reason: %s", strings.Join(res, ", "))
		c <- result
		return
	}

	// fetch (potential) reason
	txt, err := r.util.GetTxtRecords(godnsbl.Reverse(reason) + "." + result.Rbl)
	if err != nil {
		logger.Error("error occurred fetching TXT record", slog.String("msg", err.Error()))

		result.Error = true
		result.ErrorType = err
		c <- result
		return
	}

	if len(txt) > 0 {
		result.Text = strings.Join(txt, ", ")
	}
	c <- result
}
