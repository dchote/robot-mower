package control

import (
	//"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	publishInterval = 2
)

type MowerControllerStruct struct {
	wsClients  map[*wsClientStruct]bool
	register   chan *wsClientStruct
	unregister chan *wsClientStruct
	commands   chan []byte

	publishTicker *time.Ticker
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
	MowerState.Drive.Direction = "stopped"

	MowerState.Cutter.Speed = 45

	// mower controller
	MowerController = &MowerControllerStruct{
		wsClients:     make(map[*wsClientStruct]bool),
		register:      make(chan *wsClientStruct),
		unregister:    make(chan *wsClientStruct),
		commands:      make(chan []byte),
		publishTicker: time.NewTicker(publishInterval * time.Second),
	}

	go MowerController.run()

	go func() {
		for {
			select {
			case <-MowerController.publishTicker.C:
				PublishState()
			}
		}
	}()
}

func StopController() {

}

func PublishState() {
	//log.Println("publishing state")

	message, _ := json.Marshal(StateMessage{MowerStateStruct: MowerState, Namespace: "mower", Mutation: "setMowerState"})

	for client := range MowerController.wsClients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(MowerController.wsClients, client)
		}
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
			} else {
				// we want to stay in this processing loop, so never return out

				if strings.Compare(commandMessage.Method, "setMowerDriveSpeed") == 0 {
					MowerState.Drive.Speed, _ = strconv.Atoi(commandMessage.Value)
				} else if strings.Compare(commandMessage.Method, "setMowerCutterSpeed") == 0 {
					MowerState.Cutter.Speed, _ = strconv.Atoi(commandMessage.Value)
				} else if strings.Compare(commandMessage.Method, "requestDirectionStart") == 0 {
					// TODO actual callout logic, right now we'll just update state
					MowerState.Drive.Direction = commandMessage.Value
				} else if strings.Compare(commandMessage.Method, "requestDirectionStop") == 0 {
					// TODO actual callout logic, right now we'll just update state
					MowerState.Drive.Direction = "stopped"
				}

				// send updated state immediately
				go PublishState()
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
		}
	}
}
