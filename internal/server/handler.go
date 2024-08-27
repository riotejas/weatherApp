package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log/slog"
	"net/http"
	"os"
	wamiddleware "weatherApp/internal/middleware"
)

func handleError(w http.ResponseWriter, r *http.Request, err error, status int) {
	http.Error(w, err.Error(), status)
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(wamiddleware.QueryParamsToContext)

	r.Route("/v1", func(r chi.Router) {
		// GET forecast w/ validation
		r.With(wamiddleware.ValidateQueryParams(wamiddleware.ForecastSearchRules)).Get("/forecast", s.v1ForecastHandler)

	})

	// app health
	r.Get("/health", s.HealthHandler)

	// allow user to simulate an error response
	r.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(400), 400)
	})

	r.Get("/doc", s.DocHandler)

	return r
}

func (s *Server) v1ForecastHandler(w http.ResponseWriter, r *http.Request) {
	body, err := s.nc.Forecast(r.Context())
	if err != nil {
		slog.Error("ForecastHandler", "error", err)
		handleError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Checking health")
	// todo: what can we check for health wise?

	res := map[string]string{
		"message": "App is healthy",
	}
	jsonResp, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResp)
}

func (s *Server) DocHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./weather_app.yaml")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Disposition", "attachment; filename=weather_app.yaml")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "File copy error", http.StatusInternalServerError)
		return
	}
}
