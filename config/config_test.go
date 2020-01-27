package config

import (
	"errors"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		File  string
		Key   string
		Error error
	}{
		{
			File:  "./../targets.ini",
			Key:   "targets",
			Error: nil,
		},
		{
			File:  "./../rbls.ini",
			Key:   "rbl",
			Error: nil,
		},
		{
			File:  "./does-not-exists.ini",
			Key:   "foo",
			Error: errors.New("Section does not exists"),
		},
		{
			File:  "./../targets.ini",
			Key:   "blah",
			Error: errors.New("Section does not exists"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.File, func(t *testing.T) {
			_, err := LoadFile(tt.File, tt.Key)
			if tt.Error == nil {
				if err != nil {
					t.Errorf("Could not load '%s' with key '%s': %s", tt.File, tt.Key, err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error for: '%s', but got none.", tt.File)
				}
			}
		})
	}
}
