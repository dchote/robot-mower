package handlers

import (
	"net"
	"net/http"

	"github.com/dchote/robot-mower/src/config"

	"github.com/labstack/echo"
)

type JSONResponse map[string]interface{}

func Health() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, JSONResponse{
			"status": "OK",
		})
	}
}

func Config() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, config.Config)
	}
}

func Endpoints() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, JSONResponse{
			"camera": "http://" + GetLocalIP() + config.Config.APIServer.ListenAddress + "/camera",
			"ws":     "ws://" + GetLocalIP() + config.Config.APIServer.ListenAddress + "/ws",
		})
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
