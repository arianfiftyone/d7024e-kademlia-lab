package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/arianfiftyone/src/kademlia"
)

func main() {

	NODE_PORT_STR := os.Getenv("NODE_PORT")
	NODE_PORT, _ := strconv.Atoi(NODE_PORT_STR)

	fmt.Printf("Starting node...\n")

	err := kademlia.Listen("", NODE_PORT)
	if err != nil {
		panic(err)

	}

}
