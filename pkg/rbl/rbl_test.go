package rbl_test

import (
	"os"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/dnsbl_exporter/pkg/rbl"
	x "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestRblSuite(t *testing.T) {
	t.Run("run=invalid_ip", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stderr))

		d := dns.New(new(x.Client), "0.0.0.0:53", logger)

		r := rbl.New(d, logger)
		r.Update("this.is.not.an.ip", []string{"cbl.abuseat.org"})

		assert.Equal(t, 0, len(r.Results), "Got a result, but shouldn't have: %v", r.Results)
	})

	t.Run("run=error_result", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stderr))

		d := dns.New(new(x.Client), "0.0.0.0:53", logger)

		r := rbl.New(d, logger)
		r.Update("79.214.198.85", []string{"zen.spamhaus.org"})

		assert.Equal(t, 1, len(r.Results))

		result := r.Results[0]

		// assert the right RBL is in there
		assert.Equal(t, "zen.spamhaus.org", result.Rbl)

		// this is not an error as in transport/dialer
		assert.False(t, result.Error)
		assert.NoError(t, result.ErrorType)

		// but the IP is listed
		assert.True(t, result.Listed)
		assert.Equal(t, "https://www.spamhaus.org/query/ip/79.214.198.85", result.Text)
	})
}
