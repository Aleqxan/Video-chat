package webrtc

import (

)

func RoomConn(c *websocket.Conn, p *Peers){
	var config webrtc.Configuration

	PeerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Print(err)
		return
	}

	newPeer := PeerConnectionState{
		PeerConnection: PeerConnection
		WebSocket: &ThreadSafeWriter{},
		Conn: c,
		Mutex: sync.Mutex{},
	}
}