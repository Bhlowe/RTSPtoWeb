package main

import (
	"github.com/gin-gonic/gin"
)

//HTTPAPIServerStreams function return client map
func HTTPAPIServerClients(c *gin.Context) {
	c.IndentedJSON(200, Message{Status: 1, Payload: Storage.ClientList()})
}

// Disconnect a user with userID c::uuid
func HTTPAPIServerDeleteClient(c *gin.Context) {
	// Disconnect a user with userID c::uuid
	c.IndentedJSON(200, Message{Status: 1})
}

//HTTPAPIServerStreams function return client map
func HTTPAPIServerAuthorizeClient(c *gin.Context) {
	// create a userInfoRecord with clientID and streamID and channelID;

	c.IndentedJSON(200, Message{Status: 1})
}

//HTTPAPIServerStreams function return client map
func HTTPAPIServerClientInfo(c *gin.Context) {
	// get ClientInfoStrct
	// if exists, return it as JSON.
	// if it does not exist, return error or status 0
	c.IndentedJSON(200, Message{Status: 1})
}

