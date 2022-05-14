package collector

import (
	"fmt"
	"net"
	"sync"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/godnsbl"
	log "github.com/sirupsen/logrus"
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
}

// NewRbl ... factory
func NewRbl(util *dns.DNSUtil) Rbl {
	var results []Rblresult

	rbl := Rbl{
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

			log.Debugf("Working blacklist %s (ip: %s)", source, ip)

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

	log.Debugf("Trying to query blacklist '%s' for %s", blacklist, ip)

	lookup := fmt.Sprintf("%s.%s", ip, blacklist)

	res, err := rbl.util.GetARecords(lookup)
	if len(res) > 0 {
		result.Listed = true

		txt, _ := net.LookupTXT(lookup)
		if len(txt) > 0 {
			result.Text = txt[0]
		}
	}

	if err != nil {
		result.Error = true
		result.ErrorType = err
	}

}

func (rbl *Rbl) lookup(rblList string, targetHost string) []Rblresult {
	var ips []string

	addr := net.ParseIP(targetHost)
	if addr == nil {
		ipsA, err := rbl.util.GetARecords(targetHost)
		if err != nil {
			log.Errorln(err)
		}

		ips = ipsA
	} else {
		log.Infoln("We had an ip", addr.String())
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
			log.Errorf("Unable to parse IP: %s", ip)
			continue
		}

		// reverse it, for the look up
		revIP := godnsbl.Reverse(ValidIPAddress)

		rbl.query(revIP, rblList, &res)

		rbl.Results = append(rbl.Results, res)
	}

	return rbl.Results
}
