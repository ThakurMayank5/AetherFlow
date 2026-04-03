package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThakurMayank5/AetherFlow/internal/models"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue() (*RedisQueue, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Testing connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisQueue{client: rdb}, nil
}

func (q *RedisQueue) Enqueue(job models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.client.LPush(ctx, "jobs", data).Err()
}

func (q *RedisQueue) EnqueueWithPriority(job models.Job, priority models.JobPriority) error {
	data, err := json.Marshal(job)

	if err != nil {
		return err
	}

	queueName := string("jobs:" + priority)
	return q.client.LPush(ctx, queueName, data).Err()
}

func (q *RedisQueue) Dequeue() (string, error) {
	result, err := q.client.BRPop(ctx, 0, "jobs").Result()
	if err != nil {
		return "", err
	}

	return result[1], nil
}

func (q *RedisQueue) DequeueWithPriority() (string, error) {
	result, err := q.client.BRPop(ctx, 0,
		"jobs:HIGH",
		"jobs:MEDIUM",
		"jobs:LOW",
	).Result()

	if err != nil {
		return "", err
	}

	return result[1], nil
}

func (q *RedisQueue) SaveJob(job models.Job, id string) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.client.Set(ctx, "job:"+id, data, 0).Err()
}

func (q *RedisQueue) GetJob(id string) (string, error) {
	return q.client.Get(ctx, "job:"+id).Result()
}
