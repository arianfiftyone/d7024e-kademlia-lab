package kademlia

import (
	"testing"
	"time"
)

func TestJoin(t *testing.T) {
	kademliaBootsrap := NewKademlia("127.0.0.1", 3001, true, "nil", 0)

	go kademliaBootsrap.Start()

	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 3002, false, "127.0.0.1", 3001)

	kademlia.Join()

	// assert.Fail(t, "Join failed!")
}
