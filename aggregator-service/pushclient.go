package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
)

var start = "-"
var prevStart = start

// maintaining mutex to handle concurrent read write on same var
var (
	redis_lock bool
	mutex      sync.Mutex
)

func updateRedisLock(b bool) {
	mutex.Lock()
	defer mutex.Unlock()

	redis_lock = b
}

func getRedisLock() bool {
	mutex.Lock()
	defer mutex.Unlock()

	return redis_lock
}

// must be from env
const redis_key string = "mystream"

func InitPushClient() {
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var logEntry []map[string]interface{}
		err = json.Unmarshal([]byte(string(body)), &logEntry)
		if err != nil {
			http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
			fmt.Println("err ", err)
			return
		}
		// Push the log entry to the Redis stream
		for _, entry := range logEntry {
			_, err := client.XAdd(context.Background(), &redis.XAddArgs{
				Stream: redis_key,
				ID:     "*",
				Values: entry,
			}).Result()

			if err != nil {
				http.Error(w, "Failed to push log entry to Redis stream", http.StatusInternalServerError)
				fmt.Println("redis error is ", err)
				return
			}
		}

		count += 1
		fmt.Fprint(w, "Log entry pushed to Redis stream successfully")
	})
}

/*
checks redis streams
if count > 1, processes it
if count is > 1000, it keeps processing it and locks further execution of checkStream
function until all logs are process at this instance itself
*/
func (rs *RedisStream) CheckStream() error {
	defer func() {
		updateRedisLock(false)
	}()

	updateRedisLock(true)
	fmt.Println("Checking Redis...")
	c, err := rs.client.XLen(context.Background(), redis_key).Result()
	if err != nil {
		return fmt.Errorf("failed to get stream length: %v", err)
	}

	for c >= 1 {
		if c > 1000 {
			c = 1000
		}

		fmt.Printf("Processing %d messages...\n", c)

		items, err := rs.client.XRangeN(context.Background(), redis_key, start, "+", c).Result()
		if err != nil {
			return fmt.Errorf("failed to get stream items: %v", err)
		}

		var buff = []redis.XMessage{}
		var id_buff = map[string]redis.XMessage{}
		todel := make([]string, 0, len(items))

		for _, item := range items {
			x, f := id_buff[item.ID]
			// Just a precautionary check to prevent duplicate ids.
			if f {
				fmt.Println("found a duplicate key ", item, " \n ", x)
				continue
			}

			buff = append(buff, item)
			id_buff[item.ID] = item
			todel = append(todel, item.ID)
			start = item.ID
		}

		if len(buff) > 0 {
			if err := db.PushToPostgres(buff); err != nil {
				return fmt.Errorf("failed to push messages to PostgreSQL: %v", err)
			}
		} else {
			return nil
		}

		_, er := rs.client.XDel(context.Background(), redis_key, todel...).Result()
		if er != nil {
			return nil
		}
		fmt.Println(len(buff), "messages pushed to PostgreSQL")
		fmt.Println("Deleted messages from Redis")

		c, err = rs.client.XLen(context.Background(), redis_key).Result()
		if err != nil {
			return fmt.Errorf("failed to get stream length: %v", err)
		}
	}

	return nil
}

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
