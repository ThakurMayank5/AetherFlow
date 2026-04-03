package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ThakurMayank5/AetherFlow/internal/models"
	"github.com/ThakurMayank5/AetherFlow/internal/queue"
)

func main() {
	q, err := queue.NewRedisQueue()
	if err != nil {
		log.Fatal("Failed to create Redis queue:", err)
	}

	fmt.Println("Worker started...")

	for {
		data, err := q.Dequeue()
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		var job models.Job
		err = json.Unmarshal([]byte(data), &job)
		if err != nil {
			log.Println("Invalid job:", err)
			continue
		}

		job.Status = "PROCESSING"
		q.SaveJob(job, job.ID)

		err = processJob(job)
		if err != nil {
			log.Println("Error processing job:", err)
		}
	}
}

func processJob(job models.Job) error {
	fmt.Printf("Processing job: %+v\n", job)

	// simulate work
	fmt.Println("Job done:", job.ID)

	return nil
}
