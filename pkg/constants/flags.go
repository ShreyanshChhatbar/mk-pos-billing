package constants

var (
	AllFlagStatuses   = AllValues(FlagStatusPending, FlagStatusCompleted, FlagStatusSubmitted)
	AllFlagPriorities = AllValues(FlagPriorityHigh, FlagPriorityMedium, FlagPriorityLow)
)

type FlagStatus string

const (
	FlagStatusPending   FlagStatus = "PENDING"
	FlagStatusCompleted FlagStatus = "COMPLETED"
	FlagStatusSubmitted FlagStatus = "SUBMITTED"
)

type FlagPriority string

const (
	FlagPriorityHigh   FlagPriority = "HIGH"
	FlagPriorityMedium FlagPriority = "MEDIUM"
	FlagPriorityLow    FlagPriority = "LOW"
)
