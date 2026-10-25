package ip_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/internal/ip"
	"github.com/stretchr/testify/assert"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func TestExpand(t *testing.T) {
	log := discardLogger()

	t.Run("dns", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"mail.example.org"}, log)
		assert.Len(t, hosts, 1)
	})

	t.Run("ip", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"8.8.8.8"}, log)
		assert.Len(t, hosts, 1)
	})

	t.Run("cidr: /24", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"1.2.3.4/24"}, log)
		assert.Len(t, hosts, 254)
	})

	t.Run("cidr: /32", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"1.1.1.1/32"}, log)
		assert.Len(t, hosts, 1)
		assert.Contains(t, hosts, "1.1.1.1")
	})

	t.Run("cidr: /31", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"192.168.1.0/31"}, log)
		assert.Len(t, hosts, 2)
		assert.Contains(t, hosts, "192.168.1.0")
		assert.Contains(t, hosts, "192.168.1.1")
	})

	t.Run("cidr: /30", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"10.0.0.0/30"}, log)
		// Excludes network (.0) and broadcast (.3)
		assert.Len(t, hosts, 2)
		assert.Contains(t, hosts, "10.0.0.1")
		assert.Contains(t, hosts, "10.0.0.2")
	})

	t.Run("cidr: /29", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"172.16.0.0/29"}, log)
		// 8 IPs - 2 (network/broadcast)
		assert.Len(t, hosts, 6)
	})

	t.Run("cidr: /28", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"10.1.1.0/28"}, log)
		// 16 IPs - 2 (network/broadcast)
		assert.Len(t, hosts, 14)
	})

	t.Run("cidr: /16", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"192.168.0.0/16"}, log)
		// 65536 IPs - 2 (network/broadcast)
		assert.Len(t, hosts, 65534)
	})
}

// TestExpandLargePrefixWarns asserts that a CIDR expanding to >= 1024
// hosts emits a WARN-level log record naming the offending CIDR and
// host count, and that smaller CIDRs stay silent.
func TestExpandLargePrefixWarns(t *testing.T) {
	t.Run("/24 stays silent", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelWarn}))

		ip.ExpandCIDRs([]string{"10.0.0.0/24"}, log)
		assert.Empty(t, buf.String(), "/24 (254 hosts) must not warn")
	})

	t.Run("/16 warns with cidr and host count", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelWarn}))

		ip.ExpandCIDRs([]string{"10.0.0.0/16"}, log)
		out := buf.String()
		assert.Contains(t, out, "level=WARN")
		assert.Contains(t, out, "cidr=10.0.0.0/16")
		assert.Contains(t, out, "hosts=65534")
	})
}
