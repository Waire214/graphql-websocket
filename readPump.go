package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// In the readPump Goroutine, the client will read new messages send over the
// WebSocket connection. It will do so in an endless loop until the client is
// disconnected. When the connection is closed, the client will call its own
// disconnect method to clean up.

const (

	// max wait time when writing message to peer
	writeWait = 10 * time.Second

	// max time till next pong from peer
	pongWait = 60 * time.Second

	// send ping interval, must be less than pong wait time
	pingPeriod = (pongWait * 9) / 10

	// maximum message size allowed from peer
	maxMessageSize = 1000
)

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	close(client.send)
	client.conn.Close()
}

func (client *Client) readPump() {

	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	//start endless read loop, waiting for messages from client
	for {

		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.wsServer.broadcast <- jsonMessage
	}
}

// Upon receiving new messages the client will push them in the WsServer
// broadcast channel.

