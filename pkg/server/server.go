package server

import (
	"net/http"

	"github.com/Evertras/live-leaderboards/pkg/api"
	"github.com/labstack/echo/v4"
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
