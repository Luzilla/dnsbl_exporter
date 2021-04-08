package collector

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestUpdate(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	rbl := NewRbl("0.0.0.0:53")
	rbl.Update("this.is.not.an.ip", []string{"cbl.abuseat.org"})

	if len(rbl.Results) > 0 {
		t.Errorf("Got a result, but shouldn't have: %v", rbl.Results)
	}
}
