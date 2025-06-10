package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetIp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"client_ip":       c.ClientIP(),
		"x-real-ip":       c.Request.Header.Get("X-Real-IP"),
		"x-forwarded-for": c.Request.Header.Get("X-Forwarded-For"),
	})
}

func SetupRoutes(r *gin.Engine, path string) {
	api := r.Group(path)

	ipRoutes := api.Group(
		"/ip",
	)
	{
		ipRoutes.GET("/", GetIp)
	}
}

func main() {
	r := gin.Default()
	r.TrustedPlatform = "X-Real-IP"
	SetupRoutes(r, "api")
	r.Run("0.0.0.0:" + "8080")
}
