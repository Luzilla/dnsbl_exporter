package resolvconf_test

import (
	"path/filepath"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/internal/resolvconf"
	"github.com/stretchr/testify/assert"
)

func TestGetServers(t *testing.T) {
	testCases := []struct {
		path     string
		expected []string
	}{
		{
			path:     "no-server.conf",
			expected: []string(nil),
		},
		{
			path:     "two-servers.conf",
			expected: []string{"1.1.1.1", "8.8.8.8"},
		},
	}

	for _, tc := range testCases {
		servers, err := resolvconf.GetServers(filepath.Join("fixtures", tc.path))
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, servers)
	}
}
