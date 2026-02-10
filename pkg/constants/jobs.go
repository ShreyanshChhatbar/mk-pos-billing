package constants

type JobExecutionStatus string

const (
	JobStatusPending   JobExecutionStatus = "PENDING"
	JobStatusActive    JobExecutionStatus = "ACTIVE"
	JobStatusCompleted JobExecutionStatus = "COMPLETED"
	JobStatusFailed    JobExecutionStatus = "FAILED"
	JobStatusRetrying  JobExecutionStatus = "RETRYING"
)

var AllJobExecutionStatuses = AllValues(
	JobStatusPending,
	JobStatusActive,
	JobStatusCompleted,
	JobStatusFailed,
	JobStatusRetrying,
)
