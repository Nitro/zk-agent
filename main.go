package main

var cluster map[string]*zkNode

var config *ZKConfig
var opts *AgentOpts

func main() {
	//log.SetLevel(log.DebugLevel)
	//opts = parseCommandLine()
	opts = parseCommandLine()
	config = parseConfig(*opts.ZkConfigFile)
	if *opts.Command == "run-checks" {
		runChecks()
	} else if *opts.Command == "run-sensu" {
		runSensu()
	}
}
