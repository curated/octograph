package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/curated/octograph/graph"
	"github.com/labstack/echo"
)

// New creates the server
func New() *Server {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	s := &Server{
		Echo:   echo.New(),
		Logger: logger,
		Graph:  graph.New(logger),
	}
	s.Echo.GET("/fetch", s.process)
	return s
}

// Server serves HTTP requests
type Server struct {
	Echo   *echo.Echo
	Logger *log.Logger
	Graph  *graph.Graph
}

// Start initializes the server
func (s *Server) Start() {
	s.Logger.Fatal(s.Echo.Start(":1323"))
}

func (s *Server) process(c echo.Context) error {
	i, err := s.Graph.Fetch()
	s.Logger.Printf("Fetched body of length: %d, err: %v", i, err)
	return c.String(http.StatusOK, fmt.Sprintf("len=%d err=%v", i, err))
}
