package server

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	http.Server
}

const port = 5000

func NewServer(e *echo.Echo) *Server {
	return &Server{
		http.Server{
			Addr:              ":" + strconv.Itoa(port),
			Handler:           e,
			ReadTimeout:       30 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
			WriteTimeout:      30 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	log.Println("start serving in :" + strconv.Itoa(port))
	return s.ListenAndServe()
}
