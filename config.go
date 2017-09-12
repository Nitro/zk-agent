package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

type ZKConfig struct {
	ZkAddresses []string `toml:"zk_addresses"`
}

func parseConfig(path string) *ZKConfig {
	var config ZKConfig

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err.Error())
	}
	return &config
}
