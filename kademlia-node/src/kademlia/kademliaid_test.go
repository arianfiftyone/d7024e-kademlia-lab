package kademlia

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomKademliaidInRange(t *testing.T) {
	lowerBound := NewKademliaID("000000000000000000000000000000000000001E")
	highBound := NewKademliaID("0000000000000000000000000000000000000020")
	kademliaId, err := NewRandomKademliaIDInRange(lowerBound, highBound)
	fmt.Println(kademliaId)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.True(t, kademliaId.Less(highBound) || kademliaId.Equals(highBound))
	assert.True(t, lowerBound.Less(kademliaId) || lowerBound.Equals(kademliaId))

}

func TestNewRandomKademliaidInRange2(t *testing.T) {
	lowerBound := NewRandomKademliaID()
	highBound := NewRandomKademliaID()

	if highBound.Less(lowerBound) {
		temp := lowerBound
		lowerBound = highBound
		highBound = temp
	}

	kademliaId, err := NewRandomKademliaIDInRange(lowerBound, highBound)
	fmt.Println(lowerBound)

	fmt.Println(kademliaId)

	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.True(t, kademliaId.Less(highBound) || kademliaId.Equals(highBound))
	assert.True(t, lowerBound.Less(kademliaId) || lowerBound.Equals(kademliaId))

}
