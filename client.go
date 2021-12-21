package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

//client represents the websocket client at the server
type Client struct {
	//the actual websocket connection
	conn *websocket.Conn
	wsServer *WsServer
	send chan []byte
}

func newClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		conn: conn,
		wsServer: wsServer,
		send: make(chan []byte, 256),
	}
}

//ServeWs handles websocket requests from clients requests
//The upgrader is used to upgrade the HTTP server connection to the WebSocket
//protocol. the return value of this function is a WebSocket connection.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer)

	go client.writePump()
	go client.readPump()

	wsServer.register <- client

	fmt.Println("New Client joined the hub!")
	fmt.Println(client)

}


