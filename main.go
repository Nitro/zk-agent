package main

import (
	log "github.com/Sirupsen/logrus"
)

var cluster map[string]*zkNode

var config *GoavailConfig
var opts *GoavailOpts

func main() {
	log.SetLevel(log.DebugLevel)
	opts = parseCommandLine()
	config = parseConfig(*opts.ZkConfigFile)
	if *opts.Command == "run-checks" {
		runChecks()
	} else if *opts.Command == "run-sensu" {
		runSensu()
	}
}
