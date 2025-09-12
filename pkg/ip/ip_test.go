package ip_test

import (
	"testing"

	"github.com/Luzilla/dnsbl_exporter/pkg/ip"
	"github.com/stretchr/testify/assert"
)

func TestExpand(t *testing.T) {
	t.Run("dns", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"mail.example.org"})
		assert.Len(t, hosts, 1)
	})

	t.Run("ip", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"8.8.8.8"})
		assert.Len(t, hosts, 1)
	})

	t.Run("cidr: /24", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"1.2.3.4/24"})
		assert.Len(t, hosts, 254)
	})

	t.Run("cidr: /32", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"1.1.1.1/32"})
		assert.Len(t, hosts, 1)
		assert.Contains(t, hosts, "1.1.1.1")
	})

	t.Run("cidr: /31", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"192.168.1.0/31"})
		assert.Len(t, hosts, 2)
		assert.Contains(t, hosts, "192.168.1.0")
		assert.Contains(t, hosts, "192.168.1.1")
	})

	t.Run("cidr: /30", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"10.0.0.0/30"})
		// Excludes network (.0) and broadcast (.3)
		assert.Len(t, hosts, 2)
		assert.Contains(t, hosts, "10.0.0.1")
		assert.Contains(t, hosts, "10.0.0.2")
	})

	t.Run("cidr: /29", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"172.16.0.0/29"})
		// 8 IPs - 2 (network/broadcast)
		assert.Len(t, hosts, 6)
	})

	t.Run("cidr: /28", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"10.1.1.0/28"})
		// 16 IPs - 2 (network/broadcast)
		assert.Len(t, hosts, 14)
	})

	t.Run("cidr: /16", func(t *testing.T) {
		hosts := ip.ExpandCIDRs([]string{"192.168.0.0/16"})
		// 65536 IPs - 2 (network/broadcast)
		assert.Len(t, hosts, 65534)
	})
}
