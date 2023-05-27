package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var start = "-"

// docker run --name mypostgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -p 4432:5432 -d postgres

type LogEntry struct {
	ID     string `json:"ID"`
	Values struct {
		Date string `json:"date"`
		Log  string `json:"log"`
	} `json:"Values"`
}

type HP struct {
	Timestamp   float64 `json:"timestamp"`
	EventSource string  `json:"event_source"`
	Event       Event   `json:"payload"`
}

type Event struct {
	ID        int    `gorm:"column:id;primaryKey" json:"id"`
	UnixTS    int    `gorm:"column:unix_ts" json:"unix_ts"`
	UserID    int    `gorm:"column:user_id" json:"user_id"`
	EventName string `gorm:"column:event_name" json:"event_name"`
}

// TableName specifies the table name for the Event struct
func (Event) TableName() string {
	return "events"
}

// RedisStream represents the Redis stream and provides methods to interact with it.
type RedisStream struct {
	client *redis.Client
}

// PostgreSQLDB represents the PostgreSQL database and provides methods to interact with it.
type PostgreSQLDB struct {
	db *gorm.DB
}

var db *PostgreSQLDB

// NewRedisStream creates a new RedisStream instance.
func NewRedisStream(client *redis.Client) *RedisStream {
	return &RedisStream{
		client: client,
	}
}

// CheckStream checks the Redis stream and processes the messages.
func (rs *RedisStream) CheckStream() error {
	fmt.Println("Checking Redis...")
	c, err := rs.client.XLen(context.Background(), "mystream").Result()
	if err != nil {
		return fmt.Errorf("failed to get stream length: %v", err)
	}

	if c >= 5 {
		fmt.Println("Found more than 5 messages:", c)
		items, err := rs.client.XRangeN(context.Background(), "mystream", start, "+", 1000).Result()
		if err != nil {
			return fmt.Errorf("failed to get stream items: %v", err)
		}

		var buff []redis.XMessage
		prevStart := start
		todel := make([]string, 0, len(items))

		for _, item := range items {
			buff = append(buff, item)
			todel = append(todel, item.ID)
			start = items[len(items)-1].ID
		}

		if err := db.PushToPostgres(buff); err != nil {
			return fmt.Errorf("failed to push messages to PostgreSQL: %v", err)
		}

		c, er := rs.client.XDel(context.Background(), "mystream", todel...).Result()
		fmt.Println(c, er, prevStart, start)
	}

	return nil
}

// NewPostgreSQLDB creates a new PostgreSQLDB instance.
func NewPostgreSQLDB() *PostgreSQLDB {
	// PostgreSQL connection details
	dsn := "host=localhost user=user password=pass dbname=postgres port=4432 sslmode=disable"

	// Establish connection to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate the table if it doesn't exist
	if err := db.AutoMigrate(&Event{}); err != nil {
		log.Fatalf("Failed to perform auto-migration: %v", err)
	}
	fmt.Println("TABLE AVAILABLE ")
	return &PostgreSQLDB{
		db: db,
	}
}

// PushToPostgres pushes the messages to PostgreSQL.
func (pg *PostgreSQLDB) PushToPostgres(b []redis.XMessage) error {
	jsonData, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	var logs []LogEntry
	if err := json.Unmarshal(jsonData, &logs); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	var payloadData []Event
	for _, entry := range logs {
		var payload HP

		if err := json.Unmarshal([]byte(entry.Values.Log), &payload); err != nil {
			log.Println("Failed to parse log entry:", err)
			continue
		}

		payloadData = append(payloadData, payload.Event)
	}

	if err := pg.db.Create(&payloadData).Error; err != nil {
		return fmt.Errorf("failed to insert batch into the database: %v", err)
	}

	fmt.Println("Batch insert completed.")

	return nil
}

func main() {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	db = NewPostgreSQLDB()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		rs := NewRedisStream(client)
		for range time.Tick(time.Second * 30) {
			if err := rs.CheckStream(); err != nil {
				log.Println("Error occurred:", err)
			}
		}
	}()
	wg.Wait()
}
