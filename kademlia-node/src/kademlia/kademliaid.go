package kademlia

import (
	"encoding/hex"
	"errors"
	"math/rand"
)

// the static number of bytes in a KademliaID
const IDLength = 20

// type definition of a KademliaID
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}

	return &newKademliaID
}

// NewRandomKademliaID returns a new instance of a random KademliaID,
// change this to a better version if you like
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

// NewRandomKademliaIDInRange returns a new instance of a random KademliaID, within a given range
func NewRandomKademliaIDInRange(lowBound *KademliaID, highBound *KademliaID) (*KademliaID, error) {
	if !lowBound.Less(highBound) {
		return nil, errors.New("the low bound is not lower than the high bound")
	}

	newKademliaID := &KademliaID{}
	for i := 0; i < IDLength; i++ {
		min := lowBound[i]
		max := highBound[i]
		diff := max - min
		randomInt := rand.Intn(int(diff) + 1)
		randomInt = randomInt + int(min)
		newKademliaID[i] = byte(randomInt)
	}

	return newKademliaID, nil

}

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// CalcDistance returns a new instance of a KademliaID that is built
// through a bitwise XOR operation betweeen kademliaID and target
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a KademliaID
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}
