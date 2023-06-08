package dns_test

import (
	"testing"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	x "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestDNSSuite(t *testing.T) {
	t.Run("test=a", func(t *testing.T) {
		d := dns.New(new(x.Client), "0.0.0.0:53", slog.Default())

		aRecords, err := d.GetARecords("relay.heise.de")

		assert.NoError(t, err)
		assert.Greater(t, len(aRecords), 0)
	})

	t.Run("test=txt", func(t *testing.T) {
		d := dns.New(new(x.Client), "0.0.0.0:53", slog.Default())

		// DTAG dial-up
		txtRecords, err := d.GetTxtRecords("85.198.214.79.zen.spamhaus.org")
		assert.NoError(t, err)
		assert.Greater(t, len(txtRecords), 0)
	})
}
