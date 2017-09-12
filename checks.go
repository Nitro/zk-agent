package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/upfluence/sensu-client-go/sensu"
	"github.com/upfluence/sensu-client-go/sensu/check"
	"github.com/upfluence/sensu-client-go/sensu/handler"
	"github.com/upfluence/sensu-go/sensu/transport/rabbitmq"
)

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
	}
	log.Warnln("Check the ZK leader.  Not enough synced followers", zkLeader)
	return false
}

func numLeadersFollowers(zkAddrs []string) (int, int) {
	numLeaders := 0
	numFollowers := 0
	for _, node := range zkAddrs {
		if cluster[node].leader {
			numLeaders += 1
		}
		if cluster[node].follower {
			numFollowers += 1
		}
	}
	return numLeaders, numFollowers
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

func runChecks() {
	cluster = make(map[string]*zkNode)
	initCluster(config.ZkAddresses)
	gatherZkMetrics(config.ZkAddresses)
	leaderOk, zkLeader := findLeader(config.ZkAddresses)
	if leaderOk {
		log.Debugln("Leader is:", zkLeader)
	}
	nodesOk := checkNodes(config.ZkAddresses)
	clusterOk := clusterInSync(config.ZkAddresses, zkLeader)
	numLeaders, numFollowers := numLeadersFollowers(config.ZkAddresses)
	log.Infoln("Number of Leaders:", numLeaders, "Number of Followers:", numFollowers, "Size of Cluster:", len(config.ZkAddresses))
	if nodesOk && clusterOk && leaderOk && numLeaders == 1 && (numLeaders+numFollowers == len(config.ZkAddresses)) {
		log.Infoln("The ZK Cluster is in healthy state")
	} else {
		log.Warnln("The ZK Cluster is not healthy")
	}
}

func SensuCheck() check.ExtensionCheckResult {
	cluster = make(map[string]*zkNode)
	initCluster(config.ZkAddresses)

	gatherZkMetrics(config.ZkAddresses)
	leaderOk, zkLeader := findLeader(config.ZkAddresses)
	nodesOk := checkNodes(config.ZkAddresses)
	clusterOk := clusterInSync(config.ZkAddresses, zkLeader)
	numLeaders, numFollowers := numLeadersFollowers(config.ZkAddresses)

	if !nodesOk || !clusterOk || !leaderOk || numLeaders != 1 || (numLeaders+numFollowers != len(config.ZkAddresses)) {
		return handler.Error(fmt.Sprintf("ZK CLuster is not healthy.  number of leaders: %s, number of followers: %s, size of cluster: %s",
			numLeaders, numFollowers, len(config.ZkAddresses)))
	}
	return handler.Ok("The ZK Cluster is in healthy state")
}

func runSensu() {
	sensuConfig, err := sensu.NewConfigFromFile(sensu.ExtractFlags(), *opts.SensuConfigFile)

	if err != nil {
		log.Fatalf("Error reading Sensu config: %s", err)
	}

	log.Printf("URL: %s", sensuConfig.RabbitMQURI())

	transport, err := rabbitmq.NewRabbitMQTransport(sensuConfig.RabbitMQURI())
	if err != nil {
		log.Fatalf("Error reading Sensu config: %s", err)
	}

	sensuClient := sensu.NewClient(transport, sensuConfig)

	check.Store["zookeeper_check"] = &check.ExtensionCheck{SensuCheck}

	sensuClient.Start()
}
