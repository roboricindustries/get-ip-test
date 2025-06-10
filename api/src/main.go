package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func InitializeLogger() *logrus.Logger {
	logFile, err := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalf("Failed to create log file: %v", err)
	}

	logger := logrus.New()
	logger.SetOutput(logFile)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		if l, exists := c.Get("logger"); exists {
			log := l.(*logrus.Logger)
			log.WithFields(logrus.Fields{
				"latency":   duration.Seconds(),
				"client_ip": c.ClientIP(),
				"method":    c.Request.Method,
				"path":      c.Request.URL.Path,
				"status":    c.Writer.Status(),
			}).Info("[GIN] request handled")
		}
	}
}

func InjectLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	}
}

func GetIp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"client_ip":       c.ClientIP(),
		"x-real-ip":       c.Request.Header.Get("X-Real-IP"),
		"x-forwarded-for": c.Request.Header.Get("X-Forwarded-For"),
	})
}

func SetupRoutes(r *gin.Engine, path string) {
	api := r.Group(path)

	rate, err := limiter.NewRateFromFormatted("5-M")
	if err != nil {
		panic(err)
	}

	// In-memory store (for production, use Redis)
	store := memory.NewStore()

	// Create limiter instance
	limiterInstance := limiter.New(
		store, rate,
		limiter.WithClientIPHeader("X-Real-IP"),
	)

	// Apply middleware
	rateMiddleware := ginlimiter.NewMiddleware(limiterInstance)

	ipRoutes := api.Group(
		"/ip",
	)
	ipRoutes.Use(rateMiddleware)
	{
		ipRoutes.GET("/", GetIp)
	}
}

func main() {
	logger := InitializeLogger()
	r := gin.Default()
	r.Use(InjectLogger(logger))
	r.Use(LoggerMiddleware())
	r.TrustedPlatform = "X-Real-IP"
	SetupRoutes(r, "api")
	r.Run("0.0.0.0:" + "8080")
}
