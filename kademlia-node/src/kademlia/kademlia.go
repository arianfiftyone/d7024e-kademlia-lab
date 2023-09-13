package kademlia

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

type Kademlia struct {
	Network          Network
	KademliaNode     *KademliaNode
	isBootstrap      bool
	bootstrapContact *Contact
}

var mutex sync.Mutex

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

func (kademlia *Kademlia) QueryAlphaContacts(contactsToQuery []Contact, queriedContacts *[]Contact, target Contact, foundContactsChannel chan []Contact) {
	for i := 0; i < len(contactsToQuery); i++ {
		go func(contactToQuery Contact) {
			mutex.Lock()
			*queriedContacts = append(*queriedContacts, contactToQuery)
			mutex.Unlock()
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

func (kademlia *Kademlia) getContactsToQuery(queriedContacts *[]Contact, closestToTargetList *[]Contact) []Contact {
	mutex.Lock()
	contactsToQuery := []Contact{}
	currentAmountToQuery := 0
	for _, contact := range *closestToTargetList {
		if currentAmountToQuery >= NumberOfAlphaContacts {
			break
		}

		isQueried := false
		for _, queriedContact := range *queriedContacts {
			if contact.ID == queriedContact.ID {
				isQueried = true
			}
		}
		if !isQueried {
			contactsToQuery = append(contactsToQuery, contact)
			currentAmountToQuery++
		}
	}
	mutex.Unlock()
	return contactsToQuery
}
func (kademlia *Kademlia) lookupRound(target *Contact, lookupCompleteChannel chan bool, stop *bool, previousClosestToTargetList []Contact, queriedContacts *[]Contact, closestToTargetList *[]Contact) {
	contactsToQuery := kademlia.getContactsToQuery(queriedContacts, closestToTargetList)
	mutex.Lock()
	if *stop {
		mutex.Unlock()
		return
	}
	mutex.Unlock()

	foundContactsChannel := make(chan []Contact)

	kademlia.QueryAlphaContacts(contactsToQuery, queriedContacts, *target, foundContactsChannel)
	timesTimedOut := 0

	for i := 0; i < len(contactsToQuery); i++ {
		select {
		case foundContacts := <-foundContactsChannel:
			mutex.Lock()

			kClosest := kademlia.getKClosest(*closestToTargetList, foundContacts, target.ID, NumberOfClosestNodesToRetrieved)
			*closestToTargetList = kClosest

			mutex.Unlock()
			go kademlia.lookupRound(target, lookupCompleteChannel, stop, *closestToTargetList, queriedContacts, closestToTargetList)

		case <-time.After(time.Second * 4):
			timesTimedOut++
		}

	}
	mutex.Lock()
	if (len(previousClosestToTargetList) != 0 && containsAll(*closestToTargetList, previousClosestToTargetList)) || timesTimedOut >= len(contactsToQuery) {
		*stop = true
		mutex.Unlock()
		lookupCompleteChannel <- true
	} else {
		mutex.Unlock()
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	queriedContacts := new([]Contact)

	var closestToTargetList *[]Contact
	alphaClosest := kademlia.KademliaNode.RoutingTable.FindClosestContacts(target.ID, NumberOfAlphaContacts)
	closestToTargetList = &alphaClosest

	lookupCompleteChannel := make(chan bool)
	stop := false
	go kademlia.lookupRound(target, lookupCompleteChannel, &stop, []Contact{}, queriedContacts, closestToTargetList)
	lookupComplete := <-lookupCompleteChannel

	if !lookupComplete {
		return nil, errors.New("Something went wrong!")
	}
	mutex.Lock()
	kClosest := *closestToTargetList
	mutex.Unlock()
	contactsToQuery := kademlia.getContactsToQuery(queriedContacts, closestToTargetList)

	for i := 0; i < len(contactsToQuery); i++ {
		foundContacts, err := kademlia.Network.SendFindContactMessage(&kademlia.KademliaNode.RoutingTable.me, &contactsToQuery[i], target.ID)

		if err != nil {
			log.Printf("Failed find node: %v\n", err)
		} else {
			kClosest = kademlia.getKClosest(kClosest, foundContacts, target.ID, NumberOfClosestNodesToRetrieved)

		}

	}

	return kClosest, nil
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
