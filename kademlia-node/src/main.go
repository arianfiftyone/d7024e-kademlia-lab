package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/arianfiftyone/src/kademlia"
)

func main() {

	NODE_PORT_STR := os.Getenv("NODE_PORT")
	NODE_PORT, _ := strconv.Atoi(NODE_PORT_STR)

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	ips, err := net.LookupIP(hostname)
	if err != nil {
		panic(err)

	}
	ip := ips[0].String()
	kademliaInstance := kademlia.NewKademlia(ip, NODE_PORT)

	fmt.Printf("Starting node...\n")
	kademliaInstance.Start()

}
