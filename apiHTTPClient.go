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
func HTTPAPIServerAuthorizeUser(c *gin.Context) {
	// TODO: Fix Logs
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "XXX", // TODO
		"stream":  c.Param("stream"),
		"channel": c.Param("channel"),
		"cid":     c.Param("cid"),
		"func":    "HTTPAPIServerAuthorizeUser",
	})

	err := Storage.Clients.addClient(c.Param("cid"), c.Param("stream"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "HTTPAPIServerAuthorizeUser",
		}).Errorln(err.Error())

	} else {
		c.IndentedJSON(200, Message{Status: 1})
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
