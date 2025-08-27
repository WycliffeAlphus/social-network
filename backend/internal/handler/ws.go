package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"backend/internal/model"

	"github.com/gorilla/websocket"
)

var (
	users     = make(map[string]*websocket.Conn)
	broadcast = make(chan model.Message)
	mutex     = &sync.Mutex{}
)

// upgrader upgrades http conns to websocket conns
var upgrader = websocket.Upgrader{
	ReadBufferSize:  3000,
	WriteBufferSize: 3000,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000" // only allow local development origin
	},
}

// handleWebSocket uses the upgrader to upgrade the http conn
// then reads and writes messages from and to the ws
func WebSocketConnection(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hit")
		// upgrade initial get request to a WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		defer conn.Close()

		fmt.Println("Client connected via WebSocket!")

		var loginMsg model.Message

		jsonErr := conn.ReadJSON(&loginMsg)
		if jsonErr != nil {
			log.Println("Login message read error:", err)
			return
		}
		id := loginMsg.From

		mutex.Lock()
		users[id] = conn
		mutex.Unlock()

		log.Println(id, "connected")

		// Clean up on disconnect
		defer func() {
			mutex.Lock()
			delete(users, id)
			mutex.Unlock()
			log.Println(id, "disconnected")
		}()

		for {
			// read message from user
			var msg model.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("Read error from user: ", err)
				break
			}
			fmt.Println("here is the message",msg)

			broadcast <- msg
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast


		mutex.Lock()
		conn, ok := users[msg.From]
		mutex.Unlock()
		if ok {
			fmt.Println("we good")
			err := conn.WriteJSON("hello thereeeeee")
			if err != nil {
				log.Println("Error sending typing notification to", msg.From+":", err)
				conn.Close()
				delete(users, msg.From)
			}
		}
	}
}
