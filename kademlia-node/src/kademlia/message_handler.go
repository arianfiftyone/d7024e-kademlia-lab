package kademlia

import (
	"encoding/json"
	"strconv"

	"github.com/arianfiftyone/src/logger"
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
		logger.Log("Error when unmarshaling `message` message: " + err.Error())
		return nil, err
	}
	logger.Log("MessageType: " + string(message.MessageType))

	if err := message.MessageType.IsValid(); err != nil {
		return nil, err
	} else {
		messageHandler.kademliaNode.updateRoutingTable(message.From)
	}

	switch message.MessageType {

	case PING:
		var ping Ping

		json.Unmarshal(rawMessage, &ping)

		logger.Log(ping.From.Ip + " sent you a ping")

		pong := NewPongMessage(messageHandler.kademliaNode.RoutingTable.Me)
		bytes, err := json.Marshal(pong)
		if err != nil {
			logger.Log("Error when unmarshaling `pong` message: " + err.Error())
			return nil, err
		}

		return bytes, nil

	case FIND_NODE:
		var findN FindNode

		json.Unmarshal(rawMessage, &findN)

		logger.Log(findN.From.Ip + " wants to find your k closest nodes.")
		closestKNodesList := messageHandler.kademliaNode.RoutingTable.FindClosestContacts(findN.ID, NumberOfClosestNodesToRetrieved)

		bytes, err := json.Marshal(NewFoundContactsMessage(messageHandler.kademliaNode.RoutingTable.Me, closestKNodesList))
		if err != nil {
			logger.Log("Error when marshaling `closetsKNodesList`: " + err.Error())
			return nil, err
		}

		return bytes, nil

	case FIND_DATA:
		var findData FindData

		json.Unmarshal(rawMessage, &findData)

		logger.Log(findData.FromAddress + " wants to find a value.")

		data, err := messageHandler.kademliaNode.DataStore.Get(findData.Key)
		if err != nil {
			closestKNodesList := messageHandler.kademliaNode.RoutingTable.FindClosestContacts(findData.ID, NumberOfClosestNodesToRetrieved)
			bytes, err := json.Marshal(NewFoundDataMessage(messageHandler.kademliaNode.RoutingTable.Me, closestKNodesList, ""))
			if err != nil {
				logger.Log("Error when marshaling `closetsKNodesList`: " + err.Error())
				return nil, err
			}
			return bytes, nil

		} else {
			bytes, err := json.Marshal(NewFoundDataMessage(messageHandler.kademliaNode.RoutingTable.Me, nil, data))
			if err != nil {
				logger.Log("Error when marshaling `data`: " + err.Error())
				return nil, err
			}

			return bytes, nil
		}

	case STORE:
		var store Store

		json.Unmarshal(rawMessage, &store)

		messageHandler.kademliaNode.DataStore.Insert(store.Key, store.Value)

		logger.Log(store.FromAddress + " wants to to store an object at the K(=" + strconv.Itoa(NumberOfClosestNodesToRetrieved) + ") nodes nearest to the hash of the data object in question")

		newStoreResponse := NewStoreResponseMessage(messageHandler.kademliaNode.RoutingTable.Me)
		bytes, err := json.Marshal(newStoreResponse)
		if err != nil {
			logger.Log("Error when marshaling `newStoreResponse`: " + err.Error())
			return nil, err
		}

		return bytes, nil

	default:
		errorMessage := NewErrorMessage(messageHandler.kademliaNode.RoutingTable.Me)
		bytes, err := json.Marshal(errorMessage)
		if err != nil {
			logger.Log("Error when marshaling `errorMessage`: " + err.Error())
			return nil, err
		}
		return bytes, nil
	}
}
