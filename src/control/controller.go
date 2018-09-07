package control

import (
	"encoding/json"
	//"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dchote/robot-mower/src/control/drivers"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"

	"gobot.io/x/gobot"
	//"gobot.io/x/gobot/api"
	//"gobot.io/x/gobot/drivers/gpio"
	//"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const (
	publishInterval = 1000 * time.Millisecond
)

type MowerControllerStruct struct {
	wsClients    map[*wsClientStruct]bool
	wsRegister   chan *wsClientStruct
	wsUnregister chan *wsClientStruct
	wsCommands   chan []byte

	wsPublishTicker *time.Ticker

	robotPlatform *gobot.Robot

	mpuData *drivers.MPUData
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

	// initialize the hardware platform devices
	r := raspi.NewAdaptor()
	ina := drivers.NewINA219Driver(r)
	mpu := drivers.NewMPU9250Driver(r)

	robotWork := func() {
		// we will want to sample our IMU at ~8hz (125ms) ALL i2c devices need to be read in here
		gobot.Every(1000*time.Millisecond, func() {
			// read voltage and current
			val, err := ina.GetBusVoltage()
			if err == nil {
				MowerState.Battery.Voltage = val
			}
			val, err = ina.GetCurrent()
			if err == nil {
				MowerState.Battery.Current = val
			}

			// read IMU data
			err = mpu.GetData()
			if err == nil {
				SetIMUValues(mpu.Data)
			}

		})
	}

	InitMowerState()
	InitFilters()
	UpdateSystemState()

	// mower controller
	MowerController = &MowerControllerStruct{
		wsClients:    make(map[*wsClientStruct]bool),
		wsRegister:   make(chan *wsClientStruct),
		wsUnregister: make(chan *wsClientStruct),
		wsCommands:   make(chan []byte),

		wsPublishTicker: time.NewTicker(publishInterval),

		robotPlatform: gobot.NewRobot("Mower",
			[]gobot.Connection{r},
			[]gobot.Device{ina, mpu},
			robotWork),
	}

	time.Sleep(1 * time.Second)

	// start the robotPlatform loop
	go MowerController.robotPlatform.Start()

	// websocket client transmit/recieve loop
	go MowerController.wsClientLoop()

	// websocket client data publishing loop, this will live elsewhere
	go wsPublishLoop()
}

func StopController() {
	MowerController.robotPlatform.Stop()
}

func InitMowerState() {
	sysInfo, _ := host.Info()
	MowerState.Platform.Hostname = sysInfo.Hostname
	MowerState.Platform.OperatingSystem = sysInfo.OS
	MowerState.Platform.Platform = sysInfo.Platform

	MowerState.Battery.Status = "Unknown"
	MowerState.Battery.VoltageNominal = 24.3
	MowerState.Battery.VoltageWarn = 23.0
	MowerState.Battery.Voltage = 23.5
	MowerState.Battery.Current = 0.1

	MowerState.Compass.Status = "Unknown"
	MowerState.Compass.Bearing = "NE"

	MowerState.GPS.Status = "Unknown"
	MowerState.GPS.Coordinates = "40.780715, -78.007729"

	MowerState.Drive.Speed = 100
	MowerState.Drive.Direction = "stopped"

	MowerState.Cutter.Speed = 0
}

func UpdateSystemState() {
	MowerState.Platform.CPULoad.Count, _ = cpu.Counts(false)

	cpuLoad, _ := cpu.Percent(publishInterval, false)
	MowerState.Platform.CPULoad.Total = cpuLoad[0]

	perCPU, _ := cpu.Percent(publishInterval, true)
	// TODO do this better, this is a hax but I dont know the right way to do it right now.
	MowerState.Platform.CPULoad.Core1 = perCPU[0]
	if MowerState.Platform.CPULoad.Count >= 2 {
		MowerState.Platform.CPULoad.Core2 = perCPU[1]
	}
	if MowerState.Platform.CPULoad.Count >= 3 {
		MowerState.Platform.CPULoad.Core3 = perCPU[2]
	}
	if MowerState.Platform.CPULoad.Count >= 4 {
		MowerState.Platform.CPULoad.Core4 = perCPU[3]
	}
	if MowerState.Platform.CPULoad.Count >= 5 {
		MowerState.Platform.CPULoad.Core5 = perCPU[4]
	}
	if MowerState.Platform.CPULoad.Count >= 6 {
		MowerState.Platform.CPULoad.Core6 = perCPU[5]
	}
	if MowerState.Platform.CPULoad.Count >= 7 {
		MowerState.Platform.CPULoad.Core7 = perCPU[6]
	}
	if MowerState.Platform.CPULoad.Count >= 8 {
		MowerState.Platform.CPULoad.Core8 = perCPU[7]
	}

	loadInfo, _ := load.Avg()
	MowerState.Platform.LoadAverage.Load1 = loadInfo.Load1
	MowerState.Platform.LoadAverage.Load5 = loadInfo.Load5
	MowerState.Platform.LoadAverage.Load15 = loadInfo.Load15

	memInfo, _ := mem.VirtualMemory()
	MowerState.Platform.MemoryUsage.Total = memInfo.Total
	MowerState.Platform.MemoryUsage.Available = memInfo.Available

	diskInfo, _ := disk.Usage("/")
	MowerState.Platform.DiskUsage.Total = diskInfo.Total
	MowerState.Platform.DiskUsage.Free = diskInfo.Free
}

func wsPublishLoop() {
	for {
		select {
		case <-MowerController.wsPublishTicker.C:
			UpdateSystemState()
			wsPublishState()
		}
	}
}

func wsPublishState() {
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

func (m *MowerControllerStruct) wsClientLoop() {
	for {
		select {
		case client := <-m.wsRegister:
			m.wsClients[client] = true
		case client := <-m.wsUnregister:
			if _, ok := m.wsClients[client]; ok {
				delete(m.wsClients, client)
				close(client.send)
			}
		case message := <-m.wsCommands:
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
				go wsPublishState()
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
	client.controller.wsRegister <- client

	go client.writeWebSocket()
	go client.readWebSocket()

	return nil
}

func (c *wsClientStruct) readWebSocket() {
	defer func() {
		c.controller.wsUnregister <- c
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

		c.controller.wsCommands <- message
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
