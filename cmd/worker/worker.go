package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ThakurMayank5/AetherFlow/internal/models"
	"github.com/ThakurMayank5/AetherFlow/internal/queue"
	"github.com/redis/go-redis/v9"
)

func main() {
	q, err := queue.NewRedisQueue()
	if err != nil {
		log.Fatal("Failed to create Redis queue:", err)
	}

	go delayedJobsPoller(q)

	fmt.Println("Workers are starting...")

	// Spawning multiple worker goroutines
	for i := 0; i < 5; i++ {
		go workerLoop(i, q)
	}

	select {}

}

func workerLoop(id int, q *queue.RedisQueue) {

	fmt.Printf("Worker %d started\n", id)

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

		err = processJob(job, id)
		if err != nil {
			log.Println("Job failed:", err)
			handleJobFailure(q, job)
		} else {
			job.Status = "SUCCESS"
			q.SaveJob(job, job.ID)
		}
	}
}
func processJob(job models.Job, workerId int) error {
	log.Printf("Worker %d processing job: %+v\n", workerId, job)

	// simulate work
	log.Println("Job done:", job.ID)

	// return fmt.Errorf("error simulation")
	return nil
}

func handleJobFailure(q *queue.RedisQueue, job models.Job) {
	job.Retries++

	if job.Retries < job.MaxRetries {
		fmt.Println("Retrying job:", job.ID)

		job.Status = "PENDING"
		q.SaveJob(job, job.ID)
		q.EnqueueWithPriority(job, job.Priority)

	} else {
		fmt.Println("Moving Job to Dead Letter Queue:", job.ID)

		job.Status = "FAILED"
		q.SaveJob(job, job.ID)

		data, err := json.Marshal(job)

		if err != nil {
			log.Println("Failed to marshal job for DLQ:", err)
			return
		}

		ctx := q.GetContext()

		q.Client.LPush(ctx, "jobs:dlq", data)
	}
}

func delayedJobsPoller(q *queue.RedisQueue) {

	for {
		now := time.Now().Unix()

		ctx := q.GetContext()

		jobs, err := q.Client.ZRangeArgs(ctx, redis.ZRangeArgs{
			Key:     "jobs:delayed",
			Start:   "0",
			Stop:    fmt.Sprintf("%d", now),
			ByScore: true,
		}).Result()

		if err != nil {
			log.Println("Error fetching delayed jobs:", err)
			continue
		}

		for _, jobStr := range jobs {

			var job models.Job
			err := json.Unmarshal([]byte(jobStr), &job)
			if err != nil {
				log.Println("Invalid delayed job:", err)
				continue
			}

			q.EnqueueWithPriority(job, job.Priority)
			q.Client.ZRem(ctx, "jobs:delayed", jobStr)
		}

		time.Sleep(1 * time.Second)
	}

}
