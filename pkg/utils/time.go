package utils

import "time"

const DefaultDateTimeLayout = "2006-01-02 15:04:05"

// FormatTime formats a time with the provided layout. Falls back to a standard yyyy-mm-dd HH:MM:SS layout if empty.
func FormatTime(t time.Time, layout string) string {
	if t.IsZero() {
		return ""
	}
	if layout == "" {
		layout = DefaultDateTimeLayout
	}
	return t.UTC().Format(layout)
}
