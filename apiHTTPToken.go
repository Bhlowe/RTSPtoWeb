package main

import (
	"github.com/gin-gonic/gin"
)

//HTTPAPIServerStreams function return client map
func HTTPAPIServerTokens(c *gin.Context) {
	c.IndentedJSON(200, Message{Status: 1, Payload: Storage.ClientList()})
}

func HTTPAPIServerSetToken(c *gin.Context) {
	// get parameters:
	// token, stream_id
	// add to map
	c.IndentedJSON(200, Message{Status: 1, Payload: Storage.ClientList()})
}
