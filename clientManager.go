package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// This class keeps a list of clients (viewers) who ware watching a stream.
// A token based reservation system can be enabled to pre-authorize a client ID access to a stream:channel
// We also support disconnecting a client (viewer) by cid (clientID)

// TODO: Since I am updating the contents of the ClientInfoMap, should probably
// use a pointer to the object... not just replace the item (ClientInfoST) each time it is modified.
//

func (obj *ClientInfoMapST) init() {
	fmt.Println("init clients", time.Now().Unix(), obj.ClientInfoMap)
	obj.ClientInfoMap = make(map[string]ClientInfoST)

	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for range ticker.C {
			obj.removeStaleClients()
		}
	}()
}

/*func (obj *ClientInfoMapST) streamChannelID(stream string, channel string) string {
	return stream + ":" + channel
}
*/
func (obj *ClientInfoMapST) requireClientID() bool {
	return false
}

func (obj *ClientInfoMapST) addClient(stream string, channel string, cid string) error {

	if !Storage.StreamChannelExist(stream, channel) {
		return ErrorStreamNotFound
	}
	obj.mutex.Lock()
	defer obj.mutex.Unlock()

	_, exist := obj.ClientInfoMap[cid]

	if exist {
		return ErrorStreamAlreadyExists // TODO, fix
	}

	var out ClientInfoST
	tm := time.Now()

	if cid == "" {
		copy, _ := generateUUID()
		cid = copy
	}

	out.Channel = channel
	out.StreamId = stream
	out.ClientId = cid
	out.LastTime = tm
	out.Bytes = 0
	out.Start = tm
	out.Mode = -1
	obj.ClientInfoMap[cid] = out

	return nil

}

func (obj *ClientInfoMapST) checkOrCreateCID(stream string, channel string, cid string, mode int) (ClientInfoST, error) {

	var out ClientInfoST

	if !Storage.StreamChannelExist(stream, channel) {
		return out, ErrorStreamNotFound
	}

	fmt.Println("checkOrCreateCID " + cid)

	tm := time.Now()
	obj.mutex.Lock()
	defer obj.mutex.Unlock()

	if stream == "" || channel == "" {
		return out, ErrorStreamNotFound
	}

	if cid == "" {
		copy, _ := generateUUID()
		cid = copy
	}

	v, exist := obj.ClientInfoMap[cid]
	if exist {
		// double check stream and channel.
		if stream != v.StreamId || channel != v.Channel {
			return out, ErrorClientNotAuthorized
		}
		// update mode..
		v.Mode = mode
		obj.ClientInfoMap[cid] = out

		return v, nil
	}

	if obj.requireClientID() {
		return out, ErrorClientNotAuthorized
	}

	out.Channel = channel
	out.StreamId = stream
	out.ClientId = cid
	out.LastTime = tm
	out.Bytes = 0
	out.Start = tm
	out.Mode = mode

	obj.ClientInfoMap[cid] = out

	return out, nil
}

func (obj *ClientInfoMapST) getClientInfo(cid string) (ClientInfoST, bool) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	out, exist := obj.ClientInfoMap[cid]
	return out, exist // obj.ClientInfoMap[cid]
}

func (obj *ClientInfoMapST) hasClientInfo(cid string) bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	_, exist := obj.ClientInfoMap[cid]
	return exist
}

// Return a copy of the map.
func (obj *ClientInfoMapST) copyMap() map[string]ClientInfoST {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	tmp := make(map[string]ClientInfoST)
	for key, value := range obj.ClientInfoMap {
		tmp[key] = value
	}
	return tmp
}

func (obj *ClientInfoMapST) logPackets(cid string, bytes int) error {

	obj.mutex.Lock()
	v, exist := obj.ClientInfoMap[cid]
	if !exist {
		return ErrorStreamChannelNotFound
	}
	//if (v.Disconnected) {
	//	return ErrorStreamChannelNotFound
	//}
	v.Bytes += bytes
	v.LastTime = time.Now()
	obj.ClientInfoMap[cid] = v // TODO: Use pointer instead?
	obj.mutex.Unlock()

	// fmt.Println("logPackets:", cid, bytes, v.LastTime, v.Bytes)
	return nil
}

func (t ClientInfoST) String() string {
	e, err := json.Marshal(t)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(e)
}

func (t ClientInfoST) Seconds() float64 {
	dur := t.LastTime.Sub(t.Start)
	return dur.Seconds()
}

func (t ClientInfoST) Now() time.Time {
	return time.Now()
}

// return 1 if removed, return 0 if not removed.
func (obj *ClientInfoMapST) streamClosed(cid string, reason string, err error) int {

	obj.mutex.Lock()
	defer obj.mutex.Unlock()

	_, exist := obj.ClientInfoMap[cid]
	if !exist {
		return 0
	}
	fmt.Println("streamClosed", cid, reason, err)
	delete(obj.ClientInfoMap, cid)
	return 1
}

func (obj *ClientInfoMapST) removeStaleClients() {
	if len(obj.ClientInfoMap) == 0 {
		return
	}

	var expired []string // an empty list
	now := time.Now()
	expireBefore := now.Add(time.Duration(-60) * time.Second)

	obj.mutex.Lock()

	for key, value := range obj.ClientInfoMap {
		if value.LastTime.Before(expireBefore) {
			expired = append(expired, key)
			delta := value.LastTime.Sub(value.Start)
			// value.LastTime - value.Start
			fmt.Println("remove stale client, last,start,delta, seconds", value.String(), value.LastTime, value.Start, delta, value.Seconds())

			// ok to delete item from map while iterating through it? delete now.. otherwise.. delete later as list.
		}
	}
	obj.mutex.Unlock()

	if len(expired) > 0 {
		for _, cid := range expired {
			obj.streamClosed(cid, "expired stream", nil)
		}
		fmt.Println("removed Stale Clients:", len(expired), len(obj.ClientInfoMap))
	}
}
