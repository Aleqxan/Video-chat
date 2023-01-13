type Room struct {
	Peers *Peers
	Hub *chat.Hub
}

type Peers struct {
	ListLock sync.RWMutex
	Connections []PeerConnectionState
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type PeerConnectionState struct{
	PeerConnection *webrtc.PeerConnection
	websocket *ThreadsSafeWriter
}

type ThreadsSafeWriter struct{
	Conn *websocket.Conn 
	Mutex sync.Mutex
}

func (t *ThreadsSafeWriter) WriteJSON (v interace()) error{
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	return t.Conn.WriteJSON(v)
}

func (p *Peers) AddTrack (t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP{
	p.ListLock.Lock()
	defer func(){
		p.ListLock.Unlock()
		p.SignalPeerConnections()
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(). t.StreamID())

		if err != nil {
			log.Println(err.Error())
		}
		p.TrackLocals[t.ID()] = trackLocal
		return trackLocal
}

func (p *Peers)RemoveTrack (t *webrtc.TrackLocalStaticRTP){
	p.ListLock.Lock()
	defer func(){
		p.ListLock.Unlock()
		p.SignalPeerConnections()
	}()
	delete(p.TrackLocals, t.ID())
}

func (p *Peers)SignalPeerConnection() {

}

func (p *Peers) DispatchKeyFrame() {

}

type websocketMessage struct{
	Event string `json:"event"`
	Data string `json:"data"`
}