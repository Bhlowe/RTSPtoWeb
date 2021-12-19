package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//HTTPAPIServerStreams function return client map
func HTTPAPIServerClients(c *gin.Context) {
	c.IndentedJSON(200, Message{Status: 1, Payload: Storage.Clients.copyMap()})
}

// Disconnect a user with userID c::cid
func HTTPAPIServerDisconnectUser(c *gin.Context) {
	// Disconnect a user with clientID ("cid")
	status := Storage.Clients.streamClosed(c.Param("cid"), "api_disconnect", nil)
	c.IndentedJSON(200, Message{Status: status})
}

//HTTPAPIServerStreams function return client map
// 	privat.GET("/client/add/:stream", HTTPAPIServerAuthorizeUser)
func HTTPAPIServerAuthorizeUser(c *gin.Context) {
	// TODO: Fix Logs
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "XXX", // TODO
		"stream":  c.Param("stream"),
		"channel": c.Param("channel"),
		"cid":     c.Param("cid"),
		"func":    "HTTPAPIServerAuthorizeUser",
	})
	var cid = c.Param("cid")
	if len(cid) == 0 {
		cid, _ = generateUUID()
	}
	var channel = c.Param("channel")
	if len(channel) == 0 {
		channel = "0"
	}
	stream := c.Param("stream")
	err := Storage.Clients.addClient(stream, channel, cid)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "HTTPAPIServerAuthorizeUser",
		}).Errorln(err.Error())

	} else {
		info, found := Storage.Clients.getClientInfo(cid)
		if found {
			c.IndentedJSON(200, Message{Status: 1, Payload: info})
		}
	}
}

//HTTPAPIServerStreams function return client map
func HTTPAPIServerClientInfo(c *gin.Context) {
	// get ClientInfo

	info, found := Storage.Clients.getClientInfo(c.Param("cid"))
	if found {
		// j, err2 := json.Marshal(info)
		c.IndentedJSON(200, Message{Status: 1, Payload: info})

	} else {
		c.IndentedJSON(500, Message{Status: 0})
	}
	// if exists, return it as JSON.
	// if it does not exist, return error or status 0
}
