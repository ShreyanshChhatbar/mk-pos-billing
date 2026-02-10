package constants

// FormFieldType represents the type of a form field
type FormFieldType string

const (
	FormFieldTypeText        FormFieldType = "text"
	FormFieldTypeSelect      FormFieldType = "select"
	FormFieldTypeMultiSelect FormFieldType = "multi_select"
	FormFieldTypeNumber      FormFieldType = "number"
	FormFieldTypeURL         FormFieldType = "url"
	FormFieldTypeEmail       FormFieldType = "email"
	FormFieldTypePhoneNumber FormFieldType = "phone_number"
	FormFieldTypeDate        FormFieldType = "date"
	FormFieldTypeFiles       FormFieldType = "files"
)
