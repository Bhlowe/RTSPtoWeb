package main

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/sirupsen/logrus"
)

var Storage = NewStreamCore()

//Default stream  type
const (
	MSE = iota
	WEBRTC
	RTSP
)

//Default stream status type
const (
	OFFLINE = iota
	ONLINE
)

//Default stream errors
var (
	Success                         = "success"
	ErrorStreamNotFound             = errors.New("stream not found")
	ErrorStreamAlreadyExists        = errors.New("stream already exists")
	ErrorStreamChannelAlreadyExists = errors.New("stream channel already exists")
	ErrorStreamNotHLSSegments       = errors.New("stream hls not ts seq found")
	ErrorStreamNoVideo              = errors.New("stream no video")
	ErrorStreamNoClients            = errors.New("stream no clients")
	ErrorStreamRestart              = errors.New("stream restart")
	ErrorStreamStopCoreSignal       = errors.New("stream stop core signal")
	ErrorStreamStopRTSPSignal       = errors.New("stream stop rtsp signal")
	ErrorStreamChannelNotFound      = errors.New("stream channel not found")
	ErrorStreamChannelCodecNotFound = errors.New("stream channel codec not ready, possible stream offline")
	ErrorStreamsLen0                = errors.New("streams len zero")
	ErrorClientNotAuthorized        = errors.New("client not authorized")
)

//StorageST main storage struct
type StorageST struct {
	mutex   sync.RWMutex
	Server  ServerST            `json:"server" groups:"api,config"`
	Streams map[string]StreamST `json:"streams,omitempty" groups:"api,config"`
	Clients ClientInfoMapST     // added BHL
}

//ServerST server storage section
type ServerST struct {
	Debug              bool         `json:"debug" groups:"api,config"`
	LogLevel           logrus.Level `json:"log_level" groups:"api,config"`
	HTTPDemo           bool         `json:"http_demo" groups:"api,config"`
	HTTPDebug          bool         `json:"http_debug" groups:"api,config"`
	HTTPLogin          string       `json:"http_login" groups:"api,config"`
	HTTPPassword       string       `json:"http_password" groups:"api,config"`
	HTTPDir            string       `json:"http_dir" groups:"api,config"`
	HTTPPort           string       `json:"http_port" groups:"api,config"`
	RTSPPort           string       `json:"rtsp_port" groups:"api,config"`
	HTTPS              bool         `json:"https" groups:"api,config"`
	HTTPSPort          string       `json:"https_port" groups:"api,config"`
	HTTPSCert          string       `json:"https_cert" groups:"api,config"`
	HTTPSKey           string       `json:"https_key" groups:"api,config"`
	HTTPSAutoTLSEnable bool         `json:"https_auto_tls" groups:"api,config"`
	HTTPSAutoTLSName   string       `json:"https_auto_tls_name" groups:"api,config"`
	Token              Token        `json:"token,omitempty" groups:"api,config"`
}

//Token auth
type Token struct {
	Enable bool `json:"enable,omitempty" groups:"api,config"`
}

//ServerST stream storage section
type StreamST struct {
	Name     string               `json:"name,omitempty" groups:"api,config"`
	Channels map[string]ChannelST `json:"channels,omitempty" groups:"api,config"`
}

type ChannelST struct {
	Name             string `json:"name,omitempty" groups:"api,config"`
	URL              string `json:"url,omitempty" groups:"api,config"`
	OnDemand         bool   `json:"on_demand,omitempty" groups:"api,config"`
	Debug            bool   `json:"debug,omitempty" groups:"api,config"`
	Status           int    `json:"status,omitempty" groups:"api"`
	runLock          bool
	codecs           []av.CodecData
	sdp              []byte
	signals          chan int
	hlsSegmentBuffer map[int]SegmentOld
	hlsSegmentNumber int
	clients          map[string]ClientST
	ack              time.Time
	hlsMuxer         *MuxerHLS `json:"-"`
}

//ClientST client storage section
type ClientST struct {
	mode              int
	signals           chan int
	outgoingAVPacket  chan *av.Packet
	outgoingRTPPacket chan *[]byte
	socket            net.Conn
}

// ClientInfo .. Keep separate or merge into ClientST?
type ClientInfoST struct {
	ClientId string    `json:"id,omitempty" `
	StreamId string    `json:"stream,omitempty" groups:"api,config"`
	Channel  string    `json:"channel,omitempty" groups:"api,config"`
	Mode     int       `json:"mode,omitempty" groups:"api,config"`
	Bytes    int       `json:"bytes,omitempty" groups:"api,config"`
	Start    time.Time `json:"start,omitempty" groups:"api,config"`
	LastTime time.Time `json:"last_time,omitempty" groups:"api,config"`
}

//Map of ClientInfo records.
type ClientInfoMapST struct {
	mutex         sync.RWMutex
	ClientInfoMap map[string]*ClientInfoST
}

//SegmentOld HLS cache section
type SegmentOld struct {
	dur  time.Duration
	data []*av.Packet
}
