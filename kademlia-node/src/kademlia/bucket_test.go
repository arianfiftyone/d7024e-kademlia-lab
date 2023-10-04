package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewRandomKademliaID(), "", 0)
	contact1 := NewContact(NewRandomKademliaID(), "", 0)
	contact2 := NewContact(NewRandomKademliaID(), "", 0)
	contact3 := NewContact(NewRandomKademliaID(), "", 0)

	bucket.AddContact(contact)
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	bucket.AddContact(contact3)

	assert.True(t, bucket.Contains(contact))
}

func TestDoesNotContain(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewRandomKademliaID(), "", 0)
	contact1 := NewContact(NewRandomKademliaID(), "", 0)
	contact2 := NewContact(NewRandomKademliaID(), "", 0)
	contact3 := NewContact(NewRandomKademliaID(), "", 0)

	bucket.AddContact(contact1)
	bucket.AddContact(contact2)
	bucket.AddContact(contact3)

	assert.False(t, bucket.Contains(contact))
}
