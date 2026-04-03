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
		data, err := q.DequeueWithPriority()
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
			log.Println("Job failed:", err)
			handleJobFailure(q, job)
		} else {
			job.Status = "SUCCESS"
			q.SaveJob(job, job.ID)
		}
	}
}

func processJob(job models.Job) error {
	fmt.Printf("Processing job: %+v\n", job)

	// simulate work
	fmt.Println("Job done:", job.ID)

	// return fmt.Errorf("error simulation")
	return nil
}

func handleJobFailure(q *queue.RedisQueue, job models.Job) {
	job.Retries++

	if job.Retries < job.MaxRetries {
		fmt.Println("Retrying job:", job.ID)

		job.Status = "PENDING"
		q.SaveJob(job, job.ID)
		q.Enqueue(job)

	} else {
		fmt.Println("Job failed permanently:", job.ID)

		job.Status = "FAILED"
		q.SaveJob(job, job.ID)
	}
}
