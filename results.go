package main

import (
	"bufio"
	"net"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type mntr struct {
	zk_avg_latency                int64
	zk_max_latency                int64
	zk_min_latency                int64
	zk_packets_received           int64
	zk_packets_sent               int64
	zk_num_alive_connections      int64
	zk_outstanding_requests       int64
	zk_server_state               string
	zk_znode_count                int64
	zk_watch_count                int64
	zk_ephemerals_count           int64
	zk_approximate_data_size      int64
	zk_open_file_descriptor_count int64
	zk_max_file_descriptor_count  int64
	zk_followers                  int
	zk_synced_followers           int
	zk_pending_syncs              int
}

type srvr struct {
	latency     string
	received    int64
	sent        int64
	connections int64
	outstanding int64
	zxid        string
	mode        string
	nodeCount   int64
}

type zkNode struct {
	zkAddr    string
	leader    bool
	follower  bool
	zkVersion string
	ruok      bool
	srvrCmd   *srvr
	mntrCmd   *mntr
}

func initCluster(zkNodes []string) {
	for _, node := range zkNodes {
		cluster[node] = &zkNode{zkAddr: node, leader: false, srvrCmd: &srvr{}, mntrCmd: &mntr{}}
	}
}

func readOutputs(conn net.Conn, addr string, cmd string) {
	log.Debugln("Reading outputs")
	var val64 int64
	var valint int
	var err error

	if conn == nil {
		cluster[addr].ruok = false
		return
	}

	message := bufio.NewScanner(conn)
	for message.Scan() {
		if cmd == "mntr" {
			output := strings.Split(message.Text(), "\t")
			log.Debugln(output)
			if output[0] == "zk_version" {
				cluster[addr].zkVersion = strings.Split(output[1], ",")[0]
				continue
			}

			if output[0] == "zk_followers" || output[0] == "zk_synced_followers" || output[0] == "zk_pending_syncs" {
				valint, err = strconv.Atoi(output[1])
				if err != nil {
					log.Fatalln("Could not parse output value", err)
				}
			}
			if output[0] != "zk_server_state" {
				val64, err = strconv.ParseInt(output[1], 10, 64)
				if err != nil {
					log.Fatalln("Could not parse output value", err)
				}
			}

			switch {
			case output[0] == "zk_avg_latency":
				cluster[addr].mntrCmd.zk_avg_latency = val64
			case output[0] == "zk_max_latency":
				cluster[addr].mntrCmd.zk_max_latency = val64
			case output[0] == "zk_min_latency":
				cluster[addr].mntrCmd.zk_min_latency = val64
			case output[0] == "zk_packets_received":
				cluster[addr].mntrCmd.zk_packets_received = val64
			case output[0] == "zk_packets_sent":
				cluster[addr].mntrCmd.zk_packets_sent = val64
			case output[0] == "zk_num_alive_connections":
				cluster[addr].mntrCmd.zk_num_alive_connections = val64
			case output[0] == "zk_outstanding_requests":
				cluster[addr].mntrCmd.zk_outstanding_requests = val64
			case output[0] == "zk_server_state":
				cluster[addr].mntrCmd.zk_server_state = output[1]
			case output[0] == "zk_znode_count":
				cluster[addr].mntrCmd.zk_znode_count = val64
			case output[0] == "zk_watch_count":
				cluster[addr].mntrCmd.zk_watch_count = val64
			case output[0] == "zk_ephemerals_count":
				cluster[addr].mntrCmd.zk_ephemerals_count = val64
			case output[0] == "zk_approximate_data_size":
				cluster[addr].mntrCmd.zk_approximate_data_size = val64
			case output[0] == "zk_open_file_descriptor_count":
				cluster[addr].mntrCmd.zk_open_file_descriptor_count = val64
			case output[0] == "zk_max_file_descriptor_count":
				cluster[addr].mntrCmd.zk_max_file_descriptor_count = val64
			case output[0] == "zk_followers":
				cluster[addr].mntrCmd.zk_followers = valint
			case output[0] == "zk_synced_followers":
				cluster[addr].mntrCmd.zk_synced_followers = valint
			case output[0] == "zk_pending_syncs":
				cluster[addr].mntrCmd.zk_pending_syncs = valint
			}

			// } else if cmd == "srvr" {
			// 	output := strings.Split(message.Text(), ":")
			// 	log.Debugln(output)
			// 	if output[0] == "Zookeeper version" {
			// 		continue
			// 	}
			// 	if output[0] != "Latency min/avg/max" && output[0] != "Mode" && output[0] != "Zxid" {
			// 		val, err = strconv.ParseInt(output[1], 10, 64)
			// 		if err != nil {
			// 			log.Fatalln("Could not parse output value", err)
			// 		}
			// 	}
			//
			// 	switch {
			// 	case output[0] == "Latency min/avg/max":
			// 		cluster[addr].srvrCmd.latency = output[1]
			// 	case output[0] == "Received":
			// 		cluster[addr].srvrCmd.received = val
			// 	case output[0] == "Sent":
			// 		cluster[addr].srvrCmd.sent = val
			// 	case output[0] == "Connections":
			// 		cluster[addr].srvrCmd.connections = val
			// 	case output[0] == "Outstanding":
			// 		cluster[addr].srvrCmd.outstanding = val
			// 	case output[0] == "Zxid":
			// 		cluster[addr].srvrCmd.zxid = output[1]
			// 	case output[0] == "Mode":
			// 		cluster[addr].srvrCmd.mode = output[1]
			// 	case output[0] == "Node Count":
			// 		cluster[addr].srvrCmd.nodeCount = val
			// 	}
			// }
			if cmd == "ruok" {
				if output[0] == "imok" {
					cluster[addr].ruok = true
				} else {
					cluster[addr].ruok = false
				}

			}
			if err := message.Err(); err != nil {
				log.Fatalln(err)
			}
			if cluster[addr].mntrCmd.zk_server_state == "leader" {
				cluster[addr].leader = true
			}
			if cluster[addr].mntrCmd.zk_server_state == "follower" {
				cluster[addr].follower = true
			}
		}
	}
}
