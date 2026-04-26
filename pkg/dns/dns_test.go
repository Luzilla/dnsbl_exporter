package dns_test

import (
	"errors"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/internal/tests"
	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/foxcpp/go-mockdns"
	x "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDNSSuite(t *testing.T) {
	t.Run("test=a", func(t *testing.T) {
		dnsMock := tests.CreateDNSMock(t)
		defer dnsMock.Close()

		d := tests.CreateDNSUtil(t, dnsMock.LocalAddr())

		aRecords, err := d.GetARecords("relay.heise.de")
		require.NoError(t, err)
		assert.Greater(t, len(aRecords), 0)
	})

	t.Run("test=txt", func(t *testing.T) {
		dnsMock := tests.CreateDNSMock(t)
		defer dnsMock.Close()

		d := tests.CreateDNSUtil(t, dnsMock.LocalAddr())

		txtRecords, err := d.GetTxtRecords("10.0.0.127.zen.spamhaus.org")
		require.NoError(t, err)
		assert.Equal(t, 1, len(txtRecords))
	})
}

// TestRcodeAsError: a DNS response carrying a non-success rcode (SERVFAIL,
// REFUSED, ...) must surface as an error so the `errors` metric fires.
// Per RFC 5782, NXDOMAIN means "not listed" and must NOT surface as an
// error.
func TestRcodeAsError(t *testing.T) {
	t.Run("servfail-surfaces-as-error", func(t *testing.T) {
		// mockdns returns SERVFAIL when Zone.Err is non-nil.
		srv, err := mockdns.NewServer(map[string]mockdns.Zone{
			"broken.example.com.": {Err: errors.New("simulated outage")},
		}, true)
		require.NoError(t, err)
		defer srv.Close()

		d := tests.CreateDNSUtil(t, srv.LocalAddr())

		_, err = d.GetARecords("broken.example.com")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "SERVFAIL")
	})

	t.Run("nxdomain-stays-clean", func(t *testing.T) {
		// Unknown names in the mock yield NXDOMAIN; per RFC 5782 that's
		// "not listed", not an error.
		dnsMock := tests.CreateDNSMock(t)
		defer dnsMock.Close()

		d := tests.CreateDNSUtil(t, dnsMock.LocalAddr())

		records, err := d.GetARecords("not.in.mock.example.com")
		require.NoError(t, err)
		assert.Empty(t, records)
	})
}

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		addr string
		err  bool
	}{
		{"0.0.0.0", false},      // assert the port gets added
		{"0.0.0.0:53", false},   // standard input
		{"unbound:5353", false}, // tests
	} {
		_, err := dns.New(new(x.Client), tc.addr, tests.CreateTestLogger(t))
		assert.NoError(t, err, tc.addr)
	}
}
