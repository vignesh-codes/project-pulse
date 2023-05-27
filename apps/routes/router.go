package routes

import (
	"io"
	"log"
	"net/http"
	"os"
	"pulse-service/apps/repository/adapter"
	"pulse-service/apps/routes/client"
	"pulse-service/apps/routes/internal"
	"pulse-service/constants"
	"pulse-service/middlewares"
	"pulse-service/utils/response"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	config *Config
	router *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		config: NewConfig().SetTimeout(30),
		router: gin.New(),
	}
}

func healthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		health := struct {
			Version   string `json:"version"`
			Env       string `json:"env"`
			Service   string `json:"service"`
			Message   string `json:"message"`
			Timestamp int64  `json:"timestamp"`
		}{
			Version:   "1.0.1",
			Env:       constants.ENV,
			Service:   "pulse-service",
			Message:   "ok",
			Timestamp: time.Now().UnixMilli(),
		}
		c.JSON(200, health)
	}
}

func invalidPathHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := response.RouteNotFound()
		c.JSON(status.Status(), status)
		c.Abort()
	}
}

func (r *Router) setConfigRouters() {
	r.SetMode()
	r.EnableAPILogger()
	r.EnableCORS()
	r.EnableRecover()
	r.RouterHealth()
}

func (r *Router) SetMode() *Router {
	gin.SetMode(gin.ReleaseMode)
	return r
}

func (r *Router) SetRouters(repository *adapter.Repository) http.Handler {
	r.setConfigRouters()
	// r.SetClientRoutes(repository)
	r.SetInternalRoutes(repository)
	r.router.NoRoute(invalidPathHandler())
	return r.router
}

func (r *Router) RouterHealth() {
	r.router.GET("/health", healthHandler())
}

func (r *Router) EnableAPILogger() *Router {
	accessLogFile, fileErr := os.Create("logs/access.log")
	if fileErr != nil {
		log.Fatal(fileErr)
		return r
	}
	gin.DefaultWriter = io.MultiWriter(accessLogFile)
	errorLogFile, fileErr := os.Create("logs/error.log")
	if fileErr != nil {
		log.Fatal(fileErr)
		return r
	}
	gin.DefaultErrorWriter = io.MultiWriter(errorLogFile)
	r.router.Use(middlewares.GetRequestAndResponseLog())
	return r
}

func (r *Router) EnableCORS() *Router {
	r.router.Use(r.config.Cors())
	return r
}

func (r *Router) EnableRecover() *Router {
	r.router.Use(gin.Recovery())
	return r
}

func (r *Router) SetInternalRoutes(repository *adapter.Repository) {
	v1Group := r.router.Group("v1/internal/")
	internal.V1(v1Group, repository)
}

func (r *Router) SetClientRoutes(repository *adapter.Repository) {
	v1Group := r.router.Group("v1/")
	// v1PublicGroup := r.router.Group("v1/")
	client.V1(v1Group, repository)
	// client.V1PublicApis(v1PublicGroup, repository)
}
