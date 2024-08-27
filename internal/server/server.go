package server

import (
	"fmt"
	"net/http"
	"time"
	"weatherApp/internal/clients"
)

type Server struct {
	nc   clients.NWSDataService
	port int
}

func NewServer(port int, url string) *http.Server {
	NewSvr := &Server{
		port: port,
		nc:   clients.NewNWSDataService(url),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewSvr.port),
		Handler:      NewSvr.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
