Zk-Agent
=========

Tool for checking on the zookeeper clusters' states and health

Overview
---------

Zookeeper famously uses the "four letter words" along with netcat utility for tracking its state and usage.
This, however, becomes a bit cumbersome when trying to determine the overall state of a cluster.  This tool
takes as input a toml file listing the cluster members, performs some basic checks, and outputs the status.

Usage
-------------
Zk-agent can be run as a one time tool using the "run-checks" command or as a continuous sensu standalone check using the "run-sensu" command.  A Dockerfile is added to run the sensu client in a container. 

* zk-config flag is the toml configuration file listing the zookeeper cluster members.  See cluster.toml for an example
* sensu-config flag is the sensu client configuration file to be uses when the tool is run as a client sensu check.  See sensu-config.json for an example.

Example Usage
-------------
```
usage: zk-agent <command> [<flags>] [<args> ...]

./zk-agent run-checks --config-file prod_cluster.toml

./zk-agent run-sensu --config-file prod_cluster.toml --sensu-config config.json
```
Right now there are two commands - run-checks and run-sensu

More Details
-------------
Since many of the four letter words are a bit redundant in their output, the tool only uses 'mntr' and 'ruok' for now.
The following checks are included in the overall health:
* 'ruok' command
* is there exactly one leader
* are there exactly len(cluster) - 1 followers
* verify each node's avg latency is > 2
* verify cluster is in sync
* verify no outstanding requests on each node
