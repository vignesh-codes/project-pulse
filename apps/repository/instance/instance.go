package instance

import (
	"context"
	"fmt"
	"pulse-service/constants"
	"pulse-service/logger"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetRedisConnection() *redis.Client {
	fmt.Println("setting redis ", constants.REDIS_SERVER)
	red := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 0,
	})
	var ctx = context.Background()
	err := red.Ping(ctx).Err()
	if err != nil {
		logger.ConsoleLogger.Fatal("GetRedisConnection", zap.Any(logger.KEY_ERROR, err.Error()))
		panic(err)
	}
	logger.ConsoleLogger.Info("Creating Redis Cluster Connection: ", zap.Any(logger.KEY_KEY, constants.REDIS_SERVER))
	return red
}

func GetPSqlConnection() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", constants.POSTGRESDB_HOST,
		constants.POSTGRESDB_USER, constants.POSTGRESDB_PWD, constants.POSTGRESDB_DB, constants.POSTGRESDB_PORT)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Creating PostgreSql connection")
	return client
}
