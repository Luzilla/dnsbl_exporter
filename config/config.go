package config

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// ValidateConfig validate the supplied configuration, e.g. check if we have "server="
func ValidateConfig(cfg *ini.File, section string) error {
	configSection, err := cfg.GetSection(section)
	if err != nil {
		return errors.New("Section does not exists")
	}

	if !configSection.HasKey("server") {
		return errors.New("Please add a few server= entries to your rbls.ini")
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
func GetRbls(cfg *ini.File) []string {
	return cfg.Section("rbl").Key("server").ValueWithShadows()
}

// GetTargets returns all targets from the config
func GetTargets(cfg *ini.File) []string {
	return cfg.Section("targets").Key("server").ValueWithShadows()
}

// LoadFile ...
func LoadFile(path string, section string) (*ini.File, error) {
	log.Debugln("Loading configuration...", path)

	cfg, err := loadConfig(path)
	if err != nil {
		return nil, err
	}

	err = ValidateConfig(cfg, section)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
