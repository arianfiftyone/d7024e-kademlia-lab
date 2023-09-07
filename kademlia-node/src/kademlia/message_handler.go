package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type MessageHandler interface {
	HandleMessage(rawMessage []byte) []byte
}

type MessageHandlerImplementation struct {
	kademliaNode *KademliaNode
}

func (messageHandler *MessageHandlerImplementation) HandleMessage(rawMessage []byte) []byte {
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

	case FIND_DATA:
		var findData FindData

		json.Unmarshal(rawMessage, &findData)

		fmt.Println(findData.FromAddress + " wants to find a value.")

		data, err := messageHandler.kademliaNode.DataStore.Get(findData.Key)
		if err != nil {
			closestKNodesList := messageHandler.kademliaNode.RoutingTable.FindClosestContacts(findData.ID, NumberOfClosestNodesToRetrieved)
			bytes, err2 := json.Marshal(closestKNodesList)
			if err2 != nil {
				fmt.Print(err2)
			}
			return bytes

		} else {
			bytes, err3 := json.Marshal(data)
			if err3 != nil {
				fmt.Print(err3)
			}

			return bytes
		}

	default:
		errorMessage := NewErrorMessage()
		bytes, err := json.Marshal(errorMessage)
		if err != nil {
			fmt.Println(err)
		}
		return bytes
	}
}
