package messago

import "errors"

var RoomDoesntExist = errors.New("Room Doesn't exist")
var RoomAlreadyCreated = errors.New("Room already created !")
var serverList map[string]*Server
var MessagoConfig Config

type Config struct {
	MessagePersistent bool
	//Zero means unlimited Room pool.
	RoomPoolSize      int
}

//Default unlimited Room pool.
func Init(Config Config) {
	MessagoConfig = Config
	//Default make Unlimited pool
	serverList = make(map[string]*Server, MessagoConfig.RoomPoolSize)
}

func GetServer(name string) (*Server, error) {
	if server, ok := serverList[name]; ok {
		return server, nil
	}
	return &Server{}, RoomDoesntExist
}

func (this *Config) NewServer(Name string) *Server {
	server,err := GetServer(Name)
	if err == nil {
		return server
	}
	serverList[Name] = &Server{
		Publish:make(chan *Message),
		Register:make(chan *Client),
		Unregister:make(chan *Client),
		Clients:make(map[*Client]bool),
		Messages:make(map[int]*Message),
	}
	//running server
	go serverList[Name].Run()

	return serverList[Name]

}
