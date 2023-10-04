package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
)

const (
	KeySize = IDLength
)

type Key struct {
	Hash [KeySize]byte
}

func NewKey(value string) *Key {
	bytes := []byte(value)
	hash := sha1.Sum(bytes)

	return &Key{
		hash,
	}
}

func (key *Key) GetHashString() string {
	return hex.EncodeToString(key.Hash[:])
}

func (key *Key) GetKademliaIdRepresentationOfKey() *KademliaID {
	return GenerateNewKademliaID(key.GetHashString())
}

func GetKeyRepresentationOfKademliaId(id *KademliaID) *Key {
	str := id.String()
	decoded, _ := hex.DecodeString(str)
	key := Key{}
	for i := 0; i < IDLength; i++ {
		key.Hash[i] = decoded[i]
	}
	return &key
}
