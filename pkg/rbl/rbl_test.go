package rbl_test

import (
	"net"
	"sync"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/internal/tests"
	"github.com/Luzilla/dnsbl_exporter/pkg/rbl"
	"github.com/stretchr/testify/assert"
)

func TestRblSuite(t *testing.T) {
	t.Run("run=valid", func(t *testing.T) {
		dnsMock := tests.CreateDNSMock(t)
		defer dnsMock.Close()

		logger := tests.CreateTestLogger(t)
		d := tests.CreateDNSUtil(t, dnsMock.LocalAddr())
		r := rbl.New(d, logger)

		resolver := rbl.NewRBLResolver(logger, d)

		targets := make(chan rbl.Target)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			wg.Wait()
			close(targets)
		}()
		go resolver.Do("relay.heise.de", targets, wg.Done)

		for ip := range targets {
			for _, blocklist := range []string{"cbl.abuseat.org", "zen.spamhaus.org"} {
				c := make(chan rbl.Result)
				r.Update(ip, blocklist, c)
				res := <-c
				close(c)

				assert.False(t, res.Error)
				assert.NoError(t, res.ErrorType)
			}
		}

		// assert.Equal(t, 1, len(r.Results), "Got more than one result, but shouldn't have: %v", r.Results)
	})

	t.Run("run=multiple", func(t *testing.T) {
		dnsMock := tests.CreateDNSMock(t)
		defer dnsMock.Close()

		logger := tests.CreateTestLogger(t)
		d := tests.CreateDNSUtil(t, dnsMock.LocalAddr())
		r := rbl.New(d, logger)

		hosts := []rbl.Target{
			{
				Host: "127.0.0.4",
				IP:   net.ParseIP("127.0.0.4"),
			},
			{
				Host: "127.0.0.2",
				IP:   net.ParseIP("127.0.0.2"),
			},
			{
				Host: "127.0.0.10",
				IP:   net.ParseIP("127.0.0.10"),
			},
		}

		// 3 hosts, and 2 RBLs => 6 results
		results := make([]rbl.Result, 0)
		for _, ip := range hosts {
			for _, blocklist := range []string{"cbl.abuseat.org", "zen.spamhaus.org"} {
				c := make(chan rbl.Result)
				defer close(c)

				r.Update(ip, blocklist, c)
				results = append(results, <-c)
			}

		}

		assert.Len(t, results, 6)
	})

	t.Run("run=error_result", func(t *testing.T) {
		dnsMock := tests.CreateDNSMock(t)
		defer dnsMock.Close()

		logger := tests.CreateTestLogger(t)
		r := rbl.New(tests.CreateDNSUtil(t, dnsMock.LocalAddr()), logger)
		c := make(chan rbl.Result)

		// d-tag dial-up IP
		target := rbl.Target{
			Host: "79.214.198.85",
			IP:   net.ParseIP("79.214.198.85"),
		}
		r.Update(target, "zen.spamhaus.org", c)

		result := <-c

		// assert the right RBL is in there
		assert.Equal(t, "zen.spamhaus.org", result.Rbl)

		// this is not an error as in transport/dialer
		assert.False(t, result.Error)
		assert.NoError(t, result.ErrorType)

		// but the IP is listed
		assert.True(t, result.Listed)
		assert.Contains(t, result.Text, "https://www.spamhaus.org/")
	})
}

func TestResolver(t *testing.T) {
	dnsMock := tests.CreateDNSMock(t)
	defer dnsMock.Close()

	logger := tests.CreateTestLogger(t)

	resolver := rbl.NewRBLResolver(logger, tests.CreateDNSUtil(t, dnsMock.LocalAddr()))

	c := make(chan rbl.Target)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Wait()
		close(c)
	}()
	go resolver.Do("relay.heise.de", c, wg.Done)

	for ip := range c {
		assert.Equal(t, "relay.heise.de", ip.Host)
		assert.NotEmpty(t, ip.IP.String())
	}
}
