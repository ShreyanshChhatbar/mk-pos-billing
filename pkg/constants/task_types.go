package constants

// TaskType represents an async task type identifier
type TaskType string

const (
	// Core tasks
	TaskTypeCoreTest TaskType = "core:test"

	// Form tasks
	TaskTypeProcessFormSubmission TaskType = "forms:process_submission"
)

var AllTaskTypes = AllValues(
	TaskTypeCoreTest,
	TaskTypeProcessFormSubmission,
)

// QueueName represents an async queue name
type QueueName string

const (
	QueueCritical QueueName = "critical"
	QueueDefault  QueueName = "default"
	QueueLow      QueueName = "low"
)

var AllQueueNames = AllValues(
	QueueCritical,
	QueueDefault,
	QueueLow,
)
