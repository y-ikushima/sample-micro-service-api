package internal

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	systemsHandler "sample-micro-service-api/apps/backend/app-service/internal/handler/systems"
	"sample-micro-service-api/package-go/database"
	"sample-micro-service-api/package-go/logging"
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
	// Zap Logger middleware
	s.router.Use(s.zapLoggerMiddleware())

	// Recovery middleware with Zap
	s.router.Use(s.zapRecoveryMiddleware())

	// CORS middleware
	s.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
}

// zapLoggerMiddleware はzapを使用したGinロガーミドルウェア
func (s *Server) zapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// リクエスト処理
		c.Next()

		// レスポンス時間とログ出力
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// Cloud Logging標準のhttpRequestフィールドを使用
		httpReq := logging.HttpRequest{
			RequestMethod: method,
			RequestUrl:    path,
			Status:        statusCode,
			ResponseSize:  string(rune(bodySize)),
			RemoteIp:      clientIP,
			Latency:       latency.String(),
			Protocol:      c.Request.Proto,
			UserAgent:     c.Request.UserAgent(),
			Referer:       c.Request.Referer(),
		}

		// ログレベルをステータスコードに基づいて決定
		if statusCode >= 500 {
			logging.LogHttpRequest("HTTP Request", httpReq, 
				zap.String("level", "ERROR"),
			)
		} else if statusCode >= 400 {
			logging.LogHttpRequest("HTTP Request", httpReq,
				zap.String("level", "WARNING"),
			)
		} else {
			logging.LogHttpRequest("HTTP Request", httpReq)
		}
	}
}

// zapRecoveryMiddleware はzapを使用したGinリカバリーミドルウェア
func (s *Server) zapRecoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, recovered interface{}) {
		logging.Error("Panic recovered",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("clientIP", c.ClientIP()),
			zap.Any("panic", recovered),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	})
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
	isConnected := s.dbClient.IsConnected()
	
	// ヘルスチェックログ
	logging.Debug("Health check requested",
		zap.Bool("database_connected", isConnected),
		zap.String("client_ip", c.ClientIP()),
	)

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "app-service",
		"timestamp": time.Now().Format(time.RFC3339),
		"database":  isConnected,
	})
} 