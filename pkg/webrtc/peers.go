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
	p.ListLock.Lock(){
		p.ListLock.Unlock()
		p.DispatchKeyFrame()
	}()

	attemptSync := func() (tryAgain bool){
		for i := range p.Connections{
			if p.Connections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed{
				p.Connections - append(p.Connections[:i], p.Connections[i+1]...)
				log.Println("a". p.Connections)
				return true
			}

			existingSenders := map[string]bool{}
			for _, sender := range p.Connections[i].PeerConnection.GetSenders(){
				if sender.Track() == nil {
					continue
				}

				existingSenders[senders.Track().ID()] = true

				if _, ok := p.TrackLocals[sender.Track().ID()]; !ok {
					if err := p.Connections[i].PeerConnection.RemoveTrack(Sender); err != nil {
						return true
					}
				}
			}

			for _, receiver := range p.Connections[i].PeerConnection.GetReceivers(){
				if receiver.Track() == nil {
					continue 
				}
				existingSenders[receiver.Track().ID()] = true
			}

			for trackID := range p.TrackLocals{
				if _, ok := existingSenders[trackID]; !ok{
					if _, err := p.Connections[i].PeerConnection.AddTrack(p.TrackLocals[trackID]); err!= nil{
						return true
					}
				}
			}
		}
	}
}

func (p *Peers) DispatchKeyFrame() {

}

type websocketMessage struct{
	Event string `json:"event"`
	Data string `json:"data"`
}