package kademlia

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type MessageHandler interface {
	HandleMessage(rawMessage []byte) ([]byte, error)
}

type MessageHandlerImplementation struct {
	kademliaNode *KademliaNode
}

func (messageHandler *MessageHandlerImplementation) HandleMessage(rawMessage []byte) ([]byte, error) {
	var message Message

	err := json.Unmarshal(rawMessage, &message)
	if err != nil {
		log.Printf("Error when unmarshaling `message` message: %v\n", err)
		return nil, err
	}
	fmt.Printf("MessageType: %s \n", message.MessageType)

	if err := message.MessageType.IsValid(); err != nil {
		return nil, err
	} else {
		messageHandler.kademliaNode.updateRoutingTable(message.Contact)
	}

	switch message.MessageType {

	case PING:
		var ping Ping

		json.Unmarshal(rawMessage, &ping)

		fmt.Println(ping.Contact.Ip + " sent you a ping")

		pong := NewPongMessage(messageHandler.kademliaNode.RoutingTable.me)
		bytes, err := json.Marshal(pong)
		if err != nil {
			log.Printf("Error when marshaling `pong` message: %v\n", err)
			return nil, err
		}

		return bytes, nil

	case FIND_NODE:
		var findN FindNode

		json.Unmarshal(rawMessage, &findN)

		fmt.Println(findN.FromAddress + " wants to find your k closest nodes.")
		closestKNodesList := messageHandler.kademliaNode.RoutingTable.FindClosestContacts(findN.ID, NumberOfClosestNodesToRetrieved)

		bytes, err := json.Marshal(closestKNodesList)
		if err != nil {
			log.Printf("Error when marshaling `closetsKNodesList`: %v\n", err)
			return nil, err
		}

		return bytes, nil

	case FIND_DATA:
		var findData FindData

		json.Unmarshal(rawMessage, &findData)

		fmt.Println(findData.FromAddress + " wants to find a value.")

		data, err := messageHandler.kademliaNode.DataStore.Get(findData.Key)
		if err != nil {
			closestKNodesList := messageHandler.kademliaNode.RoutingTable.FindClosestContacts(findData.ID, NumberOfClosestNodesToRetrieved)
			bytes, err := json.Marshal(closestKNodesList)
			if err != nil {
				log.Printf("Error when marshaling `closetsKNodesList`: %v\n", err)
				return nil, err
			}
			return bytes, nil

		} else {
			bytes, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error when marshaling `data`: %v\n", err)
				return nil, err
			}

			return bytes, nil
		}

	case STORE:
		var store Store

		json.Unmarshal(rawMessage, &store)

		messageHandler.kademliaNode.DataStore.Insert(store.Key, store.Value)

		fmt.Println(store.FromAddress + " wants to to store an object at the K(=" + strconv.Itoa(NumberOfClosestNodesToRetrieved) + ") nodes nearest to the hash of the data object in question")

		newStoreResponse := NewStoreResponseMessage(messageHandler.kademliaNode.RoutingTable.me)
		bytes, err := json.Marshal(newStoreResponse)
		if err != nil {
			log.Printf("Error when marshaling `newStoreResponse`: %v\n", err)
			return nil, err
		}

		return bytes, nil

	default:
		errorMessage := NewErrorMessage(messageHandler.kademliaNode.RoutingTable.me)
		bytes, err := json.Marshal(errorMessage)
		if err != nil {
			log.Printf("Error when marshaling `errorMessage`: %v\n", err)
			return nil, err
		}
		return bytes, nil
	}
}
