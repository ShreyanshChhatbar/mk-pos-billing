package validation

import (
	"mk-pos-billing/pkg/constants"
	"mk-pos-billing/pkg/utils"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate = validator.New()

func init() {
	// Register custom validators
	validate.RegisterValidation("escalation_type", validateEscalationType)
	validate.RegisterValidation("frequency_type", validateFrequencyType)
	validate.RegisterValidation("schedule_type", validateScheduleType)
	validate.RegisterValidation("assignable_type", validateAssignableType)
	validate.RegisterValidation("flag_status", validateFlagStatus)
	validate.RegisterValidation("task_priority", validateTaskPriority)
}

// validateEscalationType validates if the value is a valid EscalationType
func validateEscalationType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return constants.IsValidEnum(value, constants.AllEscalationTypes...)
}

// validateFrequencyType validates if the value is a valid FrequencyType
func validateFrequencyType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return constants.IsValidEnum(value, constants.AllFrequencyTypes...)
}

// validateScheduleType validates if the value is a valid ScheduleType
func validateScheduleType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return constants.IsValidEnum(value, constants.AllScheduleTypes...)
}

// validateAssignableType validates if the value is a valid AssignableType
func validateAssignableType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return constants.IsValidEnum(value, constants.AllAssignableTypes...)
}

// validateFlagStatus validates if the value is a valid FlagStatus
func validateFlagStatus(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return constants.IsValidEnum(value, constants.AllFlagStatuses...)
}

// validateTaskPriority validates if the value is a valid TaskPriority
func validateTaskPriority(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return constants.IsValidEnum(value, constants.AllTaskPriorities...)
}

func Validate[U any, Q any, B any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var uri U
		var query Q
		var body B

		if err := c.ShouldBindUri(&uri); err != nil {
			SendErrorResponse(c, uri, err)
			return
		}

		if err := c.ShouldBindQuery(&query); err != nil {
			SendErrorResponse(c, query, err)
			return
		}

		if err := c.ShouldBindJSON(&body); err != nil && err.Error() != "EOF" {
			SendErrorResponse(c, body, err)
			return
		}

		if err := validateValue(uri); err != nil {
			SendErrorResponse(c, uri, err)
			return
		}
		if err := validateValue(query); err != nil {
			SendErrorResponse(c, query, err)
			return
		}
		if err := validateValue(body); err != nil {
			SendErrorResponse(c, body, err)
			return
		}

		c.Set("validated_uri", uri)
		c.Set("validated_query", query)
		c.Set("validated_body", body)
		c.Next()
	}
}

func validateValue(v interface{}) error {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}
	return validate.Struct(v)
}

func ValidateQuery[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindQuery(&req); err != nil {
			SendErrorResponse(c, req, err)
			return
		}
		if err := validate.Struct(req); err != nil {
			SendErrorResponse(c, req, err)
			return
		}
		c.Set("request.query", req)
		c.Next()
	}
}

func ValidateBody[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			SendErrorResponse(c, req, err)
			return
		}
		if err := validate.Struct(req); err != nil {
			SendErrorResponse(c, req, err)
			return
		}
		c.Set("request.body", req)
		c.Next()
	}
}

func ValidateMultiTypeBody(types ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, t := range types {
			// Create a new instance of the type
			val := reflect.New(reflect.TypeOf(t)).Interface()

			// Try to bind - ShouldBindBodyWith allows multiple reads
			if err := c.ShouldBindBodyWith(val, binding.JSON); err == nil {
				// Binding succeeded for this type, now validate the struct
				if err := validate.Struct(val); err != nil {
					// Binding worked but validation failed - return specific validation errors
					SendErrorResponse(c, val, err)
					c.Abort()
					return
				}

				// Success: Both binding and validation passed
				// If it's a pointer, get the underlying value for the context
				result := val
				if reflect.TypeOf(val).Kind() == reflect.Ptr {
					result = reflect.ValueOf(val).Elem().Interface()
				}
				c.Set("validated_body", result)
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, utils.APIResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request body format (Does not match any allowed types)",
		})
	}
}

func GetValidParamData[T any](c *gin.Context) T {
	v, _ := c.Get("validated_query")
	val, _ := v.(T)
	return val
}

func GetValidBodyData[T any](c *gin.Context) T {
	v, _ := c.Get("validated_body")
	val, _ := v.(T)
	return val
}

func GetValidUriData[T any](c *gin.Context) T {
	v, _ := c.Get("validated_uri")
	val, _ := v.(T)
	return val
}

func ValidateRequestData[T any](c *gin.Context, key string) T {
	if v, ok := c.Get(key); ok {
		if val, ok := v.(T); ok {
			return val
		}
	}
	var empty T
	return empty
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func SendErrorResponse(c *gin.Context, req interface{}, err error) {
	zap.L().Error("Error Type:", zap.Error(err))

	if ve, ok := err.(validator.ValidationErrors); ok {
		errors := formatValidationError(req, ve)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, utils.APIResponse{
			Success:    false,
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Validation failed",
			Errors:     errors,
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, utils.APIResponse{
		Success:    false,
		StatusCode: http.StatusUnprocessableEntity,
		Message:    "Invalid request format",
		Errors:     gin.H{"error": err.Error()},
	})
}

func formatValidationError(req interface{}, err error) map[string]string {
	errors := map[string]string{}
	var customMsgs map[string]string
	var attributes map[string]string

	if m, ok := req.(interface{ Messages() map[string]string }); ok {
		customMsgs = m.Messages()
	}
	if a, ok := req.(interface{ Attributes() map[string]string }); ok {
		attributes = a.Attributes()
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			fieldName := getJSONOrFormTag(req, fe.StructField())
			errors[fieldName] = buildFieldErrorMessage(fe, fieldName, customMsgs, attributes)
		}
	} else {
		errors["error"] = err.Error()
	}

	return errors
}

func buildFieldErrorMessage(fe validator.FieldError, field string, customMsgs map[string]string, attrs map[string]string) string {
	if msg, exists := customMsgs[field+"."+fe.Tag()]; exists {
		return msg
	}
	if label, ok := attrs[field]; ok {
		field = label
	}
	switch fe.Tag() {
	case "required", "required_without":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	default:
		return fe.Error()
	}
}

func getJSONOrFormTag(obj interface{}, structField string) string {
	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, ok := t.FieldByName(structField)
	if !ok {
		return strings.ToLower(structField)
	}

	if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
		return strings.Split(jsonTag, ",")[0]
	}
	if formTag := field.Tag.Get("form"); formTag != "" && formTag != "-" {
		return strings.Split(formTag, ",")[0]
	}
	return strings.ToLower(structField)
}
