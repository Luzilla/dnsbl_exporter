package collector

import (
	"fmt"
	"net"
	"sync"

	"github.com/Luzilla/godnsbl"
	"github.com/miekg/dns"
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
	Resolver  string
	Results   []Rblresult
	DNSClient *dns.Client
}

// NewRbl ... factory
func NewRbl(resolver string) Rbl {
	client := new(dns.Client)

	rbl := Rbl{
		Resolver:  resolver,
		DNSClient: client,
	}

	return rbl
}

func (rbl *Rbl) createQuestion(target string, record uint16) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(target), record)

	return msg
}

func (rbl *Rbl) makeQuery(msg *dns.Msg) (*dns.Msg, error) {
	host, port, err := net.SplitHostPort(rbl.Resolver)
	if err != nil {
		if err.Error() == "missing port in address" {
			port = "53"
		}
	}

	result, rt, err := rbl.DNSClient.Exchange(msg, net.JoinHostPort(host, port))
	log.Debugln("Roundtrip", rt) // fixme -> histogram

	return result, err
}

func (rbl *Rbl) getARecords(target string) ([]string, error) {
	msg := rbl.createQuestion(target, dns.TypeA)

	result, err := rbl.makeQuery(msg)

	var list []string

	if err == nil && len(result.Answer) > 0 {
		for _, ans := range result.Answer {
			if t, ok := ans.(*dns.A); ok {
				log.Debugf("We have an A-Record %s for %s", t.A.String(), target)
				list = append(list, t.A.String())
			}
		}
	}

	return list, err
}

func (rbl *Rbl) getTxtRecords(target string) ([]string, error) {
	msg := rbl.createQuestion(target, dns.TypeTXT)

	result, err := rbl.makeQuery(msg)

	var list []string

	if len(result.Answer) > 0 {
		for _, ans := range result.Answer {
			if t, ok := ans.(*dns.TXT); ok {
				for _, txt := range t.Txt {
					list = append(list, txt)
				}
			}
		}
	}

	return list, err
}

func (rbl *Rbl) query(ip string, blacklist string, result *Rblresult) {
	result.Listed = false

	log.Debugf("Trying to query blacklist '%s' for %s", blacklist, ip)

	lookup := fmt.Sprintf("%s.%s", ip, blacklist)

	res, err := rbl.getARecords(lookup)
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
		ipsA, err := rbl.getARecords(targetHost)
		if err != nil {
			log.Debugln(err)
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

		revIP := godnsbl.Reverse(net.ParseIP(ip))

		rbl.query(revIP, rblList, &res)

		rbl.Results = append(rbl.Results, res)
	}

	return rbl.Results
}

// Update runs the checks for an against against all "rbls"
func (rbl *Rbl) Update(ip string, rbls []string) {
	// from: godnsbl
	wg := &sync.WaitGroup{}

	for _, source := range rbls {

		wg.Add(1)
		go func(source string, ip string) {
			defer wg.Done()

			log.Debugf("Working blacklist %s", source)

			results := rbl.lookup(source, ip)
			if len(results) == 0 {
				rbl.Results = []Rblresult{}
			} else {
				rbl.Results = results
			}
		}(source, ip)
	}

	wg.Wait()

	return
}
