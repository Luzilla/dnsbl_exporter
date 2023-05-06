package config_test

import (
	"os"
	"testing"

	"github.com/Luzilla/dnsbl_exporter/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		file    string
		key     string
		success bool
	}{
		{
			file:    "./../targets.ini",
			key:     "targets",
			success: true,
		},
		{
			file:    "./../rbls.ini",
			key:     "rbl",
			success: true,
		},
		{
			file:    "./does-not-exists.ini",
			key:     "foo",
			success: false,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.file, func(t *testing.T) {
			c := &config.Config{
				Logger: slog.New(slog.NewTextHandler(os.Stderr)),
			}
			_, err := c.LoadFile(tc.file)
			if tc.success {
				assert.NoError(t, err, "tc: "+tc.file)
			} else {
				assert.Error(t, err, "tc: "+tc.file)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	c := &config.Config{
		Logger: slog.New(slog.NewTextHandler(os.Stderr)),
	}

	cfg, err := c.LoadFile("./../targets.ini")
	assert.NoError(t, err)

	// ensure we return an error when the config section does not exist
	err = c.ValidateConfig(cfg, "blah")
	assert.Error(t, err)
}
