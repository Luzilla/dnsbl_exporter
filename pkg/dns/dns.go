package dns

import (
	"net"

	x "github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type DNSUtil struct {
	client   *x.Client
	resolver string
}

func New(client *x.Client, resolver string) *DNSUtil {
	return &DNSUtil{
		client:   client,
		resolver: resolver,
	}
}

// leaving this note for future me: maybe asking for As is not enough?
// what about CNAMEs, or AAAAs, etc..
func (d *DNSUtil) GetARecords(target string) ([]string, error) {
	msg := createQuestion(target, x.TypeA)

	result, err := d.makeQuery(msg)

	var list []string

	if err == nil && len(result.Answer) > 0 {
		for _, ans := range result.Answer {
			if t, ok := ans.(*x.A); ok {
				log.Debugf("We have an A-Record %s for %s", t.A.String(), target)
				list = append(list, t.A.String())
			}
		}
	}

	return list, err
}

func (d *DNSUtil) GetTxtRecords(target string) ([]string, error) {
	msg := createQuestion(target, x.TypeTXT)

	result, err := d.makeQuery(msg)

	var list []string

	if len(result.Answer) > 0 {
		for _, ans := range result.Answer {
			if t, ok := ans.(*x.TXT); ok {
				list = append(list, t.Txt...)
			}
		}
	}

	return list, err
}

func (d *DNSUtil) makeQuery(msg *x.Msg) (*x.Msg, error) {
	host, port, err := net.SplitHostPort(d.resolver)
	if err != nil {
		if err.Error() == "missing port in address" {
			port = "53"
		}
	}

	result, rt, err := d.client.Exchange(msg, net.JoinHostPort(host, port))
	log.Debugln("Roundtrip", rt) // fixme -> histogram

	return result, err
}

func createQuestion(target string, record uint16) *x.Msg {
	msg := new(x.Msg)
	msg.SetQuestion(x.Fqdn(target), record)

	return msg
}
