package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/techievee/opensimsim/config"
)

type Server struct {
	E   *echo.Echo
	cfg *config.Server
}

func NewServerModule(cfg *config.Server) *Server {
	e := echo.New()
	e.Use(middleware.BodyLimit(cfg.Limit))
	if cfg.LogLevel == "debug" {
		e.Logger.SetLevel(log.DEBUG)
	} else if cfg.LogLevel == "info" {
		e.Logger.SetLevel(log.INFO)
	} else {
		e.Logger.SetLevel(log.ERROR)
	}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType},
		AllowMethods: []string{echo.GET},
	}))

	return &Server{E: e, cfg: cfg}
}

func (s *Server) StartServer() {
	s.E.Logger.Fatal(s.E.Start(s.cfg.Listen))
}
