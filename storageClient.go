package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

func (obj *StorageST) ClientAdd(info ClientInfoST) (string, chan *av.Packet, chan *[]byte, error) {
	return obj._ClientAdd(info.StreamId, info.Channel, info.ClientId, info.Mode)
}

//ClientAdd Add New Client to Translations
func (obj *StorageST) _ClientAdd(streamID string, channelID string, cid string, mode int) (string, chan *av.Packet, chan *[]byte, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return "", nil, nil, ErrorStreamNotFound
	}

	/*	if err != nil {
			return "", nil, nil, ErrorClientNotAuthorized
		}
	*/
	//Generate UUID client
	chAV := make(chan *av.Packet, 2000)
	chRTP := make(chan *[]byte, 2000)
	channelTmp, ok := streamTmp.Channels[channelID]
	if !ok {
		return "", nil, nil, ErrorStreamNotFound
	}

	channelTmp.clients[cid] = ClientST{mode: mode, outgoingAVPacket: chAV, outgoingRTPPacket: chRTP, signals: make(chan int, 100)}
	channelTmp.ack = time.Now()
	streamTmp.Channels[channelID] = channelTmp
	obj.Streams[streamID] = streamTmp
	return cid, chAV, chRTP, nil

}

//ClientDelete Delete Client
func (obj *StorageST) ClientDelete(streamID string, cid string, channelID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[streamID]; ok {
		delete(obj.Streams[streamID].Channels[channelID].clients, cid)
	}
}

//ClientHas check is client ext
func (obj *StorageST) ClientHas(streamID string, channelID string) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return false
	}
	channelTmp, ok := streamTmp.Channels[channelID]
	if !ok {
		return false
	}
	if time.Now().Sub(channelTmp.ack).Seconds() > 30 {
		return false
	}
	return true
}
