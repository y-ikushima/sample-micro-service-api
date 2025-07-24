package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"sample-micro-service-api/package-go/database"
)

type Server struct {
	dbClient *database.Client
	router   *gin.Engine
}

func NewServer(dbClient *database.Client) *Server {
	// Set Gin mode from environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	server := &Server{
		dbClient: dbClient,
		router:   gin.New(),
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
	s.router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "*"
		}

		// Check if origin is allowed
		if allowedOrigins == "*" || strings.Contains(allowedOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", os.Getenv("CORS_ALLOW_METHODS"))
		c.Header("Access-Control-Allow-Headers", os.Getenv("CORS_ALLOW_HEADERS"))
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

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