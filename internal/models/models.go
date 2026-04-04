package models

type Job struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Payload    string      `json:"payload"`
	Status     JobStatus   `json:"status"`
	Retries    int         `json:"retries"`
	MaxRetries int         `json:"max_retries"`
	Priority   JobPriority `json:"priority"`
	Delayed    bool        `json:"delayed"`

	// This delay is in seconds from the time the job is created.
	Delay int64 `json:"delay"`
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "PENDING"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusSuccess    JobStatus = "SUCCESS"
	JobStatusFailed     JobStatus = "FAILED"
)

type JobPriority string

const (
	JobPriorityLow    JobPriority = "LOW"
	JobPriorityMedium JobPriority = "MEDIUM"
	JobPriorityHigh   JobPriority = "HIGH"
)
