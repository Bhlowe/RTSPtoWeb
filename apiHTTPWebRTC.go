package main

import (
	"fmt"
	"time"

	webrtc "github.com/deepch/vdk/format/webrtcv3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//HTTPAPIServerStreamWebRTC stream video over WebRTC
func HTTPAPIServerStreamWebRTC(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_webrtc",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"cid":     c.Query("cid"),
		"func":    "HTTPAPIServerStreamWebRTC",
	})
	cid := c.Query("cid")

	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}
	Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
	codecs, err := Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamCodecs",
		}).Errorln(err.Error())
		return
	}

	// BHL, added to maintain list of clients and check for authorization
	info, err := Storage.Clients.checkClient(c.Param("uuid"), c.Param("channel"), cid, WEBRTC)
	// TODO: Make better error.
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}

	muxerWebRTC := webrtc.NewMuxer(webrtc.Options{})
	answer, err := muxerWebRTC.WriteHeader(codecs, c.PostForm("data"))
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		Storage.Clients.removeClient(info.StreamId, "WriteHeader", err)
		return
	}
	_, err = c.Writer.Write([]byte(answer))
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "Write",
		}).Errorln(err.Error())
		Storage.Clients.removeClient(info.StreamId, "WriteAnswer", err)
		return
	}
	go func() {

		if err != nil {
			c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
			requestLogger.WithFields(logrus.Fields{
				"call": "ClientAdd",
			}).Errorln(err.Error())
			return
		}

		cid, ch, _, err := Storage.ClientAdd(info)

		fmt.Println("webrtc Storage.ClientAdd ", info.ClientId)

		if err != nil {
			c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
			requestLogger.WithFields(logrus.Fields{
				"call": "ClientAdd",
			}).Errorln(err.Error())
			return
		}

		defer Storage.Clients.removeClient(info.ClientId, "deferred webrtc close", nil)

		defer Storage.ClientDelete(c.Param("uuid"), cid, c.Param("channel"))
		var videoStart bool
		noVideo := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-noVideo.C:
				c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNoVideo.Error()})
				requestLogger.WithFields(logrus.Fields{
					"call": "ErrorStreamNoVideo",
				}).Errorln(ErrorStreamNoVideo.Error())
				return
			case pck := <-ch:
				if pck.IsKeyFrame {
					noVideo.Reset(10 * time.Second)
					videoStart = true
				}
				if !videoStart {
					continue
				}

				err = muxerWebRTC.WritePacket(*pck)
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "WritePacket",
					}).Errorln(err.Error())
					Storage.Clients.removeClient(info.ClientId, "WritePacket", err)
					return
				}

				err = Storage.Clients.logPackets(cid, len(pck.Data))
				if err != nil {
					requestLogger.WithFields(logrus.Fields{
						"call": "logPackets",
					}).Errorln(err.Error())
					return
				}

			}
		}
	}()
}
