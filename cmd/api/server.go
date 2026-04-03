package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ThakurMayank5/AetherFlow/internal/models"
	"github.com/ThakurMayank5/AetherFlow/internal/queue"
)

func main() {
	r := gin.Default()

	q, err := queue.NewRedisQueue()

	if err != nil {
		log.Fatal("Failed to initialize Redis queue: " + err.Error())
	}

	r.POST("/jobs", func(c *gin.Context) {
		var job models.Job

		if err := c.ShouldBindJSON(&job); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		job.ID = uuid.New().String()

		job.Status = models.JobStatusPending
		job.Retries = 0
		job.MaxRetries = 3

		if job.Priority == "" {
			job.Priority = models.JobPriorityMedium
		}

		if job.Priority != models.JobPriorityLow && job.Priority != models.JobPriorityMedium && job.Priority != models.JobPriorityHigh {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid priority"})
			return
		}

		err := q.SaveJob(job, job.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save job"})
			return
		}

		err = q.EnqueueWithPriority(job, job.Priority)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"job_id": job.ID,
		})
	})

	r.GET("/jobs/:id", func(c *gin.Context) {
		id := c.Param("id")

		data, err := q.GetJob(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
			return
		}

		var job models.Job
		json.Unmarshal([]byte(data), &job)

		c.JSON(http.StatusOK, job)
	})

	r.Run(":42069")
}
