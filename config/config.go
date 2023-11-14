package config

import (
	"errors"

	"golang.org/x/exp/slog"
	"gopkg.in/ini.v1"
)

var ErrNoSuchSection = errors.New("section does not exist")
var ErrNoServerEntries = errors.New("please add a few server= entries to your .ini")

type Config struct {
	Logger *slog.Logger
}

// ValidateConfig validate the supplied configuration, e.g. check if we have "server="
func (c *Config) ValidateConfig(cfg *ini.File, section string) error {
	configSection, err := cfg.GetSection(section)
	if err != nil {
		return ErrNoSuchSection
	}

	if !configSection.HasKey("server") {
		return ErrNoServerEntries
	}

	return nil
}

func loadConfig(path string) (*ini.File, error) {
	cfg, err := ini.ShadowLoad(path)
	if err != nil {
		return nil, err
	}

	return cfg, err
}

// GetRbls returns all rbls from the config
func (c *Config) GetRbls(cfg *ini.File) []string {
	return cfg.Section("rbl").Key("server").ValueWithShadows()
}

// GetTargets returns all targets from the config
func (c *Config) GetTargets(cfg *ini.File) []string {
	return cfg.Section("targets").Key("server").ValueWithShadows()
}

// LoadFile ...
func (c *Config) LoadFile(path string) (*ini.File, error) {
	c.Logger.Debug("loading configuration file: " + path)

	cfg, err := loadConfig(path)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
