package main

import (
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

		err := q.Enqueue(job)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"job_id": job.ID,
		})
	})

	r.Run(":42069")
}
