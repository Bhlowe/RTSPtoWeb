package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/deepch/vdk/format/ts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//HTTPAPIServerStreamHLSM3U8 send client m3u8 play list
func HTTPAPIServerStreamHLSM3U8(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_hls",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"cid":     c.Query("cid"),
		"func":    "HTTPAPIServerStreamHLSM3U8",
	})

	var cid = c.Query("cid")
	if len(cid) == 0 {
		fmt.Println("no cid for "+c.Request.RequestURI, cid)
	}

	info, err := Storage.Clients.checkClient(c.Param("uuid"), c.Param("channel"), cid, RTSP)
	if err != nil {
		fmt.Println("HTTPAPIServerStreamHLSM3U8 no cid found", cid)
		Storage.Clients.checkClient(c.Param("uuid"), c.Param("channel"), cid, RTSP)

		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "checkClient",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	cid = info.ClientId

	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}
	// BHL, get or create clientInfoRecord
	// Fail if unable to get.

	c.Header("Content-Type", "application/x-mpegURL")
	Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))

	//If stream mode on_demand need wait ready segment's
	for i := 0; i < 40; i++ {
		index, seq, err := Storage.StreamHLSm3u8(c.Param("uuid"), c.Param("channel"), cid)
		if err != nil {
			c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
			requestLogger.WithFields(logrus.Fields{
				"call": "StreamHLSm3u8",
			}).Errorln(err.Error())
			return
		}
		if seq >= 6 {

			_, err := c.Writer.Write([]byte(index))

			if err != nil {
				c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
				requestLogger.WithFields(logrus.Fields{
					"call": "Write",
				}).Errorln(err.Error())
				return
			}
			Storage.Clients.logPackets(cid, len([]byte(index)))
			// increment bytes sent in clientInfoRecord.bytes
			return
		}
		time.Sleep(1 * time.Second)
	}
}

//HTTPAPIServerStreamHLSTS send client ts segment
func HTTPAPIServerStreamHLSTS(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_hls",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"cid":     c.Query("cid"),
		"func":    "HTTPAPIServerStreamHLSTS",
	})

	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}
	codecs, err := Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamCodecs",
		}).Errorln(err.Error())
		return
	}
	outfile := bytes.NewBuffer([]byte{})
	Muxer := ts.NewMuxer(outfile)
	Muxer.PaddingToMakeCounterCont = true
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	seqData, err := Storage.StreamHLSTS(c.Param("uuid"), c.Param("channel"), stringToInt(c.Param("seq")))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamHLSTS",
		}).Errorln(err.Error())
		return
	}
	if len(seqData) == 0 {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotHLSSegments.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "seqData",
		}).Errorln(ErrorStreamNotHLSSegments.Error())
		return
	}
	for _, v := range seqData {
		v.CompositionTime = 1
		err = Muxer.WritePacket(*v)
		if err != nil {
			c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
			requestLogger.WithFields(logrus.Fields{
				"call": "WritePacket",
			}).Errorln(err.Error())
			return
		}
	}
	err = Muxer.WriteTrailer()
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteTrailer",
		}).Errorln(err.Error())
		return
	}
	_, err = c.Writer.Write(outfile.Bytes())
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "Write",
		}).Errorln(err.Error())
		return
	}
	Storage.Clients.logPackets(c.Query("cid"), len(outfile.Bytes()))

}
