package main

import (
	"fmt"
	"time"
)

func (obj *ClientInfoMapST) init() {

	obj.ClientInfoMap = make(map[string]ClientInfoST)
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for range ticker.C {
			obj.removeStaleClients()
		}
	}()

}

func (obj *ClientInfoMapST) checkOrCreateCID(stream string, channel string, cid string, mode int) (ClientInfoST, error) {

	// TODO: Move this to constructor
	if obj.ClientInfoMap == nil {
		obj.init()
	}

	var out ClientInfoST

	if !Storage.StreamChannelExist(stream, channel) {
		return out, ErrorStreamNotFound
	}

	fmt.Println("checkOrCreateCID " + cid)

	tm := time.Now().Nanosecond()
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()

	if stream == "" || channel == "" {
		return out, ErrorStreamNotFound
	}

	if cid == "" {
		copy, _ := generateUUID()
		cid = copy
	}

	v, exist := obj.ClientInfoMap[cid]
	if exist {
		return v, nil
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

func (obj *ClientInfoMapST) getClientInfo(cid string) ClientInfoST {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.ClientInfoMap[cid]
}

func (obj *ClientInfoMapST) hasClientInfo(cid string) bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	_, exist := obj.ClientInfoMap[cid]
	return exist
}

func (obj *ClientInfoMapST) removeClient(cid string) {

	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	delete(obj.ClientInfoMap, cid)
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

func (obj *ClientInfoMapST) Test() error {
	var cid = "1234"

	var info, err = obj.checkOrCreateCID("stream", "1", cid, MSE)
	if err != nil {
		return err
	}
	if info.ClientId != cid {
		return ErrorStreamNotFound
	}
	if !obj.hasClientInfo(cid) {
		return ErrorStreamNotFound
	}
	obj.removeClient(cid)
	return nil
}

func (obj *ClientInfoMapST) logPackets(cid string, bytes int) error {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	v, exist := obj.ClientInfoMap[cid]
	if !exist {
		return ErrorStreamChannelNotFound
	}
	//if (v.Disconnected) {
	//	return ErrorStreamChannelNotFound
	//}
	v.Bytes += bytes
	v.LastTime = time.Now().Nanosecond()
	return nil
}

func (obj *ClientInfoMapST) removeStaleClients() {
	fmt.Println("removeStaleClients !!")

	var expired []string // an empty list
	now := time.Now()
	expireBefore := now.Add(time.Duration(-60) * time.Second).Nanosecond()
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()

	for key, value := range obj.ClientInfoMap {
		if value.LastTime < expireBefore {
			expired = append(expired, key)
			// ok to delete item from map while iterating through it? delete now.. otherwise.. delete later as list.
		}
	}

	if len(expired) > 0 {
		obj.mutex.RLock()
		defer obj.mutex.RUnlock()
		for _, cid := range expired {
			delete(obj.ClientInfoMap, cid)
		}
		fmt.Println("removed Stale Clients:", len(expired), len(obj.ClientInfoMap))
	}

}
