package collector_test

import (
	"net"
	"strings"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/collector"
	"github.com/Luzilla/dnsbl_exporter/internal/tests"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
      # HELP luzilla_rbls_errors Whether an error occurred while testing this target against the RBL (1) or not (0)
      # TYPE luzilla_rbls_errors gauge
      luzilla_rbls_errors{hostname="1.3.3.5",ip="1.3.3.5",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_errors{hostname="1.3.3.5",ip="1.3.3.5",rbl="zen.spamhaus.org"} 0
      luzilla_rbls_errors{hostname="1.3.3.6",ip="1.3.3.6",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_errors{hostname="1.3.3.6",ip="1.3.3.6",rbl="zen.spamhaus.org"} 0
      luzilla_rbls_errors{hostname="1.3.3.7",ip="1.3.3.7",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_errors{hostname="1.3.3.7",ip="1.3.3.7",rbl="zen.spamhaus.org"} 0
      luzilla_rbls_errors{hostname="79.214.198.85",ip="79.214.198.85",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_errors{hostname="79.214.198.85",ip="79.214.198.85",rbl="zen.spamhaus.org"} 0
      luzilla_rbls_errors{hostname="relay.heise.de",ip="193.99.145.50",rbl="cbl.abuseat.org"} 0
      luzilla_rbls_errors{hostname="relay.heise.de",ip="193.99.145.50",rbl="zen.spamhaus.org"} 0
      # HELP luzilla_rbls_ips_blacklisted Blacklisted IPs
      # TYPE luzilla_rbls_ips_blacklisted gauge
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
      luzilla_rbls_targets 5
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
			# HELP luzilla_rbls_errors Whether an error occurred while testing this target against the RBL (1) or not (0)
			# TYPE luzilla_rbls_errors gauge
			luzilla_rbls_errors{hostname="dbltest.com",ip="127.0.1.2",rbl="dbl.spamhaus.org"} 0
			luzilla_rbls_errors{hostname="example.com",ip="",rbl="dbl.spamhaus.org"} 0
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

// TestCollectorErrorsHaveUniqueLabelSets: when more than one target errors
// against the same RBL in a single scrape, the `errors` series must stay
// unique. Previously the metric only carried the `rbl` label, so two errors
// against the same RBL produced duplicate label sets and Gather rather
// "collected before with the same name and label values". The fix gives
// `errors` the same `(rbl, ip, hostname)` shape as `ips_blacklisted`,
// so each (target, rbl) pair is its own series.
func TestCollectorErrorsHaveUniqueLabelSets(t *testing.T) {
	// Bind a UDP socket and close it immediately so DNS queries to its
	// address fail fast (connection refused via ICMP on localhost).
	// Every lookup against the configured RBL hits the error branch in
	// rbl.lookup, so two targets produce two error emissions for the same
	// `rbl` label set.
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	require.NoError(t, err)
	deadAddr := pc.LocalAddr()
	require.NoError(t, pc.Close())

	util := tests.CreateDNSUtil(t, deadAddr)
	logger := tests.CreateTestLogger(t)

	rbls := []string{"zen.spamhaus.org"}
	// IP targets (no resolver step) so each one drives exactly one
	// failing A-record lookup against the dead resolver.
	targets := []string{"127.0.0.1", "127.0.0.2"}

	c := collector.NewRblCollector(rbls, targets, false, util, logger)

	registry := prometheus.NewRegistry()
	registry.MustRegister(c)

	_, err = registry.Gather()
	assert.NoError(t, err, "errors must produce one unique series per (target, rbl), not collide on the rbl label")
}
