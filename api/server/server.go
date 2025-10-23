package server

import (
	"api/config"
	"api/models"
	"api/routes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Router *chi.Mux
	Config *models.ServerConfig
}

func NewServer(config models.ServerConfig) *Server {
	if config.DB == nil {
		log.Fatal("DB is nil")
	}
	if config.Cache == nil {
		log.Fatal("Cache is nil")
	}
	if config.Port == "" {
		log.Fatal("Port is empty")
	}

	return &Server{
		Router: chi.NewRouter(),
		Config: &config,
	}
}

func (s *Server) Start() error {
	s.Setup()

	log.Println("Server started on port:", s.Config.Port)
	return http.ListenAndServe(":"+s.Config.Port, s.Router)
}

func (s *Server) Setup() *Server {
	s.setupMiddleware()
	s.setupRoutes()
	return s
}

// setupMiddleware configures the middleware
func (s *Server) setupMiddleware() *Server {
	// Basic middleware
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Compress(5))
	s.Router.Use(middleware.CleanPath)
	s.Router.Use(middleware.RedirectSlashes)

	// Logger middleware only in dev
	if config.Server().Env == "dev" || config.Server().Env == "development" {
		s.Router.Use(middleware.Logger)
	}

	// Timeout
	s.Router.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.Server().ClientHost},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	return s
}

// setupRoutes configures the routes
func (s *Server) setupRoutes() *Server {
	// Global routes
	s.Router.Get("/", s.handleRoot)
	s.Router.Get("/health", s.handleHealth)

	// API routes
	routes.SetupRoutes(s.Router, s.Config)

	// 404 handler
	s.Router.NotFound(s.handleNotFound)

	return s
}

// Handler functions
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "server running",
	})
}

// handleHealth checks the health of the server
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	var db *sql.DB
	var err error

	if db, err = s.Config.DB.DB(); err != nil {
		http.Error(w, "DB unreachable", http.StatusInternalServerError)
		return
	}

	if err := db.Ping(); err != nil {
		http.Error(w, "DB unreachable", http.StatusInternalServerError)
		return
	}

	if err := s.Config.Cache.Ping(r.Context()); err != nil {
		http.Error(w, "Redis unreachable", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusNotFound, map[string]string{
		"error": "route not found",
	})
}

// Helper function to send JSON responses
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
