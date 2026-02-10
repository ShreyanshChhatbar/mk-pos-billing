package constants

type TaskTrackerStatus string

const (
	TaskTrackerStatusStart     TaskTrackerStatus = "START"
	TaskTrackerStatusEnd       TaskTrackerStatus = "END"
	TaskTrackerStatusCancelled TaskTrackerStatus = "CANCELLED"
	TaskTrackerStatusDeleted   TaskTrackerStatus = "DELETED"
	TaskTrackerStatusReleased  TaskTrackerStatus = "RELEASED"
)
