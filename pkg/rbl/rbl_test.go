package rbl_test

import (
	"os"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/dnsbl_exporter/pkg/rbl"
	x "github.com/miekg/dns"
	"golang.org/x/exp/slog"
)

func TestUpdate(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr))

	d := dns.New(new(x.Client), "0.0.0.0:53", logger)

	r := rbl.New(d, logger)
	r.Update("this.is.not.an.ip", []string{"cbl.abuseat.org"})

	if len(r.Results) > 0 {
		t.Errorf("Got a result, but shouldn't have: %v", r.Results)
	}
}
