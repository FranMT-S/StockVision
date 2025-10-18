package server

import (
	"api/config"
	"api/models"
	"api/routes"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
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
		Router: gin.Default(),
		Config: &config,
	}
}

func (s *Server) Start() error {
	s.Setup()

	log.Println("Server started on port:", s.Config.Port)
	return s.Router.Run(":" + s.Config.Port)
}

func (s *Server) Setup() *Server {
	s.configCORS()
	s.setupMiddleware(s.Router)
	s.setupRoutes()
	return s
}

// setupMiddleware configures the middleware
func (s *Server) setupMiddleware(router *gin.Engine) *Server {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(gin.Recovery())
	// Logger middleware
	if config.Server().Env == "dev" || config.Server().Env == "development" {
		router.Use(gin.Logger())
	}

	return s
}

// setupRoutes configures the routes
func (s *Server) setupRoutes() *Server {

	s.setupGlobalRoutes()
	routes.SetupRoutes(s.Router, s.Config)
	return s
}

// configCORS configures the CORS
func (s *Server) configCORS() *Server {
	s.Router.Use(cors.New(cors.Config{
		AllowOrigins:     config.Server().ClientHost,
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,

		MaxAge: 6 * time.Hour,
	}))

	return s
}

func (s *Server) setupGlobalRoutes() *Server {
	// Define a simple GET endpoint
	s.Router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "server running",
		})
	})

	s.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	s.Router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })
	return s
}
