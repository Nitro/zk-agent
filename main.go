package main

import (
	log "github.com/Sirupsen/logrus"
)

var cluster map[string]*zkNode

func gatherZkMetrics(zkAddrs []string) {
	for _, node := range zkAddrs {
		ncClient("tcp", node, "mntr")
		ncClient("tcp", node, "ruok")
	}
}

func checkNodes(zkAddrs []string) bool {
	for _, node := range zkAddrs {
		if cluster[node].ruok == false && cluster[node].mntrCmd.zk_avg_latency > 1 && cluster[node].mntrCmd.zk_outstanding_requests != 0 {
			log.Warnln("Node", node, "may be unhealthy. Avg Latency=",
				cluster[node].mntrCmd.zk_avg_latency, "Outstanding Requests: ", cluster[node].mntrCmd.zk_outstanding_requests)
			return false
		}
	}
	return true
}

func clusterInSync(zkAddrs []string, zkLeader string) bool { //bool for whether cluster is synced or not
	if cluster[zkLeader].mntrCmd.zk_followers == len(zkAddrs)-1 && cluster[zkLeader].
		mntrCmd.zk_synced_followers == len(zkAddrs)-1 && cluster[zkLeader].mntrCmd.zk_pending_syncs == 0 {
		return true
	} else {
		log.Warnln("Check the ZK leader", zkLeader)
		return false
	}
}

func findLeader(zkAddrs []string) (bool, string) { //return bool for healty state and the name of the leader
	var leaders []string
	for _, node := range zkAddrs {
		if cluster[node].leader {
			leaders = append(leaders, node)
		}
	}
	if len(leaders) > 1 {
		log.Warnln("Too many zk leaders!", leaders)
		return false, ""
	} else if len(leaders) == 0 {
		log.Warnln("No elected zk leaders!")
		return false, ""
	} else {
		return true, leaders[0]
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	opts := parseCommandLine()
	config := parseConfig(*opts.ConfigFile)
	if *opts.Command == "run-checks" {
		cluster = make(map[string]*zkNode)
		initCluster(config.ZkAddresses)
		gatherZkMetrics(config.ZkAddresses)
		leaderOk, zkLeader := findLeader(config.ZkAddresses)
		if leaderOk {
			log.Debugln("Leader is:", zkLeader)
		}
		nodesOk := checkNodes(config.ZkAddresses)
		clusterOk := clusterInSync(config.ZkAddresses, zkLeader)
		if nodesOk && clusterOk && leaderOk {
			log.Infoln("The ZK Cluster is in healthy state")
		} else {
			log.Warnln("The ZK Cluster is not healthy")
		}
	}
}
