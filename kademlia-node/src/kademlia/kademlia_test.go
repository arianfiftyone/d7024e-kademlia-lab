package kademlia

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJoinWithBootstrapOnly(t *testing.T) {
	kademliaBootsrap := NewKademlia("127.0.0.1", 2000, true, "", 0)

	go kademliaBootsrap.Start()

	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 2001, false, "127.0.0.1", 2000)

	kademlia.Join()

	contact := kademlia.KademliaNode.GetRoutingTable().Me
	assert.Equal(t, contact.ID, kademliaBootsrap.KademliaNode.GetRoutingTable().FindClosestContacts(contact.ID, 1)[0].ID, "The new node must be in the bootstraps routing table.")

	bootstrapContact := kademliaBootsrap.KademliaNode.GetRoutingTable().Me
	assert.Equal(t, bootstrapContact.ID, kademlia.KademliaNode.GetRoutingTable().FindClosestContacts(bootstrapContact.ID, 1)[0].ID, "The bootsrapt must be in the new nodes routing table.")

}

func CreateMockedJoinKademlia(kademliaID *KademliaID, ip string, port int, bootstrapContact *Contact) KademliaImplementation {
	routingTable := NewRoutingTable(NewContact(kademliaID, ip, port))
	dataStore := NewDataStore()
	kademliaNode := &KademliaNodeImplementation{
		RoutingTable: routingTable,
		DataStore:    &dataStore,
	}

	network := &NetworkImplementation{
		ip,
		port,
		&MessageHandlerImplementation{
			kademliaNode,
		},
	}
	kademliaNode.setNetwork(network)
	kademlia := KademliaImplementation{
		Network:          network,
		KademliaNode:     kademliaNode,
		isBootstrap:      false,
		bootstrapContact: bootstrapContact,
	}

	return kademlia
}

func TestJoinWithMultipleNodes(t *testing.T) {
	kademliaBootsrap := NewKademlia("127.0.0.1", 2002, true, "", 0)

	kademlia1 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 7001)
	kademlia2 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 7002)
	kademlia3 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "127.0.0.1", 7003)

	kademliaBootsrap.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)
	kademliaBootsrap.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)
	kademliaBootsrap.KademliaNode.GetRoutingTable().AddContact(kademlia3.KademliaNode.GetRoutingTable().Me)

	go kademliaBootsrap.Start()

	time.Sleep(time.Second)

	kademlia := CreateMockedJoinKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000000"), "127.0.0.1", 2001, &kademliaBootsrap.KademliaNode.GetRoutingTable().Me)

	kademlia.Join()

	contacts := kademlia.KademliaNode.GetRoutingTable().FindClosestContacts(kademlia.KademliaNode.GetRoutingTable().Me.ID, 20)
	fmt.Println(contacts)
	doesContainAll := kademlia.FirstSetContainsAllContactsOfSecondSet(contacts, []Contact{kademlia1.KademliaNode.GetRoutingTable().Me, kademlia2.KademliaNode.GetRoutingTable().Me})
	containsKademlia3 := kademlia.FirstSetContainsAllContactsOfSecondSet(contacts, []Contact{kademlia3.KademliaNode.GetRoutingTable().Me})

	assert.True(t, doesContainAll && containsKademlia3)

}

func TestLookupContact(t *testing.T) {

	bootstrap := NewKademlia("127.0.0.1", 6000, true, "", 0)

	contact1 := NewContact(GenerateNewKademliaID("0000000000000000000000000000000000000001"), "198.168.1.1", 3000)
	contact2 := NewContact(GenerateNewKademliaID("0000000000000000000000000000000000000002"), "198.168.1.2", 3000)
	contact3 := NewContact(GenerateNewKademliaID("000000000000000000000000000000000000000F"), "198.168.1.3", 3000)
	contact4 := NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000001"), "198.168.1.4", 3000)
	contact5 := NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000002"), "198.168.1.5", 3000)
	contact6 := NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000003"), "198.168.1.6", 3000)
	contact7 := NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000004"), "198.168.1.7", 3000)
	contact8 := NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000005"), "198.168.1.8", 3000)
	contact9 := NewContact(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000006"), "198.168.1.9", 3000)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact1)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact3)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact4)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact5)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact6)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact7)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact8)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact9)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(contact2)

	go bootstrap.Start()

	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 3001, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	list, _ := kademlia.LookupContact(contact1.ID)

	fmt.Println("closest To Target List")
	fmt.Println(list)

	doesContainAll := bootstrap.FirstSetContainsAllContactsOfSecondSet(list, []Contact{contact1, contact2, contact3})
	assert.True(t, doesContainAll)

}

func CreateMockedKademlia(kademliaID *KademliaID, ip string, port int) KademliaImplementation {
	routingTable := NewRoutingTable(NewContact(kademliaID, ip, port))
	dataStore := NewDataStore()
	kademliaNode := &KademliaNodeImplementation{
		RoutingTable: routingTable,
		DataStore:    &dataStore,
	}

	network := &NetworkImplementation{
		ip,
		port,
		&MessageHandlerImplementation{
			kademliaNode,
		},
	}
	kademliaNode.setNetwork(network)

	ketToStopRefreshMap := make(map[[KeySize]byte]chan bool)
	kademlia := KademliaImplementation{
		Network:             network,
		KademliaNode:        kademliaNode,
		isBootstrap:         true,
		keyToStopRefreshMap: ketToStopRefreshMap,
	}

	return kademlia
}
func TestLookupContact2(t *testing.T) {

	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 7000)

	kademlia1 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 7001)
	kademlia2 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 7002)
	kademlia3 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1", 7003)
	kademlia4 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000001"), "127.0.0.1", 7004)
	kademlia5 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000002"), "127.0.0.1", 7005)
	kademlia6 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000003"), "127.0.0.1", 7006)

	bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia6.KademliaNode.GetRoutingTable().Me)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia5.KademliaNode.GetRoutingTable().Me)

	kademlia6.KademliaNode.GetRoutingTable().AddContact(kademlia5.KademliaNode.GetRoutingTable().Me)
	kademlia6.KademliaNode.GetRoutingTable().AddContact(kademlia4.KademliaNode.GetRoutingTable().Me)
	kademlia6.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)

	kademlia5.KademliaNode.GetRoutingTable().AddContact(kademlia6.KademliaNode.GetRoutingTable().Me)
	kademlia5.KademliaNode.GetRoutingTable().AddContact(kademlia4.KademliaNode.GetRoutingTable().Me)
	kademlia5.KademliaNode.GetRoutingTable().AddContact(kademlia3.KademliaNode.GetRoutingTable().Me)

	kademlia4.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)
	kademlia4.KademliaNode.GetRoutingTable().AddContact(kademlia5.KademliaNode.GetRoutingTable().Me)
	kademlia4.KademliaNode.GetRoutingTable().AddContact(kademlia6.KademliaNode.GetRoutingTable().Me)

	kademlia3.KademliaNode.GetRoutingTable().AddContact(kademlia5.KademliaNode.GetRoutingTable().Me)
	kademlia3.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)
	kademlia3.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)

	kademlia2.KademliaNode.GetRoutingTable().AddContact(kademlia6.KademliaNode.GetRoutingTable().Me)
	kademlia2.KademliaNode.GetRoutingTable().AddContact(kademlia3.KademliaNode.GetRoutingTable().Me)
	kademlia2.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)

	kademlia1.KademliaNode.GetRoutingTable().AddContact(kademlia4.KademliaNode.GetRoutingTable().Me)
	kademlia1.KademliaNode.GetRoutingTable().AddContact(kademlia3.KademliaNode.GetRoutingTable().Me)
	kademlia1.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)

	go bootstrap.Start()
	go kademlia1.Start()
	go kademlia2.Start()
	go kademlia3.Start()
	go kademlia4.Start()
	go kademlia5.Start()
	go kademlia6.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4000, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	list, _ := kademlia.LookupContact(kademlia1.KademliaNode.GetRoutingTable().Me.ID)

	fmt.Println("closest To Target List")
	fmt.Println(list)

	doesContainAll := bootstrap.FirstSetContainsAllContactsOfSecondSet(list, []Contact{kademlia1.KademliaNode.GetRoutingTable().Me, kademlia2.KademliaNode.GetRoutingTable().Me, kademlia3.KademliaNode.GetRoutingTable().Me})
	assert.True(t, doesContainAll)

}

func TestLookupContact3(t *testing.T) {

	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 13000)

	kademlia1 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 13001)
	kademlia2 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 13002)
	kademlia3 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1", 13003)
	kademlia4 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000001"), "127.0.0.1", 13004)
	kademlia5 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000002"), "127.0.0.1", 13005)
	kademlia6 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000003"), "127.0.0.1", 13006)
	kademlia7 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF000000"), "127.0.0.1", 13007)
	kademlia8 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFFFFFFFFFF000000000000000000000000"), "127.0.0.1", 13008)
	kademlia9 := CreateMockedKademlia(GenerateNewKademliaID("0000000FFFFFFFFFFFFFFFF00000000000000000"), "127.0.0.1", 13009)

	kademlias := []KademliaImplementation{kademlia1, kademlia4, kademlia5, kademlia6, kademlia3, kademlia7, kademlia8, kademlia9, kademlia2}

	go bootstrap.Start()
	go kademlia1.Start()
	go kademlia2.Start()
	go kademlia3.Start()
	go kademlia4.Start()
	go kademlia5.Start()
	go kademlia6.Start()
	go kademlia7.Start()
	go kademlia8.Start()
	go kademlia9.Start()
	time.Sleep(1 * time.Second)

	for _, kademlia := range kademlias {
		bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia.KademliaNode.GetRoutingTable().Me)
		kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

		fmt.Print("Me: ")
		fmt.Println(kademlia.KademliaNode.GetRoutingTable().Me)
		list, _ := kademlia.LookupContact(kademlia.KademliaNode.GetRoutingTable().Me.ID)
		// fmt.Print("closestToTargetList: ")
		// fmt.Println(list)

		for _, contat := range list {
			kademlia.KademliaNode.GetRoutingTable().AddContact(contat)

		}

	}
	fmt.Println()
	fmt.Println("-------------------------------------")
	fmt.Println()

	kademlia := NewKademlia("127.0.0.1", 4000, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	for _, kademlia := range kademlias {
		list, _ := kademlia.LookupContact(kademlia1.KademliaNode.GetRoutingTable().Me.ID)
		doesContainAll := bootstrap.FirstSetContainsAllContactsOfSecondSet(list, []Contact{kademlia1.KademliaNode.GetRoutingTable().Me, kademlia2.KademliaNode.GetRoutingTable().Me, kademlia3.KademliaNode.GetRoutingTable().Me})
		assert.True(t, doesContainAll)

		if !doesContainAll {
			fmt.Print("Me: ")
			fmt.Println(kademlia.KademliaNode.GetRoutingTable().Me)
			fmt.Print("failed: ")
			fmt.Println(list)
		}
	}

}

func TestLookupDataFindsData(t *testing.T) {

	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 10069)

	kademlia1 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 10021)
	kademlia2 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 10022)

	bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)
	kademlia1.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)

	value := "value"
	key := GetKeyRepresentationOfKademliaId(GenerateNewKademliaID("0000000000000000000000000000000000000002")) // Sets the key to be the same as kademlia2's id

	kademlia2.KademliaNode.GetDataStore().Insert(key, value)

	go bootstrap.Start()
	go kademlia1.Start()
	go kademlia2.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4020, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	_, data, _ := kademlia.LookupData(key)

	assert.Equal(t, value, data)

}

func TestLookupDataFindsNoData(t *testing.T) {

	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 10031)

	kademlia1 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 10032)

	bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)

	value := "value"
	key := NewKey(value)

	go bootstrap.Start()
	go kademlia1.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4030, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	list, _, _ := kademlia.LookupData(key)
	fmt.Println(list)

	doesContainAll := bootstrap.FirstSetContainsAllContactsOfSecondSet(list, []Contact{kademlia1.KademliaNode.GetRoutingTable().Me, bootstrap.KademliaNode.GetRoutingTable().Me, kademlia.KademliaNode.GetRoutingTable().Me})
	assert.True(t, doesContainAll)
}

func TestStore(t *testing.T) {
	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 11000)

	kademlia1 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 11001)
	kademlia2 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 11002)
	kademlia3 := CreateMockedKademlia(GenerateNewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1", 11003)
	kademlia4 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000001"), "127.0.0.1", 11004)
	kademlia5 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000002"), "127.0.0.1", 11005)
	kademlia6 := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000003"), "127.0.0.1", 11006)

	bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia6.KademliaNode.GetRoutingTable().Me)
	bootstrap.KademliaNode.GetRoutingTable().AddContact(kademlia5.KademliaNode.GetRoutingTable().Me)

	kademlia6.KademliaNode.GetRoutingTable().AddContact(kademlia5.KademliaNode.GetRoutingTable().Me)
	kademlia6.KademliaNode.GetRoutingTable().AddContact(kademlia4.KademliaNode.GetRoutingTable().Me)
	kademlia6.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)

	kademlia5.KademliaNode.GetRoutingTable().AddContact(kademlia6.KademliaNode.GetRoutingTable().Me)
	kademlia5.KademliaNode.GetRoutingTable().AddContact(kademlia4.KademliaNode.GetRoutingTable().Me)
	kademlia5.KademliaNode.GetRoutingTable().AddContact(kademlia3.KademliaNode.GetRoutingTable().Me)

	kademlia4.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)
	kademlia4.KademliaNode.GetRoutingTable().AddContact(kademlia5.KademliaNode.GetRoutingTable().Me)
	kademlia4.KademliaNode.GetRoutingTable().AddContact(kademlia6.KademliaNode.GetRoutingTable().Me)

	kademlia3.KademliaNode.GetRoutingTable().AddContact(kademlia5.KademliaNode.GetRoutingTable().Me)
	kademlia3.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)
	kademlia3.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)

	kademlia2.KademliaNode.GetRoutingTable().AddContact(kademlia6.KademliaNode.GetRoutingTable().Me)
	kademlia2.KademliaNode.GetRoutingTable().AddContact(kademlia3.KademliaNode.GetRoutingTable().Me)
	kademlia2.KademliaNode.GetRoutingTable().AddContact(kademlia1.KademliaNode.GetRoutingTable().Me)

	kademlia1.KademliaNode.GetRoutingTable().AddContact(kademlia4.KademliaNode.GetRoutingTable().Me)
	kademlia1.KademliaNode.GetRoutingTable().AddContact(kademlia3.KademliaNode.GetRoutingTable().Me)
	kademlia1.KademliaNode.GetRoutingTable().AddContact(kademlia2.KademliaNode.GetRoutingTable().Me)

	mockedKademlias := []KademliaImplementation{
		kademlia1,
		kademlia2,
		kademlia3,
		kademlia4,
		kademlia5,
		kademlia6,
	}

	go bootstrap.Start()
	go kademlia1.Start()
	go kademlia2.Start()
	go kademlia3.Start()
	go kademlia4.Start()
	go kademlia5.Start()
	go kademlia6.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4000, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	content := "testy"
	key, err := kademlia.Store(content)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	list, _ := kademlia.LookupContact(key.GetKademliaIdRepresentationOfKey())

	fmt.Println("closest To Target List")
	fmt.Println(list)

	for _, contact := range list {
		contactInMockedKademlias := false
		for _, kademlia := range mockedKademlias {
			if kademlia.KademliaNode.GetRoutingTable().Me.ID == contact.ID {
				contactInMockedKademlias = true
				retrivedContent, err := kademlia.KademliaNode.GetDataStore().Get(key)

				if err != nil {
					assert.Fail(t, err.Error())
				}

				assert.True(t, retrivedContent == content)
				break

			}
		}
		if contactInMockedKademlias {
			_, err := kademlia.KademliaNode.GetDataStore().Get(key)

			if err == nil {
				assert.Fail(t, "nodes not found by lookup, should not get a store rpc")
			}
		}
	}

}

func TestLookupAfterJoin(t *testing.T) {
	bootstrap := NewKademlia("localhost", 1100, true, "", 0)
	go bootstrap.Start()
	time.Sleep(time.Second)

	var kademlias []*KademliaImplementation
	var allContacts []Contact
	for i := 0; i < 3; i++ {
		port := 1101 + i
		kademlia := NewKademlia("localhost", port, false, "localhost", 1100)
		allContacts = append(allContacts, kademlia.KademliaNode.GetRoutingTable().Me)
		go kademlia.Start()
		time.Sleep(time.Microsecond * 50)

		kademlias = append(kademlias, kademlia)
	}
	time.Sleep(time.Second)

	target := allContacts[len(allContacts)-1].ID

	allContacts = append(allContacts, bootstrap.KademliaNode.GetRoutingTable().Me)
	var candidates ContactCandidates
	for i, candidate := range allContacts {
		candidate.CalcDistance(target)
		allContacts[i] = candidate
	}
	candidates.Append(allContacts)

	candidates.Sort()

	count := NumberOfClosestNodesToRetrieved
	if count > candidates.Len() {
		count = candidates.Len()
	}

	expectedClosest := candidates.GetContacts(count)
	fmt.Print("Expected: ")
	fmt.Println(expectedClosest)

	for _, kademlia := range kademlias {
		contacts, err := kademlia.LookupContact(target)

		if err != nil {
			assert.Fail(t, err.Error())
		}

		containsAll := bootstrap.FirstSetContainsAllContactsOfSecondSet(contacts, expectedClosest)
		assert.True(t, containsAll)

		if !containsAll {
			fmt.Print("Failed: ")
			fmt.Println(contacts)
		}

	}

}
func TestBig(t *testing.T) {
	bootstrap := NewKademlia("localhost", 60000, true, "", 0)
	go bootstrap.Start()
	time.Sleep(time.Second)

	var kademlias []*KademliaImplementation
	for i := 0; i < 10; i++ {
		port := 60001 + i
		kademlia := NewKademlia("localhost", port, false, "localhost", 60000)
		go kademlia.Start()
		time.Sleep(time.Microsecond * 100)

		kademlias = append(kademlias, kademlia)
	}
	time.Sleep(time.Second)

	content := "hello"
	key, err := kademlias[len(kademlias)-1].Store(content)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	for _, kademlia := range kademlias {
		contacts, retrivedContent, err := kademlia.LookupData(key)

		if contacts != nil {
			closestToTargetList, _ := kademlia.LookupContact(key.GetKademliaIdRepresentationOfKey())

			fmt.Print("Me: ")
			fmt.Println(kademlia.KademliaNode.GetRoutingTable().Me)
			fmt.Print("contacts: ")
			fmt.Println(contacts)
			fmt.Print("closestToTargetList: ")
			fmt.Println(closestToTargetList)
		}

		if err != nil {
			assert.Fail(t, err.Error())
		}

		assert.Equal(t, content, retrivedContent)

	}

}

func TestRefresh(t *testing.T) {
	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 31000)

	go bootstrap.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4000, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	content := "testy"
	key, err := kademlia.Store(content)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	initTime, err := bootstrap.KademliaNode.GetDataStore().GetTime(key)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	time.Sleep((bootstrap.KademliaNode.GetDataStore().ttl / 2) + time.Millisecond*100)

	endTime, err := bootstrap.KademliaNode.GetDataStore().GetTime(key)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	fmt.Println(initTime)
	fmt.Println(endTime)
	assert.Greater(t, endTime, initTime.Add(bootstrap.KademliaNode.GetDataStore().ttl/2))
}

func TestStopRefresh(t *testing.T) {
	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 31001)

	go bootstrap.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4000, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	content := "testy"
	key, err := kademlia.Store(content)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	initTime, err := bootstrap.KademliaNode.GetDataStore().GetTime(key)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	kademlia.keyToStopRefreshMap[key.Hash] <- true

	time.Sleep((bootstrap.KademliaNode.GetDataStore().ttl / 2) + time.Millisecond*100)

	endTime, err := bootstrap.KademliaNode.GetDataStore().GetTime(key)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	fmt.Println(initTime)
	fmt.Println(endTime)

	assert.Equal(t, initTime, endTime)
}

func TestForget(t *testing.T) {
	bootstrap := CreateMockedKademlia(GenerateNewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 31011)

	go bootstrap.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4001, false, "", 0)

	kademlia.KademliaNode.GetRoutingTable().AddContact(bootstrap.KademliaNode.GetRoutingTable().Me)

	content := "testy"
	key, err := kademlia.Store(content)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	err = kademlia.Forget(key)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	time.Sleep((bootstrap.KademliaNode.GetDataStore().ttl / 2) + time.Millisecond*100)

	expectedMap := map[[KeySize]byte]string{}

	assert.Equal(t, expectedMap, kademlia.KademliaNode.GetDataStore().data)
}
