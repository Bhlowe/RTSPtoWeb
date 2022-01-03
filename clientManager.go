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
	fmt.Println("init clients")
	obj.ClientInfoMap = make(map[string]*ClientInfoST)

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

// this gets info about the map.
// also used for sanity checking the map.
func (obj *ClientInfoMapST) inspect() string {
	var s = ""
	if true {
		/*		obj.mutex.Lock()
				defer obj.mutex.Unlock()
		*/
		for key, value := range obj.ClientInfoMap {
			if len(value.StreamId) == 0 || len(value.Channel) == 0 {
				fmt.Println("Error with map entry:" + key)
			}
			line := key + "->" + value.StreamId + " " + value.ClientId + " " + value.Channel
			s += line + "\n"
		}
	}
	return s

}
func (obj *ClientInfoMapST) addClient(stream string, channel string, cid string) error {
	fmt.Println("addClient", stream, cid, channel)

	/* allow token for stream that doesn't exist (yet)
	if !Storage.StreamChannelExist(stream, channel) {
		return ErrorStreamNotFound
	}
	*/

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

	if (len(stream) == 0) || len(channel) == 0 {

	}

	out.Channel = channel
	out.StreamId = stream
	out.ClientId = cid
	out.LastTime = tm
	out.Bytes = 0
	out.Start = tm
	out.Mode = -1
	obj.ClientInfoMap[cid] = &out

	fmt.Println("addedClient", stream, cid)
	fmt.Println(obj.inspect())

	return nil
}

func (obj *ClientInfoMapST) checkClient(stream string, channel string, cid string, mode int) (*ClientInfoST, error) {

	if !Storage.StreamChannelExist(stream, channel) {
		fmt.Println("checkClient no stream for ", stream, cid)

		return nil, ErrorStreamNotFound
	}

	tm := time.Now()
	obj.mutex.Lock()
	defer obj.mutex.Unlock()

	if stream == "" || channel == "" {
		return nil, ErrorStreamNotFound
	}

	if cid == "" {
		copy, _ := generateUUID()
		cid = copy
		if obj.requireClientID() {
			fmt.Println("empty cid not allowed")
			return nil, ErrorClientNotAuthorized
		}

	}

	v, exist := obj.ClientInfoMap[cid]
	if exist {
		// double check stream and channel.
		if stream != v.StreamId || channel != v.Channel {
			return nil, ErrorClientNotAuthorized
		}
		// update mode..
		v.Mode = mode
		v.LastTime = tm

		// fmt.Println("checkClient authorized", cid, obj.inspect())

		return v, nil
	}

	fmt.Println("No CID found: ", cid, stream, obj.inspect())

	if obj.requireClientID() {
		return nil, ErrorClientNotAuthorized
	}

	var out ClientInfoST
	out.Channel = channel
	out.StreamId = stream
	out.ClientId = cid
	out.LastTime = tm
	out.Bytes = 0
	out.Start = tm
	out.Mode = mode

	obj.ClientInfoMap[cid] = &out
	fmt.Println("checkClient no auth needed", cid, obj.inspect())
	obj.inspect()

	return &out, nil
}

func (obj *ClientInfoMapST) getClient(cid string) (ClientInfoST, bool) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	out, exist := obj.ClientInfoMap[cid]
	return *out, exist // obj.ClientInfoMap[cid]
}

func (obj *ClientInfoMapST) hasClient(cid string) bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	_, exist := obj.ClientInfoMap[cid]
	return exist
}

func (obj *ClientInfoMapST) len() int {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return len(obj.ClientInfoMap)
}

// Return a copy of the map.
func (obj *ClientInfoMapST) copyMap() map[string]ClientInfoST {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	tmp := make(map[string]ClientInfoST)
	for key, value := range obj.ClientInfoMap {
		tmp[key] = *value
	}
	return tmp
}

func (t *ClientInfoST) logPackets(bytes int) {
	t.Bytes += bytes
	t.LastTime = time.Now()
}

func (obj *ClientInfoMapST) logPackets(cid string, bytes int) error {

	obj.mutex.Lock()
	defer obj.mutex.Unlock()

	v, exist := obj.ClientInfoMap[cid]
	if !exist {
		return ErrorStreamChannelNotFound
	}

	v.logPackets(bytes)
	obj.inspect()
	obj.ClientInfoMap[cid] = v

	if false {
		copy := obj.ClientInfoMap[cid]
		if copy.Bytes != v.Bytes {
			fmt.Println("failed check.")
		}
		// fmt.Println("logPackets:", cid, bytes, copy.LastTime, copy.Bytes, copy.Seconds())
	}

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
	fmt.Println("streamClosed", reason, cid)

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
			if false {
				delta := value.LastTime.Sub(value.Start)
				fmt.Println("remove stale client, last,start,delta, seconds", value.String(), value.LastTime, value.Start, delta, value.Seconds())
			}

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
