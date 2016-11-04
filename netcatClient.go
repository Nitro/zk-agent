package main

import (
	"net"

	log "github.com/Sirupsen/logrus"
)

// StartClient starts TCP connector
func ncClient(proto string, addr string, cmd string) {
	conn, err := net.Dial(proto, addr)
	if err != nil {
		log.Fatalln(err)
	}
	local := conn.LocalAddr()
	remote := conn.RemoteAddr()
	log.Debugln("Connected from ", local)
	log.Debugln("Connected to ", remote)

	conn.Write([]byte(cmd))
	readOutputs(conn, addr, cmd)
}
