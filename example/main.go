package main

import (
	"messago"
	"net/http"
	"flag"
	"text/template"
	"log"
	"fmt"
	"github.com/gorilla/websocket"
)
//This is not good way only for demo !
var roomname = ""
var username = ""

var addr = flag.String("addr",":8080","http service address")
var home = template.Must(template.ParseFiles("home.html"))
var room = template.Must(template.ParseFiles("room.html"))


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	home.Execute(w, r.Host)
}

func serveRoom(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		fmt.Errorf("%v",err)
	}
	//Take form values and put global variable
	//this method is not good way. Only for demo !
	user := r.FormValue("username")
	roomnamepost :=  r.FormValue("roomname")
	roomname = roomnamepost
	username = user
	room.Execute(w, r.Host)
}

func serveChat(w http.ResponseWriter, r *http.Request)  {
	//this code create server but if server has created before it gives created server
	server := messago.MessagoConfig.NewServer(roomname)
	//websocket upgrader
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//creating client
	client := messago.ClientInit(conn,roomname,username)
	//adding client
	server.AddClient(client)


}

func main() {
	//this is messago initialization
	//It must work at least one time
	//before create room-server and client you have to initialize config
	//can give roomPoolSize in config it means each server can handle that number of client
	//0 means unlimited client and default value is 0 you don't have to put.
	//for example : messago.Config{MessagePersistent:true,RoomPoolSize:0}
	messago.Init(messago.Config{MessagePersistent:false})
	//messago.MessagoConfig is a global variable you can create new server with this variable
	// for example : messago.MessagoConfig.NewServer("serverName")


	//these are serving static files
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/chat", serveRoom)

	//websocket is working here !
	http.HandleFunc("/ws/chat",serveChat)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

