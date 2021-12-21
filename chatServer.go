package main

//sending and receiving messages
//This file contains a WsServer type that has one map for the Clients
//registered in the server. It also has two channels, one for register
//requests and one for unregister requests.

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast chan []byte
}

//NewWebSocketServer creates a new WsServer type
func NewWebSocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast: make(chan []byte),
	}
}

// Run our websocket server, accepting various request
func (server *WsServer) Run() {
	for {

		select {

			case client := <-server.register:
				server.registerClient(client)
				
			case client := <-server.unregister:
				server.unregisterClient(client)

			case message := <- server.broadcast:
				server.broadcastToClients(message)

		}
		
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}