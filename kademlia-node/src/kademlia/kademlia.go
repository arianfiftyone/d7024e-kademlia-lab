package kademlia

import "fmt"

type Kademlia struct {
	Network          *Network
	KademliaNode     *KademliaNode
	isBootsrap       bool
	bootsrtapContact *Contact
}

const (
	BootstrapKademliaID = "FFFFFFFF00000000000000000000000000000000"
)

func NewKademlia(ip string, port int, isBootsrap bool, bootstrapIp string, bootstratPort int) *Kademlia {

	kademliaNode := NewKademliaNode(ip, port, isBootsrap)
	network := &Network{
		ip,
		port,
		&MessageHandlerImplementation{
			kademliaNode,
		},
	}
	var contact Contact
	if !isBootsrap {
		contact = NewContact(
			NewKademliaID(BootstrapKademliaID),
			bootstrapIp,
			bootstratPort,
		)
	}
	return &Kademlia{
		Network:          network,
		KademliaNode:     kademliaNode,
		isBootsrap:       isBootsrap,
		bootsrtapContact: &contact,
	}

}

func (kademlia *Kademlia) Start() {
	if !kademlia.isBootsrap {
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

	if kademlia.isBootsrap {
		fmt.Println("You are the bootstrap node!")
		return

	}

	err := kademlia.Network.SendPingMessage(kademlia.bootsrtapContact)
	if err != nil {
		return
	}

	kademlia.KademliaNode.RoutingTable.AddContact(*kademlia.bootsrtapContact)

	kademlia.LookupContact(&kademlia.KademliaNode.RoutingTable.me)

}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// A node finds k nodes to check if they are close to the hash
}
