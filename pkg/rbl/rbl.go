package rbl

import (
	"net"
	"sync"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/godnsbl"
	"golang.org/x/exp/slog"
)

// Rblresult extends godnsbl and adds RBL name
type Rblresult struct {
	Address   string
	Listed    bool
	Text      string
	Error     bool
	ErrorType error
	Rbl       string
	Target    string
}

// Rbl ... object
type Rbl struct {
	Results []Rblresult
	util    *dns.DNSUtil
	logger  *slog.Logger
}

// NewRbl ... factory
func New(util *dns.DNSUtil, logger *slog.Logger) Rbl {
	var results []Rblresult

	rbl := Rbl{
		logger:  logger,
		util:    util,
		Results: results,
	}

	return rbl
}

// Update runs the checks for an IP against against all "rbls"
func (rbl *Rbl) Update(ip string, rbls []string) {
	// from: godnsbl
	wg := &sync.WaitGroup{}

	for _, source := range rbls {
		wg.Add(1)
		go func(source string, ip string) {
			defer wg.Done()

			rbl.logger.Debug("Next up", slog.String("rbl", source), slog.String("ip", ip))

			results := rbl.lookup(source, ip)
			if len(results) == 0 {
				rbl.Results = []Rblresult{}
			} else {
				rbl.Results = results
			}
		}(source, ip)
	}

	wg.Wait()
}

func (rbl *Rbl) query(ip string, blacklist string, result *Rblresult) {
	result.Listed = false

	rbl.logger.Debug("About to query RBL", slog.String("rbl", blacklist), slog.String("ip", ip))

	lookup := ip + "." + blacklist
	rbl.logger.Debug("Built lookup", slog.String("lookup", lookup))

	res, err := rbl.util.GetARecords(lookup)
	if err != nil {
		rbl.logger.Error("error occurred fetching A record", slog.String("msg", err.Error()))
		result.Error = true
		result.ErrorType = err
		return
	}

	if len(res) == 0 {
		// ip is not listed
		return
	}

	rbl.logger.Debug("ip is listed", slog.String("ip", ip))
	result.Listed = true

	// fetch (potential) reason
	txt, err := net.LookupTXT(lookup)
	if err != nil {
		rbl.logger.Error("error occurred fetching TXT record", slog.String("msg", err.Error()))
		return
	}

	if len(txt) > 0 {
		result.Text = txt[0]
	}
}

func (rbl *Rbl) lookup(rblList string, targetHost string) []Rblresult {
	var ips []string

	addr := net.ParseIP(targetHost)
	if addr == nil {
		ipsA, err := rbl.util.GetARecords(targetHost)
		if err != nil {
			rbl.logger.Error(err.Error())
			return rbl.Results
		}

		ips = ipsA
	} else {
		rbl.logger.Info("We had an ip", slog.String("ip", addr.String()))
		ips = append(ips, addr.String())
	}

	for _, ip := range ips {
		res := Rblresult{}
		res.Target = targetHost
		res.Address = ip
		res.Rbl = rblList

		// attempt to "validate" the IP
		ValidIPAddress := net.ParseIP(ip)
		if ValidIPAddress == nil {
			rbl.logger.Error("Unable to parse IP", slog.String("ip", ip))
			continue
		}

		// reverse it, for the look up
		revIP := godnsbl.Reverse(ValidIPAddress)

		rbl.query(revIP, rblList, &res)

		rbl.Results = append(rbl.Results, res)
	}

	return rbl.Results
}
