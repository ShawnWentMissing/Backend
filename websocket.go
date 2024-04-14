package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by default for testing purposes
		return true
	},
}

type WebSocketHandler struct {
	// Player1Score int
	// Player2Score int
	// Handout      bool
}

var raspberryPiConn *websocket.Conn

func (wsh *WebSocketHandler) handleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	raspberryPiConn = conn

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from WebSocket:", err)
			break
		}
	}
}

func sendDataToRaspberryPi(data string) {
	if raspberryPiConn == nil {
		log.Println("Connection to Raspberry Pi is not established")
		return
	}

	err := raspberryPiConn.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil {
		log.Println("Error sending data to Raspberry Pi:", err)
		return
	}
}

func startWebSocketServer() {
	wsh := &WebSocketHandler{}

	http.HandleFunc("/ws", wsh.handleWebSocketConnections)

	log.Println("Starting WebSocket server on :8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting WebSocket server:", err)
	}
}

func startWebSocketClientToRaspberryPi() {
	// for testing purposes
	var err error
	raspberryPiConn, _, err = websocket.DefaultDialer.Dial("ws://raspberrypi.local:8765", nil)
	if err != nil {
		log.Fatal("Error connecting to Raspberry Pi:", err)
	}
	defer raspberryPiConn.Close()

	err = raspberryPiConn.WriteMessage(websocket.TextMessage, []byte("Hello from Go!"))
	if err != nil {
		log.Println("Error sending message to Raspberry Pi:", err)
	}
}
