package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

type GoavailConfig struct {
	ZkAddresses []string `toml:"zk_addresses"`
}

func parseConfig(path string) *GoavailConfig {
	var config GoavailConfig

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err.Error())
	}
	return &config
}
