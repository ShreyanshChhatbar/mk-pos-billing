package constants

var (
	AllFrequencyTypes    = AllValues(FrequencyTypeOnce, FrequencyTypeDaily, FrequencyTypeWeekly, FrequencyTypeMonthly)
	AllScheduleTypes     = AllValues(ScheduleTypeMorning, ScheduleTypeEvening)
	AllEscalationTypes   = AllValues(EscalationTypeFlag, EscalationTypeEscalation)
	AllAssignableTypes   = AllValues(AssignableTypeUser, AssignableTypeRole, AssignableTypeFacility)
	AllTaskCriticalities = AllValues(TaskCriticalityLow, TaskCriticalityMedium, TaskCriticalityHigh)
	AllTaskPriorities    = AllValues(TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh)
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
	TaskStatusDeleted    TaskStatus = "DELETED"
	TaskStatusReleased   TaskStatus = "RELEASED"
	TaskStatusReAudit    TaskStatus = "RE_AUDIT"
	TaskStatusResolved   TaskStatus = "RESOLVED"
	TaskStatusExpired    TaskStatus = "EXPIRED"
	TaskStatusVerified   TaskStatus = "VERIFIED"
	TaskStatusChecked    TaskStatus = "CHECKED"
)

type AuditResolutionType string

const (
	AuditResolutionTypeCreateBill AuditResolutionType = "CREATE_BILL"
	AuditResolutionTypeNoAction   AuditResolutionType = "NO_ACTION"
)

type TaskCriticality string

const (
	TaskCriticalityLow    TaskCriticality = "LOW"
	TaskCriticalityMedium TaskCriticality = "MEDIUM"
	TaskCriticalityHigh   TaskCriticality = "HIGH"
)

type FrequencyType string

const (
	FrequencyTypeOnce    FrequencyType = "ONCE"
	FrequencyTypeDaily   FrequencyType = "DAILY"
	FrequencyTypeWeekly  FrequencyType = "WEEKLY"
	FrequencyTypeMonthly FrequencyType = "MONTHLY"
)

type TaskLevelType string

const (
	TaskLevelTypeUser     TaskLevelType = "USER"
	TaskLevelTypeRole     TaskLevelType = "ROLE"
	TaskLevelTypeFacility TaskLevelType = "FACILITY"
)

type AssignableType string

const (
	AssignableTypeUser     AssignableType = "USER"
	AssignableTypeRole     AssignableType = "ROLE"
	AssignableTypeFacility AssignableType = "FACILITY"
)

type ScheduleType string

const (
	ScheduleTypeMorning ScheduleType = "MORNING"
	ScheduleTypeEvening ScheduleType = "EVENING"
)

type EscalationType string

const (
	EscalationTypeFlag       EscalationType = "FLAG"
	EscalationTypeEscalation EscalationType = "ESCALATION"
)

type TaskPolymorphicType string

const (
	TaskPolymorphicStoreVisit TaskPolymorphicType = "StoreVisits"
	TaskPolymorphicAudit      TaskPolymorphicType = "Audit"
	TaskPolymorphicFormTask   TaskPolymorphicType = "FormTask"
	TaskPolymorphicTillTask   TaskPolymorphicType = "TillTaskUsers"
	TaskPolymorphicFlagDetail TaskPolymorphicType = "FlagDetail"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "LOW"
	TaskPriorityMedium TaskPriority = "MEDIUM"
	TaskPriorityHigh   TaskPriority = "HIGH"
)

type TaskCreationType string

const (
	TaskCreationTypeCron   TaskCreationType = "CRON"
	TaskCreationTypeManual TaskCreationType = "MANUAL"
	TaskCreationTypeSystem TaskCreationType = "SYSTEM"
)

type TillTaskStatus string

const (
	TillTaskStatusPending  TillTaskStatus = "PENDING"
	TillTaskStatusChecked  TillTaskStatus = "CHECKED"
	TillTaskStatusVerified TillTaskStatus = "VERIFIED"
)

type TillTaskUserStatus string

const (
	TillTaskUserStatusPending   TillTaskUserStatus = "PENDING"
	TillTaskUserStatusCompleted TillTaskUserStatus = "COMPLETED"
	TillTaskUserStatusExpired   TillTaskUserStatus = "EXPIRED"
)

const (
	FacilityDetailKeyWSAlternateCode = "ws_alternate_code"
)

type AuditType string

const (
	AuditTaskTypeStoreAudit AuditType = "STORE_AUDIT"
	AuditTaskTypeAsmAudit   AuditType = "ASM_AUDIT"
	AuditTaskTypeReAudit    AuditType = "RE_AUDIT"
)

type AuditDetailQuantityStatus string

const (
	AuditDetailQuantityStatusEqual AuditDetailQuantityStatus = "EQUAL"
	AuditDetailQuantityStatusLess  AuditDetailQuantityStatus = "LESS"
	AuditDetailQuantityStatusMore  AuditDetailQuantityStatus = "MORE"
)

var (
	ProductLocationRack = []string{
		"Fridge",
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
		"J",
		"K",
		"OTCA",
		"OTCB",
		"OTCC",
		"OTCD",
		"OTCE",
		"OTCF",
		"OTCG",
		"OTCH",
		"OTCI",
		"UTC1",
		"UTC2",
		"UTC3",
		"UTC4",
	}

	ProductLocationShelf = []string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
	}

	ProductLocationBox = []string{
		"M",
		"N",
		"O",
	}
)

type AuditTaskDetailType string

const (
	AuditTaskDetailTypeNewProduct AuditTaskDetailType = "NEW_PRODUCT"
)
