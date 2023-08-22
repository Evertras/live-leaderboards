package server

import (
	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Server struct {
	addr string
	e    *echo.Echo
}

func New(addr string) *Server {
	e := echo.New()

	e.Logger.SetLevel(log.INFO)

	s := &Server{
		addr: addr,
		e:    e,
	}

	api.RegisterHandlers(e, s)

	return s
}

func (s *Server) ListenAndServe() error {
	return s.e.Start(s.addr)
}
