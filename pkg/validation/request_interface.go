package validation

type ValidateRequest interface {
	// Optional: custom error messages
	Messages() map[string]string

	// Optional: custom field names
	Attributes() map[string]string
}
