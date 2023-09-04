package main

import (
	"github.com/arianfiftyone/src/kademlia"
	"fmt"
)


func main() {
	fmt.Printf("Starting node...\n")
	kademlia.Listen("0.0.0.0", 3000)

}
