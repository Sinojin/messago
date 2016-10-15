package messago

import (
	"github.com/gorilla/websocket"
	"time"
	"fmt"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Client struct {
	Conn     *websocket.Conn // The websocket connection.

	Send     chan Message

	serverName   string

	Username string
}

func ClientInit(Con *websocket.Conn, ServerName string, Username string) *Client {
	return &Client{Conn:Con, serverName:ServerName, Send:make(chan Message), Username:Username}
}
func (this *Client) Run() {
	go this.ReadPump()
	this.WritePump()
}

func (client *Client) ReadPump() {
	defer func() {
		client.Conn.Close()
	}()
	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				fmt.Errorf("Reading Error : %v ", err)
			}
			break
		}
		msg := string(message)
		Message := MessageInit(client.Username, msg)
		//Wait until publish to all clients
		server, _ :=GetServer(client.serverName)
		server.Publish <- Message
	}

}

func (client *Client) WritePump() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		server,_ := GetServer(client.serverName)
		server.Remove(client)
		client.Conn.Close()

	}()
	//old message sending...
	client.sendPreviousMessageBack()

	for {
		select {
		case message := <-client.Send:
		//Sending Message to client
			client.Conn.WriteJSON(message)
		case <-ticker.C:
		//Checking Connection.
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return

			}
		}
	}

}

func (client *Client) sendPreviousMessageBack() {
	server ,_ := GetServer(client.serverName)
	index := len(server.Messages)
	for i := 0; i < index; i++ {
		message := server.Messages[i]
		client.Conn.WriteJSON(message)
	}

}