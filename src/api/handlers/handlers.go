package handlers

import (
	//"encoding/json"
	//"fmt"
	//"log"
	"net/http"
	//"strings"

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
