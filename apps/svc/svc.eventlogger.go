package svc

import (
	adapter "pulse-service/apps/repository/adapter"
	"pulse-service/constants"
	"pulse-service/logger"
	model_event_logger "pulse-service/models/model.eventlogger"

	"go.uber.org/zap"
)

type EventLoggerService struct {
	repository *adapter.Repository
}

func (svc EventLoggerService) LogEvent(payload model_event_logger.Event) error {
	logger.EventLogger.Info(constants.SERVICE_NAME, zap.Any("payload", payload))
	return nil
}
