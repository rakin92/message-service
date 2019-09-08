package config

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func setupCORSMiddleware() cors.Config {
	config := cors.DefaultConfig()
	allowedOrigin := viper.GetString("ALLOWED_ORIGIN")
	config.AllowOrigins = strings.Split(allowedOrigin, ",")
	config.AllowMethods = []string{"OPTIONS", "POST"}
	config.ExposeHeaders = []string{"Content-Type"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return config
}

// SetupRouter : setups routers
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(setupCORSMiddleware()))

	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(200, "up")
	})

	return r
}
