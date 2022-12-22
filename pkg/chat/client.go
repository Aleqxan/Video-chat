package chat

import (
	"time"
	"bytes"
	"log"
	
	"github.com/fasthttp/websocket"
)

const(
	writeWait 		= 10 * time.Second 
	pongWait 		= 60 * time.Second 
	pingPeriod 		=(pongWait * 9) / 10
	maxMessageSize 	= 512
)

var (
	newline = []byte{'\n'}
	space = []byte{' '}
)

type Client struct {
	Hub *Hub 
	Conn *websocket.Conn
	Send chan []byte
}

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func (c *Client) readPump(){
	defer func(){
		c.Hub.unregister <- c 
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit (maxMessageSize)
	c.Conn.setReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {c.Conn.setReadDeadline(time.Now().Add(pongWait)); return nil})
	for {
		_, message, err := c.Conn.ReadMessage()
		if  err != nil {
			if websocket.IsUnexpectedClosedError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(messages, newline, space, -1))
		c.Hub.broadcast <- message
	}
}

func (c *Client) writePump(){
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select{
		case message, ok := <-c.Send:
			c.Conn.SetWriterDeadline(time.Now().Add(writeWait))
			if !ok{
				return
			}
			c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(c.Send)

			for i:=0; i<n; i++{
				w.Write(newline)
				w.Write(<-c.Send)
			}
			if err := w.Close(); err != nil {
				return 
			}
		case <- ticker.C:
			c.Conn.SetWriterDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(Websocket.PingMessage, nil); err != nil {
				return 
			}	
		}
	}
}

func PeerChatConn(c *websocket.Conn, hub *Hub){
	client := &Client{Hub: hub, Conn: c, Send: make(chan []byte, 256)}
	clinet.Hu .register <- client

	go client.writePump()
	client.readPump()
}