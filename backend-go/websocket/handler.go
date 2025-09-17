// websocket.go - Comunicación en tiempo real
package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		switch msg.Type {
		case "mpi_simulation":
			handleMPISimulation(conn, msg.Payload)
		case "telemetry_update":
			handleTelemetryUpdate(conn, msg.Payload)
		case "code_execution":
			handleCodeExecution(conn, msg.Payload)
		}
	}
}

func handleMPISimulation(conn *websocket.Conn, payload interface{}) {
	// Simular comunicación MPI en tiempo real
	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		message := WSMessage{
			Type: "mpi_update",
			Payload: map[string]interface{}{
				"process": i,
				"data":    rand.Intn(100),
				"step":    i,
			},
		}
		conn.WriteJSON(message)
	}
}

// handleTelemetryUpdate simula la actualización de telemetría en tiempo real
func handleTelemetryUpdate(conn *websocket.Conn, payload interface{}) {
	// Simulación simple de actualización de telemetría
	telemetry := WSMessage{
		Type: "telemetry_update",
		Payload: map[string]interface{}{
			"speed":    rand.Intn(300),
			"rpm":      rand.Intn(15000),
			"position": rand.Intn(20),
		},
	}
	conn.WriteJSON(telemetry)
}

// handleCodeExecution simula la ejecución de código en tiempo real
func handleCodeExecution(conn *websocket.Conn, payload interface{}) {
	// Simulación simple de ejecución de código
	result := WSMessage{
		Type: "code_execution_result",
		Payload: map[string]interface{}{
			"output":  "Código ejecutado correctamente",
			"success": true,
		},
	}
	conn.WriteJSON(result)
}
