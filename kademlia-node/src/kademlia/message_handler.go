package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type MessageHandler struct {
	kademliaNode *KademliaNode
}

func (messageHandler *MessageHandler) HandleMessage(rawMessage []byte) []byte {
	var message Message
	err := json.Unmarshal(rawMessage, &message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("MessageType: %s \n", message.MessageType)

	switch message.MessageType {

	case PING:
		var ping Ping

		json.Unmarshal(rawMessage, &ping)

		fmt.Println(ping.FromAddress + " sent you a ping")

		hostname, _ := os.Hostname()
		ips, _ := net.LookupIP(hostname)
		myIp := ips[0]
		ack := NewAckPingMessage(myIp.String())
		bytes, err := json.Marshal(ack)
		if err != nil {
			fmt.Println(err)
		}

		return bytes

	case FIND_NODE:
		var findN FindNode

		json.Unmarshal(rawMessage, &findN)

		fmt.Println(findN.FromAddress + " wants to find your k closest nodes.")
		closestKNodesList := messageHandler.kademliaNode.RoutingTable.FindClosestContacts(findN.ID, NumberOfClosestNodesToRetrieved)
		bytes, err := json.Marshal(closestKNodesList)
		if err != nil {
			fmt.Println(err)
		}

		return bytes

	default:
		errorMessage := NewErrorMessage()
		bytes, err := json.Marshal(errorMessage)
		if err != nil {
			fmt.Println(err)
		}
		return bytes
	}
}
