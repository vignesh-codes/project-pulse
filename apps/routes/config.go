package routes

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	timeout time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Cors() func(c *gin.Context) {
	return func(c *gin.Context) {
		origin_header := c.Request.Header["Origin"]
		if len(origin_header) > 0 {
			c.Header("Access-Control-Allow-Origin", origin_header[0])
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accesstoken, Accept-language, Authorization, Content-Type, x-app-version,x-platform, x-client-id, x-client-secret")
		c.Header("Access-Control-Allow-Methods", "GET,HEAD,PUT,POST,PATCH,DELETE,OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (c *Config) SetTimeout(seconds int) *Config {
	c.timeout = time.Duration(seconds) * time.Second
	return c
}

func (c *Config) GetTimeout() time.Duration {
	return c.timeout
}
