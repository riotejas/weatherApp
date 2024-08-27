package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
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

	// todo: add openapi doc for endpoints
	r.Get("/doc", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("openAPI doc is TBD"))
	})

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
