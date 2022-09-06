package rbl_test

import (
	"os"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"
	"github.com/Luzilla/dnsbl_exporter/pkg/rbl"
	x "github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func TestUpdate(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	d := dns.New(new(x.Client), "0.0.0.0:53")

	r := rbl.New(d)
	r.Update("this.is.not.an.ip", []string{"cbl.abuseat.org"})

	if len(r.Results) > 0 {
		t.Errorf("Got a result, but shouldn't have: %v", r.Results)
	}
}