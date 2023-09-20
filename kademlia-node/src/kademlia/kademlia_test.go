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

	contact := kademlia.KademliaNode.RoutingTable.Me
	assert.Equal(t, contact.ID, kademliaBootsrap.KademliaNode.RoutingTable.FindClosestContacts(contact.ID, 1)[0].ID, "The new node must be in the bootstraps routing table.")

	bootstrapContact := kademliaBootsrap.KademliaNode.RoutingTable.Me
	assert.Equal(t, bootstrapContact.ID, kademlia.KademliaNode.RoutingTable.FindClosestContacts(bootstrapContact.ID, 1)[0].ID, "The bootsrapt must be in the new nodes routing table.")

}

func CreateMockedJoinKademlia(kademliaID *KademliaID, ip string, port int, bootstrapContact *Contact) Kademlia {
	routingTable := NewRoutingTable(NewContact(kademliaID, ip, port))
	dataStore := NewDataStore()
	kademliaNode := &KademliaNode{
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
	kademlia := Kademlia{
		Network:          network,
		KademliaNode:     kademliaNode,
		isBootstrap:      false,
		bootstrapContact: bootstrapContact,
	}

	return kademlia
}

func TestJoinWithMultipleNodes(t *testing.T) {
	kademliaBootsrap := NewKademlia("127.0.0.1", 2002, true, "", 0)

	kademlia1 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 7001)
	kademlia2 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 7002)
	kademlia3 := CreateMockedKademlia(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "127.0.0.1", 7003)

	kademliaBootsrap.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)
	kademliaBootsrap.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)
	kademliaBootsrap.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.Me)

	go kademliaBootsrap.Start()

	time.Sleep(time.Second)

	kademlia := CreateMockedJoinKademlia(NewKademliaID("0000000000000000000000000000000000000000"), "127.0.0.1", 2001, &kademliaBootsrap.KademliaNode.RoutingTable.Me)

	kademlia.Join()

	contacts := kademlia.KademliaNode.RoutingTable.FindClosestContacts(kademlia.KademliaNode.RoutingTable.Me.ID, 20)
	fmt.Println(contacts)
	doesContainAll := kademlia.firstSetContainsAllContactsOfSecondSet(contacts, []Contact{kademlia1.KademliaNode.RoutingTable.Me, kademlia2.KademliaNode.RoutingTable.Me})
	containsKademlia3 := kademlia.firstSetContainsAllContactsOfSecondSet(contacts, []Contact{kademlia3.KademliaNode.RoutingTable.Me})

	assert.True(t, doesContainAll && containsKademlia3)

}

type networkMock struct{}

func (network *networkMock) Listen() error {
	return nil
}
func (network *networkMock) Send(ip string, port int, message []byte, timeOut time.Duration) ([]byte, error) {
	return nil, nil
}
func (network *networkMock) SendPingMessage(from *Contact, contact *Contact) error {
	return nil
}
func (network *networkMock) SendFindContactMessage(from *Contact, contact *Contact, id *KademliaID) ([]Contact, error) {
	return nil, nil
}
func (network *networkMock) SendFindDataMessage(from *Contact, contact *Contact, key *Key) ([]Contact, string, error) {
	return nil, "", nil
}
func (network *networkMock) SendStoreMessage(from *Contact, contact *Contact, value string) bool {
	return false
}

func TestLookupContact(t *testing.T) {

	bootstrap := NewKademlia("127.0.0.1", 6000, true, "", 0)

	contact1 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "198.168.1.1", 3000)
	contact2 := NewContact(NewKademliaID("0000000000000000000000000000000000000002"), "198.168.1.2", 3000)
	contact3 := NewContact(NewKademliaID("000000000000000000000000000000000000000F"), "198.168.1.3", 3000)
	contact4 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "198.168.1.4", 3000)
	contact5 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "198.168.1.5", 3000)
	contact6 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "198.168.1.6", 3000)
	contact7 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000004"), "198.168.1.7", 3000)
	contact8 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000005"), "198.168.1.8", 3000)
	contact9 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000006"), "198.168.1.9", 3000)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact1)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact3)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact4)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact5)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact6)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact7)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact8)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact9)
	bootstrap.KademliaNode.RoutingTable.AddContact(contact2)

	go bootstrap.Start()

	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 3001, false, "", 0)

	kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.Me)

	list, _ := kademlia.LookupContact(contact1.ID)

	fmt.Println("closest To Target List")
	fmt.Println(list)

	doesContainAll := bootstrap.firstSetContainsAllContactsOfSecondSet(list, []Contact{contact1, contact2, contact3})
	assert.True(t, doesContainAll)

}

func CreateMockedKademlia(kademliaID *KademliaID, ip string, port int) Kademlia {
	routingTable := NewRoutingTable(NewContact(kademliaID, ip, port))
	dataStore := NewDataStore()
	kademliaNode := &KademliaNode{
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
	kademlia := Kademlia{
		Network:      network,
		KademliaNode: kademliaNode,
		isBootstrap:  true,
	}

	return kademlia
}
func TestLookupContact2(t *testing.T) {

	bootstrap := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 7000)

	kademlia1 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 7001)
	kademlia2 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 7002)
	kademlia3 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1", 7003)
	kademlia4 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "127.0.0.1", 7004)
	kademlia5 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "127.0.0.1", 7005)
	kademlia6 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "127.0.0.1", 7006)

	bootstrap.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.Me)
	bootstrap.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.Me)

	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.Me)
	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.Me)
	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)

	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.Me)
	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.Me)
	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.Me)

	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)
	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.Me)
	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.Me)

	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.Me)
	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)
	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)

	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.Me)
	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.Me)
	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)

	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.Me)
	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.Me)
	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)

	go bootstrap.Start()
	go kademlia1.Start()
	go kademlia2.Start()
	go kademlia3.Start()
	go kademlia4.Start()
	go kademlia5.Start()
	go kademlia6.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4000, false, "", 0)

	kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.Me)

	list, _ := kademlia.LookupContact(kademlia1.KademliaNode.RoutingTable.Me.ID)

	fmt.Println("closest To Target List")
	fmt.Println(list)

	doesContainAll := bootstrap.firstSetContainsAllContactsOfSecondSet(list, []Contact{kademlia1.KademliaNode.RoutingTable.Me, kademlia2.KademliaNode.RoutingTable.Me, kademlia3.KademliaNode.RoutingTable.Me})
	assert.True(t, doesContainAll)

}

func TestLookupContact3(t *testing.T) {

	bootstrap := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 13000)

	kademlia1 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 13001)
	kademlia2 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 13002)
	kademlia3 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1", 13003)
	kademlia4 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "127.0.0.1", 13004)
	kademlia5 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "127.0.0.1", 13005)
	kademlia6 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "127.0.0.1", 13006)
	kademlia7 := CreateMockedKademlia(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF000000"), "127.0.0.1", 13007)
	kademlia8 := CreateMockedKademlia(NewKademliaID("FFFFFFFFFFFFFFFF000000000000000000000000"), "127.0.0.1", 13008)
	kademlia9 := CreateMockedKademlia(NewKademliaID("0000000FFFFFFFFFFFFFFFF00000000000000000"), "127.0.0.1", 13009)

	kademlias := []Kademlia{kademlia1, kademlia4, kademlia5, kademlia6, kademlia3, kademlia7, kademlia8, kademlia9, kademlia2}

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
		bootstrap.KademliaNode.RoutingTable.AddContact(kademlia.KademliaNode.RoutingTable.Me)
		kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.Me)

		fmt.Print("Me: ")
		fmt.Println(kademlia.KademliaNode.RoutingTable.Me)
		list, _ := kademlia.LookupContact(kademlia.KademliaNode.RoutingTable.Me.ID)
		// fmt.Print("closestToTargetList: ")
		// fmt.Println(list)

		for _, contat := range list {
			kademlia.KademliaNode.RoutingTable.AddContact(contat)

		}

	}
	fmt.Println()
	fmt.Println("-------------------------------------")
	fmt.Println()

	kademlia := NewKademlia("127.0.0.1", 4000, false, "", 0)

	kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.Me)

	for _, kademlia := range kademlias {
		list, _ := kademlia.LookupContact(kademlia1.KademliaNode.RoutingTable.Me.ID)
		doesContainAll := bootstrap.firstSetContainsAllContactsOfSecondSet(list, []Contact{kademlia1.KademliaNode.RoutingTable.Me, kademlia2.KademliaNode.RoutingTable.Me, kademlia3.KademliaNode.RoutingTable.Me})
		assert.True(t, doesContainAll)

		if !doesContainAll {
			fmt.Print("Me: ")
			fmt.Println(kademlia.KademliaNode.RoutingTable.Me)
			fmt.Print("failed: ")
			fmt.Println(list)
		}
	}

}

func TestLookupDataFindsData(t *testing.T) {

	bootstrap := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 10020)

	kademlia1 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 10021)
	kademlia2 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 10022)

	bootstrap.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)
	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)

	value := "value"
	key := GetKeyRepresentationOfKademliaId(NewKademliaID("0000000000000000000000000000000000000002")) // Sets the key to be the same as kademlia2's id

	kademlia2.KademliaNode.DataStore.Insert(key, value)

	go bootstrap.Start()
	go kademlia1.Start()
	go kademlia2.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4020, false, "", 0)

	kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.Me)

	_, data, _ := kademlia.LookupData(key)

	assert.Equal(t, value, data)

}

func TestLookupDataFindsNoData(t *testing.T) {

	bootstrap := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 10030)

	kademlia1 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 10031)
	kademlia2 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 10032)

	bootstrap.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)
	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)

	value := "value"
	key := HashToKey(value)

	go bootstrap.Start()
	go kademlia1.Start()
	go kademlia2.Start()
	time.Sleep(time.Second)

	kademlia := NewKademlia("127.0.0.1", 4030, false, "", 0)

	kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.Me)

	list, _, _ := kademlia.LookupData(key)
	fmt.Println(list)

	doesContainAll := bootstrap.firstSetContainsAllContactsOfSecondSet(list, []Contact{kademlia1.KademliaNode.RoutingTable.Me, bootstrap.KademliaNode.RoutingTable.Me, kademlia.KademliaNode.RoutingTable.Me})
	assert.True(t, doesContainAll)
}
func TestStore(t *testing.T) {
	bootstrap := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1", 11000)

	kademlia1 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000001"), "127.0.0.1", 11001)
	kademlia2 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000002"), "127.0.0.1", 11002)
	kademlia3 := CreateMockedKademlia(NewKademliaID("0000000000000000000000000000000000000003"), "127.0.0.1", 11003)
	kademlia4 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "127.0.0.1", 11004)
	kademlia5 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000002"), "127.0.0.1", 11005)
	kademlia6 := CreateMockedKademlia(NewKademliaID("FFFFFFFF00000000000000000000000000000003"), "127.0.0.1", 11006)

	bootstrap.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.Me)
	bootstrap.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.Me)

	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.Me)
	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.Me)
	kademlia6.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)

	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.Me)
	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.Me)
	kademlia5.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.Me)

	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)
	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.Me)
	kademlia4.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.Me)

	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia5.KademliaNode.RoutingTable.Me)
	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)
	kademlia3.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)

	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia6.KademliaNode.RoutingTable.Me)
	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.Me)
	kademlia2.KademliaNode.RoutingTable.AddContact(kademlia1.KademliaNode.RoutingTable.Me)

	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia4.KademliaNode.RoutingTable.Me)
	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia3.KademliaNode.RoutingTable.Me)
	kademlia1.KademliaNode.RoutingTable.AddContact(kademlia2.KademliaNode.RoutingTable.Me)

	mockedKademlias := []Kademlia{
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

	kademlia.KademliaNode.RoutingTable.AddContact(bootstrap.KademliaNode.RoutingTable.Me)

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
			if kademlia.KademliaNode.RoutingTable.Me.ID == contact.ID {
				contactInMockedKademlias = true
				retrivedContent, err := kademlia.KademliaNode.DataStore.Get(key)

				if err != nil {
					assert.Fail(t, err.Error())
				}

				assert.True(t, retrivedContent == content)
				break

			}
		}
		if contactInMockedKademlias {
			_, err := kademlia.KademliaNode.DataStore.Get(key)

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

	var kademlias []*Kademlia
	var allContacts []Contact
	for i := 0; i < 30; i++ {
		port := 1101 + i
		kademlia := NewKademlia("localhost", port, false, "localhost", 1100)
		allContacts = append(allContacts, kademlia.KademliaNode.RoutingTable.Me)
		go kademlia.Start()
		time.Sleep(time.Microsecond * 50)

		kademlias = append(kademlias, kademlia)
	}
	time.Sleep(time.Second)

	target := allContacts[len(allContacts)-1].ID

	allContacts = append(allContacts, bootstrap.KademliaNode.RoutingTable.Me)
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

		containsAll := bootstrap.firstSetContainsAllContactsOfSecondSet(contacts, expectedClosest)
		assert.True(t, containsAll)

		if !containsAll {
			fmt.Print("Failed: ")
			fmt.Println(contacts)
		}

	}

}

func TestBig(t *testing.T) {
	time.Sleep(time.Second)

	bootstrap := NewKademlia("localhost", 1200, true, "", 0)
	go bootstrap.Start()
	time.Sleep(time.Second)

	var kademlias []*Kademlia
	var allContacts []Contact
	for i := 0; i < 30; i++ {
		port := 1201 + i
		kademlia := NewKademlia("localhost", port, false, "localhost", 1200)
		allContacts = append(allContacts, kademlia.KademliaNode.RoutingTable.Me)
		go kademlia.Start()
		time.Sleep(time.Microsecond * 100)

		kademlias = append(kademlias, kademlia)
	}
	time.Sleep(2 * time.Second)

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
			fmt.Println(kademlia.KademliaNode.RoutingTable.Me)
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
