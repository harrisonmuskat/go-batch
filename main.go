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

func worker(tasks <-chan User, wg *sync.WaitGroup) {
	defer wg.Done()

	for user := range tasks {
		sendSMSWithRetry(user, "Hello, World!", 3)
	}
}

func processWithWorkerPools(users []User, workers int) time.Duration {
	start := time.Now()
	tasks := make(chan User, len(users))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(tasks, &wg)
	}

	for _, user := range users {
		tasks <- user
	}

	close(tasks)

	wg.Wait()

	return time.Since(start)
}

func processWithBatches(users []User, batchSize int) time.Duration {
	start := time.Now()
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}
		batch := users[i:end]

		for _, user := range batch {
			sendSMSWithRetry(user, "Hello, World!", 3)
		}
	}

	return time.Since(start)
}

func main() {
	users := make([]User, 1000)
	for i := range users {
		users[i] = User{ID: i + 1, Name: fmt.Sprintf("User %d", i+1), Phone: fmt.Sprintf("123-456-%04d", i+1)}
	}

	batchSize := 3
	batchTime := processWithBatches(users, batchSize)

	workers := 3
	workerPoolTime := processWithWorkerPools(users, workers)

	fmt.Printf("Batch processing took %v\n", batchTime)
	fmt.Printf("Worker pool processing took %v\n", workerPoolTime)
}
