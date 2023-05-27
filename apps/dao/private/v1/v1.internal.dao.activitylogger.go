package v1

import (
	"net/http"
	"pulse-service/apps/repository/adapter"
	"pulse-service/apps/svc"

	model_event_logger "pulse-service/models/model.eventlogger"

	"github.com/gin-gonic/gin"
)

type EventLoggerDao struct {
	ServiceRepo *svc.ServiceRepository
}

type IEventLoggerDao interface {
	LogActivity(ctx *gin.Context, payload model_event_logger.Event)
}

func NewEventLoggerDao(repository *adapter.Repository) IEventLoggerDao {
	return &EventLoggerDao{
		ServiceRepo: svc.NewServiceRepo(repository),
	}
}

func (dao EventLoggerDao) LogActivity(ctx *gin.Context, payload model_event_logger.Event) {
	go func(payload model_event_logger.Event) {
		_ = dao.ServiceRepo.EventLoggerService.LogEvent(payload)
	}(payload)
	ctx.JSON(http.StatusOK, map[string]interface{}{"message": payload})
	ctx.Abort()
}
