package tests

import (
	"net"
	"os"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/foxcpp/go-mockdns"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"

	x "github.com/miekg/dns"
)

func CreateTestLogger(t *testing.T) *slog.Logger {
	t.Helper()
	return slog.New(slog.NewTextHandler(os.Stderr))
}

func CreateDNSUtil(t *testing.T, resolver net.Addr) *dns.DNSUtil {
	t.Helper()

	dns, err := dns.New(new(x.Client), resolver.String(), CreateTestLogger(t))
	if err != nil {
		assert.FailNow(t, "unable to create DNSUtil", "error: %s", err)
	}
	return dns
}

func CreateDNSMock(t *testing.T) *mockdns.Server {
	t.Helper()

	srv, err := mockdns.NewServer(map[string]mockdns.Zone{
		"relay.heise.de.": {
			A: []string{"193.99.145.50"},
		},
		"85.198.214.79.zen.spamhaus.org.": {
			A: []string{"127.0.0.10"},
		},
		// rbl responses
		"4.0.0.127.zen.spamhaus.org.": {
			TXT: []string{"https://www.spamhaus.org/query/ip/127.0.0.4"},
		},
		"4.0.0.127.cbl.abuseat.org.": {
			TXT: []string{"Error: open resolver; https://www.spamhaus.org/returnc/pub/2400:cb00:67:1024::a29e:713c"},
		},
		"2.0.0.127.zen.spamhaus.org.": {
			TXT: []string{"https://www.spamhaus.org/query/ip/127.0.0.2"},
		},
		"2.0.0.127.cbl.abuseat.org.": {
			TXT: []string{"https://www.spamhaus.org/query/ip/127.0.0.2"},
		},
		"10.0.0.127.zen.spamhaus.org.": {
			TXT: []string{"https://www.spamhaus.org/query/ip/127.0.0.10"},
		},
	}, true)
	if err != nil {
		assert.FailNow(t, "failed building mock", "error: %s", err)
	}

	return srv
}
