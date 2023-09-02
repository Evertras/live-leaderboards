package server

import (
	"net/http"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Repo interface {
	RoundRepo
}

type Server struct {
	e *echo.Echo
	r Repo
}

func New(r Repo) *Server {
	e := echo.New()

	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())

	// CORS is handled mostly through API Gateway layer, we don't need to
	// worry about it... we just need to accept whatever we're given
	e.Use(middleware.CORS())

	s := &Server{
		e: e,
		r: r,
	}

	api.RegisterHandlers(e, s)

	return s
}

func (s *Server) ListenAndServe(addr string) error {
	return s.e.Start(addr)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.e.ServeHTTP(w, r)
}

func (s *Server) WithLogLevel(level log.Lvl) *Server {
	s.e.Logger.SetLevel(level)
	return s
}
