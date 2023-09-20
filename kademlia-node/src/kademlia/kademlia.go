package kademlia

import (
	"errors"
	"fmt"
	"sync"

	"github.com/arianfiftyone/src/logger"
	"golang.org/x/exp/slices"
)

type Kademlia struct {
	Network          Network
	KademliaNode     *KademliaNode
	isBootstrap      bool
	bootstrapContact *Contact
}

type LookupType string

var mutex sync.Mutex

const (
	BootstrapKademliaID   = "FFFFFFFF00000000000000000000000000000000"
	NumberOfAlphaContacts = 3

	LOOKUP_CONTACT LookupType = "LOOKUP_CONTACT"
	LOOKUP_DATA    LookupType = "LOOKUP_DATA"
)

// NewKademlia gives new instance of a kademlia participant, it can start lisining for RPC's and join the network.
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

func (kademlia *Kademlia) refresh() {
	var lowerBound *KademliaID
	var highBound *KademliaID

	for i := 0; i < 1; i++ {
		if kademlia.KademliaNode.RoutingTable.Me.ID.Less(kademlia.bootstrapContact.ID) {
			lowerBound = kademlia.bootstrapContact.ID
			highBound = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
		} else {
			lowerBound = NewKademliaID("0000000000000000000000000000000000000000")
			highBound = kademlia.bootstrapContact.ID
		}

		randomKademliaIDInRnge, err := NewRandomKademliaIDInRange(lowerBound, highBound)
		if err != nil {
			return
		}
		contacts, err := kademlia.LookupContact(randomKademliaIDInRnge)
		if err != nil {
			return
		}
		for _, contact := range contacts {
			kademlia.KademliaNode.RoutingTable.AddContact(contact)
		}

	}

	for i := 0; i < 0; i++ {
		if kademlia.KademliaNode.RoutingTable.Me.ID.Less(kademlia.bootstrapContact.ID) {
			lowerBound = NewKademliaID("0000000000000000000000000000000000000000")
			highBound = kademlia.KademliaNode.RoutingTable.Me.ID
		} else {
			lowerBound = kademlia.KademliaNode.RoutingTable.Me.ID
			highBound = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
		}

		randomKademliaIDInRnge, err := NewRandomKademliaIDInRange(lowerBound, highBound)
		if err != nil {
			return
		}
		contacts, err := kademlia.LookupContact(randomKademliaIDInRnge)
		if err != nil {
			return
		}
		for _, contact := range contacts {
			kademlia.KademliaNode.RoutingTable.AddContact(contact)
		}

	}
}

func (kademlia *Kademlia) Join() {

	if kademlia.isBootstrap {
		logger.Log("You are the bootstrap node!")
		return

	}

	err := kademlia.Network.SendPingMessage(&kademlia.KademliaNode.RoutingTable.Me, kademlia.bootstrapContact)
	if err != nil {
		return
	}

	kademlia.KademliaNode.RoutingTable.AddContact(*kademlia.bootstrapContact)

	contacts, err := kademlia.LookupContact(kademlia.KademliaNode.RoutingTable.Me.ID)
	if err != nil {
		return
	}
	for _, contact := range contacts {
		kademlia.KademliaNode.RoutingTable.AddContact(contact)
	}

	kademlia.refresh()
}

func (kademlia *Kademlia) Store(content string) (*Key, error) {
	// A node finds k nodes to check if they are close to the hash

	key := HashToKey(content)
	contacts, err := kademlia.LookupContact(key.GetKademliaIdRepresentationOfKey())

	if err != nil {
		return nil, err
	}

	if len(contacts) <= 0 {
		return nil, errors.New("found no node to store the value in")
	}
	for _, contact := range contacts {
		kademlia.Network.SendStoreMessage(&kademlia.KademliaNode.RoutingTable.Me, &contact, key, content)
	}
	return key, nil
}

func (kademlia *Kademlia) QueryAlphaContacts(lookupType LookupType, contactsToQuery []Contact, queriedContacts *[]Contact, targetId KademliaID, foundContactsChannel chan []Contact, foundValueChannel chan string, queryFailedChannel chan error) {

	for i := 0; i < len(contactsToQuery); i++ {
		go func(contactToQuery Contact) {
			var foundContacts []Contact
			var err error
			var foundValue string

			switch lookupType {

			case LOOKUP_CONTACT:
				foundContacts, err = kademlia.Network.SendFindContactMessage(&kademlia.KademliaNode.RoutingTable.Me, &contactToQuery, &targetId)

			case LOOKUP_DATA:
				foundContacts, foundValue, err = kademlia.Network.SendFindDataMessage(&kademlia.KademliaNode.RoutingTable.Me, &contactToQuery, GetKeyRepresentationOfKademliaId(&targetId))

			}

			mutex.Lock()
			*queriedContacts = append(*queriedContacts, contactToQuery)
			mutex.Unlock()

			if err != nil {
				queryFailedChannel <- err
				return
			}

			if foundValue != "" {
				foundValueChannel <- foundValue
				return
			}

			foundContactsChannel <- foundContacts

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

func (kademlia *Kademlia) firstSetContainsAllContactsOfSecondSet(first []Contact, second []Contact) bool {
	result := true

	var firstIds []KademliaID
	for _, contact := range first {
		firstIds = append(firstIds, *contact.ID)
	}

	for _, contact := range second {
		if !slices.Contains(firstIds, *contact.ID) {
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

func (kademlia *Kademlia) lookupRound(lookupType LookupType, targetId *KademliaID, lookupCompleteChannel chan bool, lookupDataChannel chan string, stop *bool, previousClosestToTargetList []Contact, queriedContacts *[]Contact, closestToTargetList *[]Contact) {
	contactsToQuery := kademlia.getContactsToQuery(queriedContacts, closestToTargetList)
	mutex.Lock()
	if *stop {
		mutex.Unlock()
		return
	}
	mutex.Unlock()

	foundContactsChannel := make(chan []Contact)
	foundValueChannel := make(chan string)
	queryFailedChannel := make(chan error)

	kademlia.QueryAlphaContacts(lookupType, contactsToQuery, queriedContacts, *targetId, foundContactsChannel, foundValueChannel, queryFailedChannel)
	timesFailed := 0

Loop:
	for i := 0; i < len(contactsToQuery); i++ {
		select {
		case foundContacts := <-foundContactsChannel:
			mutex.Lock()

			kClosest := kademlia.getKClosest(*closestToTargetList, foundContacts, targetId, NumberOfClosestNodesToRetrieved)
			*closestToTargetList = kClosest

			mutex.Unlock()
			go kademlia.lookupRound(lookupType, targetId, lookupCompleteChannel, lookupDataChannel, stop, *closestToTargetList, queriedContacts, closestToTargetList)
		case foundValue := <-foundValueChannel:
			*stop = true
			lookupDataChannel <- foundValue
			break Loop

		case queryFailedError := <-queryFailedChannel:
			logger.Log("Failed to find node in channel: " + queryFailedError.Error())
			fmt.Println("")
			fmt.Println("Failed to find node in channel: " + queryFailedError.Error())

			timesFailed++

		}

	}
	mutex.Lock()
	if (len(previousClosestToTargetList) != 0 && kademlia.firstSetContainsAllContactsOfSecondSet(*closestToTargetList, previousClosestToTargetList) && kademlia.firstSetContainsAllContactsOfSecondSet(previousClosestToTargetList, *closestToTargetList)) || timesFailed >= len(contactsToQuery) {
		*stop = true
		mutex.Unlock()
		lookupCompleteChannel <- true
	} else {
		mutex.Unlock()
	}
}

func (kademlia *Kademlia) lookup(lookupType LookupType, targetId *KademliaID) ([]Contact, string, error) {
	queriedContacts := new([]Contact)

	var closestToTargetList *[]Contact
	alphaClosest := kademlia.KademliaNode.RoutingTable.FindClosestContacts(targetId, NumberOfAlphaContacts)
	closestToTargetList = &alphaClosest

	lookupCompleteChannel := make(chan bool)
	lookupDataChannel := make(chan string)

	stop := false
	stopPointer := &stop
	for {
		*stopPointer = false
		go kademlia.lookupRound(lookupType, targetId, lookupCompleteChannel, lookupDataChannel, stopPointer, []Contact{}, queriedContacts, closestToTargetList)

		select {
		case <-lookupCompleteChannel:
			break

		case foundValue := <-lookupDataChannel:
			return nil, foundValue, nil
		}

		contactsToQuery := kademlia.getContactsToQuery(queriedContacts, closestToTargetList)

		if len(contactsToQuery) <= 0 {
			break
		}

		foundContactsChannel := make(chan []Contact)
		queryFailedChannel := make(chan error)

		kademlia.QueryAlphaContacts(lookupType, contactsToQuery, queriedContacts, *targetId, foundContactsChannel, nil, queryFailedChannel)
		fmt.Println("New")

		for i := 0; i < len(contactsToQuery); i++ {
			select {
			case foundContacts := <-foundContactsChannel:
				mutex.Lock()
				*closestToTargetList = kademlia.getKClosest(*closestToTargetList, foundContacts, targetId, NumberOfClosestNodesToRetrieved)
				mutex.Unlock()
			case queryFailedError := <-queryFailedChannel:
				logger.Log("Failed to find node in channel: " + queryFailedError.Error())
				fmt.Println("Failed to find node in channel: " + queryFailedError.Error())

			}

		}
		fmt.Println("Retry!")
	}
	mutex.Lock()
	kClosest := *closestToTargetList
	mutex.Unlock()

	return kClosest, "", nil

}

func (kademlia *Kademlia) LookupContact(targetId *KademliaID) ([]Contact, error) {
	kClosest, _, err := kademlia.lookup(LOOKUP_CONTACT, targetId)
	return kClosest, err
}

func (kademlia *Kademlia) LookupData(key *Key) ([]Contact, string, error) {
	kClosest, value, err := kademlia.lookup(LOOKUP_DATA, key.GetKademliaIdRepresentationOfKey())
	return kClosest, value, err

}
