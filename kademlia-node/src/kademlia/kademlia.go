package kademlia

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/exp/slices"
)

type Kademlia struct {
	Network          Network
	KademliaNode     *KademliaNode
	isBootstrap      bool
	bootstrapContact *Contact
}

const (
	BootstrapKademliaID   = "FFFFFFFF00000000000000000000000000000000"
	NumberOfAlphaContacts = 3
)

func NewKademlia(ip string, port int, isBootstrap bool, bootstrapIp string, bootstrapPort int) *Kademlia {

	kademliaNode := NewKademliaNode(ip, port, isBootstrap)
	network := &NetworkImplementation{
		ip,
		port,
		&MessageHandlerImplementation{
			kademliaNode,
		},
	}
	kademliaNode.setNetwork(network)

	var contact Contact
	if !isBootstrap {
		contact = NewContact(
			NewKademliaID(BootstrapKademliaID),
			bootstrapIp,
			bootstrapPort,
		)
	}
	return &Kademlia{
		Network:          network,
		KademliaNode:     kademliaNode,
		isBootstrap:      isBootstrap,
		bootstrapContact: &contact,
	}

}

func (kademlia *Kademlia) Start() {
	if !kademlia.isBootstrap {
		go func() {

			kademlia.Join()

		}()

	}

	err := kademlia.Network.Listen()
	if err != nil {
		panic(err)

	}
}

func (kademlia *Kademlia) Join() {

	if kademlia.isBootstrap {
		fmt.Println("You are the bootstrap node!")
		return

	}

	err := kademlia.Network.SendPingMessage(&kademlia.KademliaNode.RoutingTable.me, kademlia.bootstrapContact)
	if err != nil {
		return
	}

	kademlia.KademliaNode.RoutingTable.AddContact(*kademlia.bootstrapContact)

	kademlia.LookupContact(&kademlia.KademliaNode.RoutingTable.me)

}

func (kademlia *Kademlia) QueryAlphaContacts(contactsToQuery []Contact, target Contact, foundContactsChannel chan []Contact) {
	for i := 0; i < len(contactsToQuery); i++ {
		go func(contactToQuery Contact) {
			foundContacts, err := kademlia.Network.SendFindContactMessage(&kademlia.KademliaNode.RoutingTable.me, &contactToQuery, target.ID)

			if err != nil {
				log.Printf("Failed to find node in channel: %v\n", err)
			} else {
				foundContactsChannel <- foundContacts
			}

		}(contactsToQuery[i])
	}
}

func (kademlia *Kademlia) getKClosest(firstList []Contact, secondList []Contact, target *KademliaID, count int) []Contact {
	var candidates ContactCandidates

	var allContacts []Contact
	var allIds []KademliaID
	for _, contact := range append(firstList, secondList...) {
		if !slices.Contains(allIds, *contact.ID) {
			allIds = append(allIds, *contact.ID)
			allContacts = append(allContacts, contact)

		}
	}

	for i, candidate := range allContacts {
		candidate.CalcDistance(target)
		allContacts[i] = candidate
	}
	candidates.Append(allContacts)

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)

}

func containsAll(first []Contact, second []Contact) bool {
	result := true

	var secondIds []KademliaID
	for _, contact := range second {
		secondIds = append(secondIds, *contact.ID)
	}

	for _, contact := range first {
		if !slices.Contains(secondIds, *contact.ID) {
			result = false
			break
		}
	}
	return result
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	queriedContactsMap := map[KademliaID]*Contact{}
	var closestToTargetList []Contact
	var previousClosestToTargetList []Contact

	contactsToQuery := kademlia.KademliaNode.RoutingTable.FindClosestContacts(target.ID, NumberOfAlphaContacts)

	for _, contactToQuery := range contactsToQuery {
		queriedContactsMap[*contactToQuery.ID] = &contactToQuery
	}
	foundContactsChannel := make(chan []Contact)

	kademlia.QueryAlphaContacts(contactsToQuery, *target, foundContactsChannel)

	var mutex sync.Mutex
	for foundContacts := range foundContactsChannel {
		mutex.Lock()
		if len(previousClosestToTargetList) != 0 && containsAll(closestToTargetList, previousClosestToTargetList) {
			break
		}
		mutex.Unlock()
		go func(foundContacts []Contact) {
			mutex.Lock()

			previousClosestToTargetList = closestToTargetList

			closestToTargetList = kademlia.getKClosest(closestToTargetList, foundContacts, target.ID, NumberOfClosestNodesToRetrieved)

			contactsToQuery := []Contact{}
			currentAmountToQuery := 0
			for _, contact := range closestToTargetList {
				if currentAmountToQuery >= NumberOfAlphaContacts {
					break
				}
				if queriedContactsMap[*contact.ID] == nil {
					contactsToQuery = append(contactsToQuery, contact)
					currentAmountToQuery++
				}
			}

			for _, contactToQuery := range contactsToQuery {
				queriedContactsMap[*contactToQuery.ID] = &contactToQuery
			}
			kademlia.QueryAlphaContacts(contactsToQuery, *target, foundContactsChannel)

			mutex.Unlock()
		}(foundContacts)

	}
	for i := 0; i < len(contactsToQuery); i++ {
		foundContacts, err := kademlia.Network.SendFindContactMessage(&kademlia.KademliaNode.RoutingTable.me, &contactsToQuery[i], target.ID)
		fmt.Println("")
		fmt.Println(foundContacts)

		if err != nil {
			log.Printf("Failed find node: %v\n", err)
		} else {
			closestToTargetList = kademlia.getKClosest(closestToTargetList, foundContacts, target.ID, NumberOfClosestNodesToRetrieved)

		}

	}

	return closestToTargetList, nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// A node finds k nodes to check if they are close to the hash
}

func findClosestNode(arr []Contact) Contact {
	closestNode := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i].Less(&closestNode) {
			closestNode = arr[i]
		}
	}
	return closestNode
}
