package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arianfiftyone/src/kademlia"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	KeySize = 20
)

// KademliaMock is a mock implementation of the kademlia.Kademlia interface.
type KademliaMock struct {
	mock.Mock
	DataStore *kademlia.DataStore
}

func (KademliaMock *KademliaMock) Start() {}

func (KademliaMock *KademliaMock) Join() {}

func (KademliaMock *KademliaMock) Store(content string) (*kademlia.Key, error) {
	return kademlia.NewKey(content), nil
}

func (KademliaMock *KademliaMock) GetKademliaNode() *kademlia.KademliaNode {
	return nil
}

func (KademliaMock *KademliaMock) FirstSetContainsAllContactsOfSecondSet(first []kademlia.Contact, second []kademlia.Contact) bool {
	return false
}

func (KademliaMock *KademliaMock) LookupContact(targetId *kademlia.KademliaID) ([]kademlia.Contact, error) {
	return nil, nil
}

func (KademliaMock *KademliaMock) LookupData(key *kademlia.Key) ([]kademlia.Contact, string, error) {
	content, err := KademliaMock.DataStore.Get(key)
	return nil, content, err
}

func (KademliaMock *KademliaMock) Forget(key *kademlia.Key) error {
	return nil
}

func TestGetObjectValidHash(t *testing.T) {

	kademliaMock := new(KademliaMock)

	// Create a new instance of DataStore
	dataStore := kademlia.NewDataStore()
	value := "kademlia"
	key := kademlia.NewKey(value)
	hash := key.GetHashString()
	dataStore.Insert(key, value)

	kademliaMock.DataStore = &dataStore
	api := NewKademliaAPI(kademliaMock)

	// Set up the Gin context for testing
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/objects/", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{
		Key:   "hash",
		Value: hash,
	})
	c.Request = req

	// Call the GetObject handler
	api.GetObject(c)

	expectedJSON := `{"value": "%s"}`
	expectedJSON = fmt.Sprintf(expectedJSON, value)

	// Verify the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, w.Body.String(), expectedJSON)
}

func TestGetObjectInvalidHash(t *testing.T) {

	kademliaMock := new(KademliaMock)
	api := NewKademliaAPI(kademliaMock)

	value := "kademlia"
	key := kademlia.NewKey(value)
	hash := key.GetHashString()
	invalidHash := hash + "extra"

	// Set up the Gin context for testing
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/objects/"+invalidHash, nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the GetObject handler
	api.GetObject(c)

	// Verify the response for an invalid hash
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Invalid hash length"}`, w.Body.String())

}

func TestPostObjectValidInput(t *testing.T) {

	kademliaMock := new(KademliaMock)
	api := NewKademliaAPI(kademliaMock)

	// Set up the Gin context for testing
	w := httptest.NewRecorder()

	input := `{"Value":"kademlia"}`
	decoder := strings.NewReader(input)
	req, _ := http.NewRequest("POST", "/objects", decoder)

	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the PostObject handler
	api.PostObject(c)

	value := "kademlia"
	key := kademlia.NewKey(value)
	hash := key.GetHashString()
	expectedJSON := `{"hash": "%s"}`
	expectedJSON = fmt.Sprintf(expectedJSON, hash)

	// Verify the response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, w.Body.String(), expectedJSON)
}

func TestPostObjectInvalidInput(t *testing.T) {

	kademliaMock := new(KademliaMock)
	api := NewKademliaAPI(kademliaMock)

	// Set up the Gin context for testing with an invalid JSON input
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/objects", strings.NewReader(`invalidjson`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the PostObject handler with invalid input
	api.PostObject(c)

	// Verify the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"Internal server error"}`, w.Body.String())
}
