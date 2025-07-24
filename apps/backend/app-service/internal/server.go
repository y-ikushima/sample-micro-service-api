package internal

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	systemsHandler "sample-micro-service-api/apps/backend/app-service/internal/handler/systems"
	"sample-micro-service-api/package-go/database"
)

type Server struct {
	dbClient       *database.Client
	router         *gin.Engine
	systemsHandler *systemsHandler.Handler
}

func NewServer(dbClient *database.Client, systemsHandler *systemsHandler.Handler) *Server {
	// Set Gin mode from environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	server := &Server{
		dbClient:       dbClient,
		router:         gin.New(),
		systemsHandler: systemsHandler,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddleware() {
	// Logger middleware
	s.router.Use(gin.Logger())

	// Recovery middleware
	s.router.Use(gin.Recovery())

	// CORS middleware
	s.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Systems endpoints
		v1.GET("/systems", s.systemsHandler.GetSystems)
		v1.POST("/systems", s.systemsHandler.CreateSystem)
		v1.GET("/systems/:id", s.systemsHandler.GetSystemById)
		v1.PUT("/systems/:id", s.systemsHandler.UpdateSystem)
		v1.DELETE("/systems/:id", s.systemsHandler.DeleteSystem)
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "app-service",
		"timestamp": gin.H{"iso": gin.H{}},
		"database":  s.dbClient.IsConnected(),
	})
} 