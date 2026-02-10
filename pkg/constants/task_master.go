package constants

type TaskMasterType string

const (
	TaskMasterTypeCustomModule TaskMasterType = "CUSTOM_MODULE"
	TaskMasterTypeForm         TaskMasterType = "FORM"
	TaskMasterTypeQuiz         TaskMasterType = "QUIZ"
)

type TaskMasterConstant string

const (
	TaskTypeStoreVisit              TaskMasterConstant = "STORE_VISIT"
	TaskMasterConstantAsmStoreAudit TaskMasterConstant = "ASM_STORE_AUDIT"
	TaskMasterConstantTillTask      TaskMasterConstant = "TILL_TASK"
	TaskMasterConstantFlagDefault   TaskMasterConstant = "FLAG_DEFAULT"
	TaskMasterConstantStoreAudit    TaskMasterConstant = "STORE_AUDIT"
)

var CustomModuleTypes = map[TaskMasterConstant]struct{}{
	TaskMasterConstantStoreAudit: {},
}

// CustomModuleDropdownTypes lists custom modules that expose dropdown helpers.
var CustomModuleDropdownTypes = map[TaskMasterConstant]struct{}{
	TaskMasterConstantStoreAudit: {},
}
