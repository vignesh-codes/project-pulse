package main

import "github.com/go-redis/redis/v8"

type RedisStream struct {
	client *redis.Client
}
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

func (Event) TableName() string {
	return "events"
}
