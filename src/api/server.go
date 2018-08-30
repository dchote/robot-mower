package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dchote/robot-mower/src/api/handlers"
	"github.com/dchote/robot-mower/src/config"

	"github.com/dchote/robot-mower/src/control"
	"github.com/dchote/robot-mower/src/vision"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e *echo.Echo
)

func StartServer(cfg config.ConfigStruct, assets *rice.Box) {
	if e != nil {
		return
	}

	// instantiate echo instance
	e = echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 30 * time.Second

	// setup middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("32M"))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// prevent caching by client (e.g. Safari)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			return next(c)
		}
	})

	if assets != nil {
		assetHandler := http.FileServer(assets.HTTPBox())
		e.GET("/", echo.WrapHandler(assetHandler))
		e.GET("/favicon.ico", echo.WrapHandler(assetHandler))
		e.GET("/css/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/js/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/fonts/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/img/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	}

	// setup API routes
	e.GET("/v1/health", handlers.Health())
	e.GET("/v1/config", handlers.Config())
	e.GET("/v1/endpoints", handlers.Endpoints())

	e.GET("/camera", echo.WrapHandler(vision.Stream))
	e.GET("/ws", control.WebSocketConnection())

	log.Println("starting server on ", cfg.APIServer.ListenAddress)
	e.Start(cfg.APIServer.ListenAddress)
}

func StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
