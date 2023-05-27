package main

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var lastTimestamp int64

// func GenerateUniqueID(str string) string {
// 	id := uuid.New([]byte(str))

// 	return id.String()
// }

func main() {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	idSet := make(map[int64]bool)
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println("Error creating Snowflake node:", err)
		return
	}
	// Number of concurrent calls
	concurrentCalls := 10

	wg.Add(concurrentCalls)

	for i := 0; i < concurrentCalls; i++ {
		go func(i int) {
			defer wg.Done()

			// str := "example" + strconv.Itoa(i)
			uniqueID := node.Generate()
			fmt.Println("uni ", uniqueID)
			mutex.Lock()
			if _, exists := idSet[uniqueID.Int64()]; exists {
				fmt.Printf("Duplicate ID found: %d\n", uniqueID)
			} else {
				idSet[uniqueID.Int64()] = true
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()

	fmt.Println("Testing complete")
}
