package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// ValidateConfig validate the supplied configuration, e.g. check if we have "server="
func ValidateConfig(cfg *ini.File, section string) {
	if !cfg.Section(section).HasKey("server") {
		log.Fatal("Please add a few server= entries to your rbls.ini")
	}
}

func loadConfig(path string) *ini.File {
	cfg, err := ini.ShadowLoad(path)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
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
func LoadFile(path string, key string) *ini.File {
	log.Debugln("Loading configuration...", path)
	cfg := loadConfig(path)

	ValidateConfig(cfg, key)

	return cfg
}
