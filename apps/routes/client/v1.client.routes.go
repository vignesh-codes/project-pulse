package client

import (
	"pulse-service/apps/repository/adapter"
	"pulse-service/logger"

	"github.com/gin-gonic/gin"
)

func V1(group *gin.RouterGroup, repository *adapter.Repository) {
	logger.ConsoleLogger.Debug("Initialising frontend v1 group routes.")

}
