package main

import (
	//"io/ioutil"
	"log"
	"os"
	"os/signal"
	//"runtime"
	//"sync"
	"syscall"
	"time"

	"github.com/dchote/robot-mower/src/api"
	"github.com/dchote/robot-mower/src/config"

	"github.com/GeertJohan/go.rice"
	"github.com/docopt/docopt-go"
)

const VERSION = "0.0.COW💩"

var (
	cfg          *config.ConfigStruct
	err          error
	staticAssets *rice.Box
)

func cliArguments() {
	usage := `
Usage: mower [options]

Options:
  -c, --config=<json>      Specify config file [default: ./config.json]
  -h, --help               Show this screen.
  -v, --version            Show version.
`
	args, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	config.ConfigFile, _ = args.String("--config")
}

func exitCleanup() {

}

func main() {
	cliArguments()

	cfg, err = config.LoadConfig(config.ConfigFile)
	if err != nil {
		log.Fatalf("Unable to load "+config.ConfigFile+" ERROR=", err)
	}

	log.Printf("Loaded config: %+v", cfg)

	staticAssets, err = rice.FindBox("api/frontend/dist")
	if err != nil {
		log.Fatalf("Static assets not found. Build them first.")
	}

	go api.StartServer(*cfg, staticAssets)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down")

	// shut down listener, with a hard timeout
	api.StopServer()

	// extra grace time
	time.Sleep(time.Second)

	os.Exit(0)
}
