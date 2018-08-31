package control

import (
	//"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type MowerControllerStruct struct {
	wsClients  map[*wsClientStruct]bool
	register   chan *wsClientStruct
	unregister chan *wsClientStruct
	commands   chan []byte
}

type wsClientStruct struct {
	controller *MowerControllerStruct
	conn       *websocket.Conn
	send       chan []byte
}

var (
	wsUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}

	MowerController *MowerControllerStruct
)

func StartController() {
	// build default state
	MowerState = new(MowerStateStruct)

	MowerState.Battery.Status = "Unknown"
	MowerState.Battery.Voltage = 24.1
	MowerState.Battery.Current = 0.1

	MowerState.Compass.Status = "Unknown"
	MowerState.Compass.Bearing = "NE"

	MowerState.GPS.Status = "Unknown"
	MowerState.GPS.Coordinates = "40.780715, -78.007729"

	MowerState.Drive.Speed = 100

	MowerState.Cutter.Speed = 45

	// mower controller
	MowerController = &MowerControllerStruct{
		wsClients:  make(map[*wsClientStruct]bool),
		register:   make(chan *wsClientStruct),
		unregister: make(chan *wsClientStruct),
		commands:   make(chan []byte),
	}

	go MowerController.run()
	go PublishState()
}

func StopController() {

}

func PublishState() {
	for {
		message, _ := json.Marshal(StateMessage{MowerStateStruct: MowerState, Namespace: "mower", Mutation: "setMowerState"})

		for client := range MowerController.wsClients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(MowerController.wsClients, client)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func (m *MowerControllerStruct) run() {
	for {
		select {
		case client := <-m.register:
			m.wsClients[client] = true
		case client := <-m.unregister:
			if _, ok := m.wsClients[client]; ok {
				delete(m.wsClients, client)
				close(client.send)
			}
		case message := <-m.commands:
			// handle the command message
			log.Println("command: " + string(message))

			var commandMessage CommandMessage
			err := json.Unmarshal(message, &commandMessage)
			if err != nil {
				log.Println("error decoding json")
				return
			}

			if strings.Compare(commandMessage.Method, "setMowerDriveSpeed") == 0 {
				MowerState.Drive.Speed = commandMessage.Value
			} else if strings.Compare(commandMessage.Method, "setMowerCutterSpeed") == 0 {
				MowerState.Cutter.Speed = commandMessage.Value
			}
		}
	}
}

func WebSocketConnection(c echo.Context) error {
	log.Println("WebSocket: " + c.RealIP() + " connected")

	conn, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &wsClientStruct{controller: MowerController, conn: conn, send: make(chan []byte, 256)}
	client.controller.register <- client

	go client.writeWebSocket()
	go client.readWebSocket()

	return nil
}

func (c *wsClientStruct) readWebSocket() {
	defer func() {
		c.controller.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.controller.commands <- message
	}
}

func (c *wsClientStruct) writeWebSocket() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				c.conn.WriteMessage(websocket.TextMessage, <-c.send)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
