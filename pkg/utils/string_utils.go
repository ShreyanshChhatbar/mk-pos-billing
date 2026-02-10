package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// Example: "My Task Name" → "MY_TASK_NAME"
func GenerateIdentifierSlug(s string) string {
	s = strings.TrimSpace(s)

	spaceRegex := regexp.MustCompile(`\s+`)
	s = spaceRegex.ReplaceAllString(s, " ")

	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, " ", "_")

	return s
}

// Example: "hello world" → "Hello world"
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Example: "My Task Name" → "MY_TASK_NAME_1234"
func CreateConstant(s string) string {
	baseSlug := GenerateIdentifierSlug(s)
	uniqueSuffix := strings.ToUpper(uuid.New().String()[0:4])
	return fmt.Sprintf("%s_%s", baseSlug, uniqueSuffix)
}

func MaskMiddle(input string) string {
	length := len(input)
	if length <= 4 {
		return strings.Repeat("*", length)
	}
	visibleStart := input[:2]
	visibleEnd := input[length-2:]
	maskedMiddle := strings.Repeat("*", length-4)
	return visibleStart + maskedMiddle + visibleEnd
}

func FormatDateDDMMYYYY(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}

	return t.Format("02/01/2006")
}

func ConvertTo12HourFormat(timeStr string) string {
	if strings.TrimSpace(timeStr) == "" {
		return ""
	}

	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return timeStr
	}

	return t.Format("03:04:05 PM")
}

func ToUint64(s string) uint64 {
	id, _ := strconv.ParseUint(s, 10, 64)
	return id
}

func Uint64ToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}

func StringPtr(s string) *string {
	return &s
}
