package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThakurMayank5/AetherFlow/internal/models"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisQueue struct {
	Client *redis.Client
}

func (q *RedisQueue) GetContext() context.Context {
	return ctx
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

	return &RedisQueue{Client: rdb}, nil
}

func (q *RedisQueue) Enqueue(job models.Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.Client.LPush(ctx, "jobs", data).Err()
}

func (q *RedisQueue) EnqueueWithPriority(job models.Job, priority models.JobPriority) error {
	data, err := json.Marshal(job)

	if err != nil {
		return err
	}

	queueName := string("jobs:" + priority)
	return q.Client.LPush(ctx, queueName, data).Err()
}

func (q *RedisQueue) EnqueueDelayed(job models.Job, delay time.Duration) error {
	data, _ := json.Marshal(job)

	execTime := time.Now().Add(delay).Unix()

	return q.Client.ZAdd(ctx, "jobs:delayed", redis.Z{
		Score:  float64(execTime),
		Member: data,
	}).Err()
}

func (q *RedisQueue) Dequeue() (string, error) {
	result, err := q.Client.BRPop(ctx, 0, "jobs").Result()
	if err != nil {
		return "", err
	}

	return result[1], nil
}

func (q *RedisQueue) DequeueWithPriority() (string, error) {
	result, err := q.Client.BRPop(ctx, 0,
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

	return q.Client.Set(ctx, "job:"+id, data, 0).Err()
}

func (q *RedisQueue) GetJob(id string) (string, error) {
	return q.Client.Get(ctx, "job:"+id).Result()
}
