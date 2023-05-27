package adapter

import (
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
)

type Repository struct {
	RedDB *RedDB
	PSql  *PSql
}

type IRedAdapter interface {
	Get(key string) ([]byte, error)
	Exists(key string) (int64, error)
	Set(key string, value []byte, expiry int) error
	HLen(key string) (int64, error)
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HSet(key, field, value string, expiry int) error
	Del(key string) error
	XRevRangeN(key, start, stop string, count int64) ([]redis.XMessage, error)
	ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) ([]redis.Z, error)
}

type RedDB struct {
	connection *redis.Client
}

type PSql struct {
	connection *gorm.DB
}

type IPSqlQueryAdapter interface {
	Select(key string) QueryBuilder
	RawQuery(queryString string) map[string]interface{}
	Exec(queryString string)
}

func RepositoryAdapter(redis *redis.Client, psqlClient *gorm.DB) *Repository {
	return &Repository{
		&RedDB{connection: redis},
		&PSql{connection: psqlClient},
	}
}
