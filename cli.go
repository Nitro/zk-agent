package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type GoavailOpts struct {
	Command    *string
	ConfigFile *string
	DryRun     *bool
	Debug      *bool
}

func parseCommandLine() *GoavailOpts {
	var opts GoavailOpts

	kingpin.CommandLine.Help = "A Zookeeper Health Checks tool."
	opts.ConfigFile = kingpin.Flag("config-file", "The configuration TOML file path").Short('f').Default("cluster.toml").String()
	kingpin.Command("run-checks", "Check the health of the cluster listed in cluster.toml")

	command := kingpin.Parse()
	opts.Command = &command

	log.Debugln("Using", *opts.ConfigFile, "for configs")

	return &opts
}
