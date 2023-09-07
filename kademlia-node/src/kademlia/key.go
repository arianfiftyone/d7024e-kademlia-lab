package kademlia

import (
	"crypto/sha1"
)

const (
	KeySize = 20
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
