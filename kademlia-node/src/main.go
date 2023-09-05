package main

import (
	"fmt"

	"github.com/arianfiftyone/src/kademlia"
)

func main() {
	fmt.Printf("Starting node...\n")
	err := kademlia.Listen("", 3000)
	if err != nil {
		panic(err)

	}

}
