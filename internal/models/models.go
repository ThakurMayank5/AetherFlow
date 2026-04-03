package models

type Job struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Payload    string    `json:"payload"`
	Status     JobStatus `json:"status"`
	Retries    int       `json:"retries"`
	MaxRetries int       `json:"max_retries"`
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "PENDING"
	JobStatusProcessing JobStatus = "PROCESSING"
	JobStatusSuccess    JobStatus = "SUCCESS"
	JobStatusFailed     JobStatus = "FAILED"
)
