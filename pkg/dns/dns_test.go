package dns_test

import (
	"testing"

	"github.com/Luzilla/dnsbl_exporter/internal/tests"
	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	x "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestDNSSuite(t *testing.T) {
	t.Run("test=a", func(t *testing.T) {
		dnsMock := tests.CreateDNSMock(t)
		defer dnsMock.Close()

		d := tests.CreateDNSUtil(t, dnsMock.LocalAddr())

		aRecords, err := d.GetARecords("relay.heise.de")

		assert.NoError(t, err)
		assert.Greater(t, len(aRecords), 0)
	})

	t.Run("test=txt", func(t *testing.T) {
		dnsMock := tests.CreateDNSMock(t)
		defer dnsMock.Close()

		d := tests.CreateDNSUtil(t, dnsMock.LocalAddr())

		txtRecords, err := d.GetTxtRecords("10.0.0.127.zen.spamhaus.org")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(txtRecords))
	})
}

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		addr string
		err  bool
	}{
		{"0.0.0.0", false},    // assert the port gets added
		{"0.0.0.0:53", false}, // standard input
	} {
		_, err := dns.New(new(x.Client), tc.addr, tests.CreateTestLogger(t))
		assert.NoError(t, err, tc.addr)
	}

}
