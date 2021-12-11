package main

import (
	"github.com/gin-gonic/gin"
)

//HTTPAPIServerStreams function return client map
func HTTPAPIServerClients(c *gin.Context) {
	c.IndentedJSON(200, Message{Status: 1, Payload: Storage.ClientList()})
}
