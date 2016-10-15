package messago

//server

type Server struct {
	Clients map[*Client]bool

	Publish    chan *Message

	Register   chan *Client

	Unregister chan *Client

	Messages   map[int]*Message
}


func (this *Server) Init(){
	//starting Server
	go this.Run()
}

func (this *Server) AddClient(Client *Client) {
	//wait until registeratuin
	this.Register <- Client
	Client.Run()
}


func (this *Server) Remove(Client *Client) {
	this.Unregister <- Client
}

func (this *Server) Run(){
	//Listener
	for {
		select {
		//message publishing into other clients
		case message := <-this.Publish:
			this.SendMessage(*message)
		//client adding to list of server client list
		case client := <-this.Register:
			this.Clients[client] = true
		//Client removing from client list and closing send channel.
		case client := <-this.Unregister:
			if _, ok := this.Clients[client]; ok {
				delete(this.Clients, client)
				close(client.Send)

			}
		}
	}
}

//Send message all clients in server
func(this *Server) SendMessage(message Message){
	for client := range this.Clients {
		client.Send <- message
	}
	if MessagoConfig.MessagePersistent == true {
		index := len(this.Messages)
		this.Messages[index] = &message
	}

}

