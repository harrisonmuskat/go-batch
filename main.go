package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type User struct {
	ID    int
	Name  string
	Phone string
}

func enrichUser(userID int) User {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	return User{
		ID:    userID,
		Name:  fmt.Sprintf("User %d", userID),
		Phone: fmt.Sprintf("123-456-%04d", userID),
	}
}

func sendSMS(User User, message string) error {
	if rand.Float32() < 0.3 {
		return fmt.Errorf("Failed to send SMS to %d", User.ID)
	}
	fmt.Printf("Sent SMS to %d: %s\n", User.ID, message)
	return nil
}

func sendSMSWithRetry(User User, message string, maxRetries int) error {
	for attempts := 1; attempts <= maxRetries; attempts++ {
		err := sendSMS(User, message)
		if err == nil {
			return nil
		}
		fmt.Printf("Failed to send SMS to %d, retry (%d)...\n", User.ID, attempts)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
	return fmt.Errorf("Failed to send SMS to %d after %d attempts", User.ID, maxRetries)
}

func processBatch(batch []User) {
	var wg sync.WaitGroup

	for _, user := range batch {
		wg.Add(1)
		go func(user User) {
			defer wg.Done()
			err := sendSMSWithRetry(user, "Hello, World!", 3)
			if err != nil {
				fmt.Println(err)
			}
		}(user)
	}

	wg.Wait()
}

func main() {
	batch := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	batchSize := 3

	for i := 0; i < len(batch); i += batchSize {
		end := i + batchSize
		if end > len(batch) {
			end = len(batch)
		}
		batch := batch[i:end]

		var enrichedBatch []User
		for _, userID := range batch {
			enrichedUser := enrichUser(userID)
			enrichedBatch = append(enrichedBatch, enrichedUser)
		}
		processBatch(enrichedBatch)
	}
}
