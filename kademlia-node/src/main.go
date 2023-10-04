package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arianfiftyone/src/api"
	"github.com/arianfiftyone/src/cli"
	"github.com/arianfiftyone/src/kademlia"
	"github.com/arianfiftyone/src/logger"
)

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "alive\n")
}

var output io.Writer = os.Stdout

func main() {
	BOOSTRAP_NODE_HOSTNAME := os.Getenv("BOOSTRAP_NODE_HOSTNAME")

	IS_BOOTSTRAP_STR := os.Getenv("IS_BOOTSTRAP")
	isBootstrap := strings.ToLower(IS_BOOTSTRAP_STR) == "true"

	var bootstrapPort int
	var bootstrapIp string

	if !isBootstrap {
		bootstrapIps, err := net.LookupIP(BOOSTRAP_NODE_HOSTNAME)
		if err != nil {
			panic(err)

		}
		BOOSTRAP_NODE_PORT_STR := os.Getenv("BOOSTRAP_NODE_PORT")
		bootstrapPort, err = strconv.Atoi(BOOSTRAP_NODE_PORT_STR)
		if err != nil {
			panic(err)
		}
		bootstrapIp = bootstrapIps[0].String()

	}

	NODE_PORT_STR := os.Getenv("NODE_PORT")
	port, err := strconv.Atoi(NODE_PORT_STR)
	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	ips, err := net.LookupIP(hostname)
	if err != nil {
		panic(err)

	}
	ip := ips[0].String()

	KademliaInstance := kademlia.NewKademlia(ip, port, isBootstrap, bootstrapIp, bootstrapPort)
	if isBootstrap {
		http.HandleFunc("/", health)
		go http.ListenAndServe(":80", nil)

	}

	logger.Log("Starting node...")

	go KademliaInstance.Start()
	time.Sleep(500 * time.Millisecond)

	api := api.NewKademliaAPI(KademliaInstance)
	go api.StartAPI()
	cli := cli.NewCli(KademliaInstance)
	cli.StartCli(output)

}
