package dns

import (
	"fmt"
	"net"
	"strings"
	"time"

	x "github.com/miekg/dns"
	"log/slog"
)

type DNSUtil struct {
	client   *x.Client
	resolver struct {
		host string
		port string
	}
	logger *slog.Logger
}

func New(client *x.Client, resolver string, logger *slog.Logger) (*DNSUtil, error) {
	host, port, err := net.SplitHostPort(resolver)
	if err != nil {
		addrErr, ok := err.(*net.AddrError)
		if !ok {
			return nil, err
		}

		if !strings.Contains(addrErr.Error(), "missing port in address") {
			return nil, err
		}

		// default to port 53
		host = resolver
		port = "53"
	}

	// set timeouts
	client.DialTimeout = (2 * time.Second)
	client.ReadTimeout = (2 * time.Second)
	client.WriteTimeout = (2 * time.Second)

	return &DNSUtil{
		client: client,
		resolver: struct {
			host string
			port string
		}{
			host, port,
		},
		logger: logger,
	}, nil
}

// leaving this note for future me: maybe asking for As is not enough?
// what about CNAMEs, or AAAAs, etc..
func (d *DNSUtil) GetARecords(target string) (list []string, err error) {
	msg := createQuestion(target, x.TypeA)

	result, err := d.makeQuery(msg)
	if err != nil {
		return
	}

	if len(result.Answer) == 0 {
		return
	}

	for _, ans := range result.Answer {
		if t, ok := ans.(*x.A); ok {
			d.logger.Debug("we have an A-record", slog.String("target", target), slog.String("v", t.A.String()))
			list = append(list, t.A.String())
		}
	}

	return
}

func (d *DNSUtil) GetTxtRecords(target string) (list []string, err error) {
	msg := createQuestion(target, x.TypeTXT)

	result, err := d.makeQuery(msg)
	if err != nil {
		return
	}

	if len(result.Answer) == 0 {
		return
	}

	for _, ans := range result.Answer {
		if t, ok := ans.(*x.TXT); ok {
			list = append(list, t.Txt...)
		}
	}

	return
}

func (d *DNSUtil) makeQuery(msg *x.Msg) (*x.Msg, error) {
	result, rt, err := d.client.Exchange(msg, net.JoinHostPort(d.resolver.host, d.resolver.port))
	if err != nil {
		return nil, err
	}

	d.logger.Debug("roundtrip",
		slog.String("question", msg.Question[0].String()),
		slog.String("rcode", x.RcodeToString[result.Rcode]),
		slog.Float64("seconds", rt.Seconds())) // fixme -> histogram

	// RFC 5782: NOERROR + records means listed; NOERROR with no records or
	// NXDOMAIN means not listed. Anything else (SERVFAIL, REFUSED, …) is an
	// outage, not a clean check, and must surface as an error.
	switch result.Rcode {
	case x.RcodeSuccess, x.RcodeNameError:
		return result, nil
	default:
		return result, fmt.Errorf("dns: %s", x.RcodeToString[result.Rcode])
	}
}

func createQuestion(target string, record uint16) *x.Msg {
	msg := new(x.Msg)
	msg.SetQuestion(x.Fqdn(target), record)

	return msg
}
