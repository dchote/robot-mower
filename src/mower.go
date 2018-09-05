package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dchote/robot-mower/src/api"
	"github.com/dchote/robot-mower/src/config"
	"github.com/dchote/robot-mower/src/control"
	"github.com/dchote/robot-mower/src/vision"

	"github.com/GeertJohan/go.rice"
	"github.com/docopt/docopt-go"
)

const VERSION = "0.0.1"

var (
	err          error
	staticAssets *rice.Box
)

func cliArguments() {
	usage := `
Usage: mower [options]

Options:
  -c, --config=<json>           Specify config file [default: ./config.json]
	-d, --camera-device=<device>  Specify the devide id of the camera [default: 0]
  -h, --help                    Show this screen.
  -v, --version                 Show version.
`
	args, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	config.ConfigFile, _ = args.String("--config")

	err = config.LoadConfig(config.ConfigFile)
	if err != nil {
		log.Fatalf("Unable to load "+config.ConfigFile+" ERROR=", err)
	}

	config.Config.Mower.CameraDeviceID, _ = args.Int("--camera-device")

	log.Printf("Config: %+v", config.Config)
}

func exitCleanup() {

}

func main() {
	cliArguments()

	staticAssets, err = rice.FindBox("frontend/dist")
	if err != nil {
		log.Fatalf("Static assets not found. Build them first.")
	}

	vision.StartVision()
	defer vision.StopVision()

	control.StartController()

	go api.StartServer(*config.Config, staticAssets)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down")

	// shut down listener, with a hard timeout
	api.StopServer()
	control.StopController()

	// extra grace time
	time.Sleep(time.Second)

	os.Exit(0)
}
