package rbl

import (
	"net"

	"github.com/Luzilla/dnsbl_exporter/pkg/dns"

	"golang.org/x/exp/slog"
)

type Target struct {
	Host string
	IP   net.IP
}

type Resolver struct {
	logger *slog.Logger
	util   *dns.DNSUtil
}

func NewRBLResolver(logger *slog.Logger, util *dns.DNSUtil) *Resolver {
	return &Resolver{
		logger: logger,
		util:   util,
	}
}

func (r *Resolver) Do(target string, c chan<- Target, done func()) {
	defer done()

	addr := net.ParseIP(target)
	if addr != nil {
		// already an IP
		r.logger.Info("we had an ip already", slog.String("ip", target))
		c <- Target{
			Host: target,
			IP:   addr,
		}
		return
	}

	ipsA, err := r.util.GetARecords(target)
	if err != nil {
		r.logger.Error("error fetching A-records for target", slog.String("msg", err.Error()))
		return
	}

	for _, i := range ipsA {
		a := net.ParseIP(i)
		if a == nil {
			r.logger.Error("address failed parsing", slog.String("ip", i))
			continue
		}
		c <- Target{
			Host: target,
			IP:   a,
		}
	}
}
