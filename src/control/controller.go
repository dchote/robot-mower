package control

import (
	"net/http"

	"github.com/labstack/echo"
)

func StartController() {

}

func StopController() {

}

func WebSocketConnection() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}
