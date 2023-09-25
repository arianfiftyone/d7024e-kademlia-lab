package kademlia

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost", 8000))

	rt.AddContact(NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost", 8001))
	rt.AddContact(NewContact(GenerateNewKademliaID("1111111100000000000000000000000000000000"), "localhost", 8002))
	rt.AddContact(NewContact(GenerateNewKademliaID("1111111200000000000000000000000000000000"), "localhost", 8002))
	rt.AddContact(NewContact(GenerateNewKademliaID("1111111300000000000000000000000000000000"), "localhost", 8002))
	rt.AddContact(NewContact(GenerateNewKademliaID("1111111400000000000000000000000000000000"), "localhost", 8002))
	rt.AddContact(NewContact(GenerateNewKademliaID("2111111400000000000000000000000000000000"), "localhost", 8002))

	contacts := rt.FindClosestContacts(GenerateNewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}
