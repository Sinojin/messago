# Messago
Messago is Instant Message Library using websocket.
I am impressed from [gorilla-websocket](https://github.com/gorilla/websocket) chat examples.
You can Create many room not only one !


#Installation
First get library
```go
go get github.com/Sinojin/messago
```

At least one time you must run in your program
Also you can see more config information in example 
```
gomessago.Init(messago.Config{MessagePersistent:false})
```

Now we have global variable like this 
```go 
messago.MessagoConfig
```

To create room (server)
```go
messago.MessagoConfig.NewServer(roomname)
```

now we have config and room 
We have to add client to server
```go
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
``` 
Making client
```go
conn, err := upgrader.Upgrade(w, r, nil)
 	if err != nil {
 		log.Println(err)
 		return
 	}
 	//creating client
 	 client := messago.ClientInit(conn,roomname,username)
```
 
Finaly add client to server
 
```go
server.AddClient(client)
```

#Plans 
You can help me in this topics
- Redis Support for storage message
- Unit tests with ginkgo (But i don't know websocket connection in test)

 






