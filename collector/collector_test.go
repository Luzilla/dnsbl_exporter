package collector_test

import (
	"strings"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/collector"
	"github.com/Luzilla/dnsbl_exporter/internal/tests"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollectorSuite(t *testing.T) {
	dnsMock := tests.CreateDNSMock(t)
	defer dnsMock.Close()
	logger := tests.CreateTestLogger(t)
	util := tests.CreateDNSUtil(t, dnsMock.LocalAddr())

	t.Run("test=ip-based", func(t *testing.T) {
		rbls := []string{"zen.spamhaus.org", "cbl.abuseat.org"}
		targets := []string{
			"79.214.198.85",  // bad
			"relay.heise.de", // good
			"1.3.3.7",        // good
			"1.3.3.7/30",     // good
		}

		c := collector.NewRblCollector(rbls, targets, false, util, logger)

		result, err := testutil.CollectAndLint(c)
		assert.Empty(t, result)
		assert.NoError(t, err)

		// take all metrics but duration as it's value is hardly predictable
		metrics := []string{}
		for _, metric := range []string{"used", "ips_blacklisted", "errors", "listed", "targets"} {
			metrics = append(metrics, collector.BuildFQName(metric))
		}
		expected := `
      # HELP luzilla_rbls_ips_blacklisted Blacklisted IPs
      # TYPE luzilla_rbls_ips_blacklisted gauge
      luzilla_rbls_ips_blacklisted{hostname="1.3.3.4",ip="1.3.3.4",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="1.3.3.4",ip="1.3.3.4",rbl="zen.spamhaus.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="1.3.3.5",ip="1.3.3.5",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="1.3.3.5",ip="1.3.3.5",rbl="zen.spamhaus.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="1.3.3.6",ip="1.3.3.6",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="1.3.3.6",ip="1.3.3.6",rbl="zen.spamhaus.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="1.3.3.7",ip="1.3.3.7",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="1.3.3.7",ip="1.3.3.7",rbl="zen.spamhaus.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="79.214.198.85",ip="79.214.198.85",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="79.214.198.85",ip="79.214.198.85",rbl="zen.spamhaus.org"} 1
      luzilla_rbls_ips_blacklisted{hostname="relay.heise.de",ip="193.99.145.50",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_ips_blacklisted{hostname="relay.heise.de",ip="193.99.145.50",rbl="zen.spamhaus.org"} 0
      # HELP luzilla_rbls_listed The number of listings in RBLs (this is bad)
      # TYPE luzilla_rbls_listed gauge
      luzilla_rbls_listed{rbl="cbl.abuseat.org"} 0
      luzilla_rbls_listed{rbl="zen.spamhaus.org"} 1
      # HELP luzilla_rbls_targets The number of targets that are being probed (configured via targets.ini or ?target=)
      # TYPE luzilla_rbls_targets gauge
      luzilla_rbls_targets 6
      # HELP luzilla_rbls_used The number of RBLs to check IPs against (configured via rbls.ini)
      # TYPE luzilla_rbls_used gauge
      luzilla_rbls_used 2
    `
		err = testutil.CollectAndCompare(c, strings.NewReader(expected), metrics...)
		assert.NoError(t, err)
	})

	t.Run("test=domain-based", func(t *testing.T) {
		rbls := []string{"dbl.spamhaus.org"}
		targets := []string{
			"dbltest.com", // bad
			"example.com", // good
		}

		c := collector.NewRblCollector(rbls, targets, true, util, logger)

		result, err := testutil.CollectAndLint(c)
		assert.Empty(t, result)
		assert.NoError(t, err)

		// take all metrics but duration as it's value is hardly predictable
		metrics := []string{}
		for _, metric := range []string{"used", "ips_blacklisted", "errors", "listed", "targets"} {
			metrics = append(metrics, collector.BuildFQName(metric))
		}
		expected := `
			# HELP luzilla_rbls_ips_blacklisted Blacklisted IPs
			# TYPE luzilla_rbls_ips_blacklisted gauge
			luzilla_rbls_ips_blacklisted{hostname="dbltest.com",ip="127.0.1.2",rbl="dbl.spamhaus.org"} 1
			luzilla_rbls_ips_blacklisted{hostname="example.com",ip="",rbl="dbl.spamhaus.org"} 0
			# HELP luzilla_rbls_listed The number of listings in RBLs (this is bad)
			# TYPE luzilla_rbls_listed gauge
			luzilla_rbls_listed{rbl="dbl.spamhaus.org"} 1
			# HELP luzilla_rbls_targets The number of targets that are being probed (configured via targets.ini or ?target=)
			# TYPE luzilla_rbls_targets gauge
			luzilla_rbls_targets 2
			# HELP luzilla_rbls_used The number of RBLs to check IPs against (configured via rbls.ini)
			# TYPE luzilla_rbls_used gauge
			luzilla_rbls_used 1
    `
		err = testutil.CollectAndCompare(c, strings.NewReader(expected), metrics...)
		assert.NoError(t, err)
	})
}
