package kademlia

type Kademlia struct {
	Network *Network
}

const (
	BootstrapKademliaID = "FFFFFFFF00000000000000000000000000000000"
)

func NewKademlia(ip string, port int) *Kademlia {

	kademliaNode := NewKademliaNode(ip, port)
	network := Network{
		ip,
		port,
		&MessageHandlerImplementation{
			kademliaNode,
		},
	}
	return &Kademlia{
		Network: &network,
	}

}

func (kademlia *Kademlia) Start() {
	err := kademlia.Network.Listen()
	if err != nil {
		panic(err)

	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
