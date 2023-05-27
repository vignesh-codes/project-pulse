package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var count int = 0

type PostgreSQLDB struct {
	db *gorm.DB
}

var client *redis.Client

var db *PostgreSQLDB

func NewPostgreSQLDB() *PostgreSQLDB {
	// PostgreSQL connection details
	dsn := "host=localhost user=user password=pass dbname=postgres port=4432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(&Event{}); err != nil {
		log.Fatalf("Failed to perform auto-migration: %v", err)
	}
	fmt.Println("TABLE AVAILABLE ")
	return &PostgreSQLDB{
		db: db,
	}
}

func NewRedisStream(client *redis.Client) *RedisStream {
	return &RedisStream{
		client: client,
	}
}

// Connects to redis and psql and starts 30s tick to check redis
func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 0,
	})
	db = NewPostgreSQLDB()
	InitPushClient()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		rs := NewRedisStream(client)
		for range time.Tick(time.Second * 30) {
			if !getRedisLock() {
				if err := rs.CheckStream(); err != nil {
					log.Println("Error occurred:", err)
				}
			}
		}
	}()

	http.ListenAndServe(":8000", nil)
	wg.Wait()
}
