package kademlia

import (
	"crypto/sha1"
)

const (
	KeySize = IDLength
)

type Key struct {
	Hash [KeySize]byte
}

func HashToKey(value string) *Key {
	bytes := []byte(value)
	hash := sha1.Sum(bytes)

	return &Key{
		hash,
	}
}

func (key *Key) GetKademliaIdRepresentationOfKey() *KademliaID {
	return NewKademliaID(string(key.Hash[:]))
}
