package ip

import (
	"encoding/binary"
	"log/slog"
	"net"
	"slices"
)

// largeExpansionThreshold is the number of hosts a single CIDR may expand
// to before we log a warning. /22 = 1022 hosts (silent), /21 = 2046 hosts
// (warns). A /16 (65534 hosts) trivially trips this and produces a scrape
// load that will get the resolver rate-limited by most public RBLs.
const largeExpansionThreshold = 1024

func getIPsFromCIDR(ipv4Net *net.IPNet) []string {
	// Convert IPNet mask and IP address to uint32
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// Compute the correct network and broadcast addresses
	network := start & mask      // Network address (e.g., 1.2.3.0)
	broadcast := network | ^mask // Broadcast address (e.g., 1.2.3.255)

	// For /32 and /31, handle specially
	if network == broadcast {
		// /32 - single IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, network)
		return []string{ip.String()}
	}
	
	if broadcast-network == 1 {
		// /31 - two IPs, both usable
		ips := make([]string, 2)
		ip1 := make(net.IP, 4)
		ip2 := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip1, network)
		binary.BigEndian.PutUint32(ip2, broadcast)
		ips[0] = ip1.String()
		ips[1] = ip2.String()
		return ips
	}

	// Exclude network and broadcast addresses for other CIDRs
	ips := make([]string, 0, broadcast-network-1)
	for i := network + 1; i < broadcast; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip.String())
	}

	return ips
}

func ExpandCIDRs(hosts []string, logger *slog.Logger) (newHosts []string) {
	for _, host := range hosts {
		// If it's a CIDR
		if _, ipNet, err := net.ParseCIDR(host); err == nil {
			ips := getIPsFromCIDR(ipNet)
			if len(ips) >= largeExpansionThreshold {
				logger.Warn("CIDR expands to a large number of hosts; scrape load and resolver rate-limiting risk",
					slog.String("cidr", host),
					slog.Int("hosts", len(ips)),
				)
			}
			for _, ip := range ips {
				// Only add IP to slice if it doesn't already exist
				if !slices.Contains(newHosts, ip) {
					newHosts = append(newHosts, ip)
				}
			}
		} else {
			// If it's not a CIDR, just check and add the host itself
			if !slices.Contains(newHosts, host) {
				newHosts = append(newHosts, host)
			}
		}
	}

	return
}
