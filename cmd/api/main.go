package main

import (
	"fmt"
	"log/slog"
	cfg "weatherApp/internal/config"

	"weatherApp/internal/clients"
	"weatherApp/internal/server"
)

func main() {
	slog.Info("Loading config")
	config := cfg.NewConfig()
	err := config.LoadConfig(clients.Vendor)
	if err != nil {
		panic(err)
	}

	slog.Info("Starting server on", "port", config.Port)
	svr := server.NewServer(config.Port, config.Url)

	err = svr.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
