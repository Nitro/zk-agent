package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type AgentOpts struct {
	Command         *string
	ZkConfigFile    *string
	SensuConfigFile *string
	DryRun          *bool
	Debug           *bool
}

func parseCommandLine() *AgentOpts {
	var opts AgentOpts

	kingpin.CommandLine.Help = "A Zookeeper Health Checks tool."
	opts.ZkConfigFile = kingpin.Flag("zk-config", "The Zookeeper configuration TOML file path").Short('z').Default("cluster.toml").String()
	opts.SensuConfigFile = kingpin.Flag("sensu-config", "The Sensu configuration file").Short('c').Default("sensu-config.json").String()
	kingpin.Command("run-checks", "Check the health of the cluster listed in cluster.toml")
	kingpin.Command("run-sensu", "Monitor the cluster by running the Sensu Client")

	command := kingpin.Parse()
	opts.Command = &command

	log.Debugln("Using", *opts.ZkConfigFile, "for configs")

	return &opts
}
