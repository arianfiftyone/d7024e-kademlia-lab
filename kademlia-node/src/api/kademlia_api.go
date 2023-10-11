package api

import (
	"net/http"

	"github.com/arianfiftyone/src/kademlia"
	"github.com/arianfiftyone/src/logger"
	"github.com/gin-gonic/gin"
)

type KademliaAPI struct {
	kademlia kademlia.Kademlia
}

func NewKademliaAPI(kademlia kademlia.Kademlia) *KademliaAPI {
	return &KademliaAPI{
		kademlia: kademlia,
	}
}

type ValueDTO struct {
	Value string `json:"value"`
}

type HashDTO struct {
	Hash string `json:"hash"`
}

// StartAPI initializes and starts the REST API using Gin.
func (kademliaAPI KademliaAPI) StartAPI() {
	logger.Log("Starting REST API...")

	router := gin.Default()

	router.GET("/objects/:hash", kademliaAPI.GetObject)
	router.POST("/objects", kademliaAPI.PostObject)

	err := router.Run(":50000")
	if err != nil {
		logger.Log("error when listening to the http server " + err.Error())
	}
}

// GetObject handles GET requests for object retrieval.
func (kademliaAPI KademliaAPI) GetObject(ctx *gin.Context) {
	// Extract the hash from the URL parameter
	hashParam := ctx.Param("hash")

	if len(hashParam) != 40 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hash length"})
	} else {

		newKademliaID, err := kademlia.NewKademliaID(hashParam)
		if err != nil {
			logger.Log("error when creating a new KademliaID: " + err.Error())
			return
		}

		// Attempt to find the value associated with the `KademliaID` in the Kademlia network
		_, value, err := kademliaAPI.kademlia.LookupData(kademlia.GetKeyRepresentationOfKademliaId(newKademliaID))

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "404 page not found"})
		} else {
			res := ValueDTO{Value: value}
			ctx.JSON(http.StatusOK, res)
		}
	}
}

// PostObject handles POST requests for object storage.
func (kademliaAPI KademliaAPI) PostObject(ctx *gin.Context) {
	var valueDTO ValueDTO

	if err := ctx.ShouldBindJSON(&valueDTO); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Store the value in the Kademlia network and get the associated key
	key, err := kademliaAPI.kademlia.Store(valueDTO.Value)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing object"})
		return
	}

	res := HashDTO{Hash: key.GetHashString()}

	ctx.Header("Location", "/objects/"+key.GetHashString())
	ctx.IndentedJSON(http.StatusCreated, res)
}
