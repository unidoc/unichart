package dataset

import (
	"fmt"
	"strconv"
	"time"
)

const (
	// defaultDateFormat is the default date format.
	defaultDateFormat = "2006-01-02"

	// defaultDateHourFormat is the date format for hour timestamp formats.
	defaultDateHourFormat = "01-02 3PM"

	// defaultDateMinuteFormat is the date format for minute range timestamp formats.
	defaultDateMinuteFormat = "01-02 3:04PM"

	// defaultFloatFormat is the default float format.
	defaultFloatFormat = "%.2f"

	// defaultPercentValueFormat is the default percent format.
	defaultPercentValueFormat = "%0.2f%%"
)

// ValueFormatterProvider is a series that has custom formatters.
type ValueFormatterProvider interface {
	GetValueFormatters() (x, y ValueFormatter)
}

// ValueFormatter is a function that takes a value and produces a string.
type ValueFormatter func(v interface{}) string

// TimeValueFormatter is a ValueFormatter for timestamps.
func TimeValueFormatter(v interface{}) string {
	return formatTime(v, defaultDateFormat)
}

// TimeHourValueFormatter is a ValueFormatter for timestamps.
func TimeHourValueFormatter(v interface{}) string {
	return formatTime(v, defaultDateHourFormat)
}

// TimeMinuteValueFormatter is a ValueFormatter for timestamps.
func TimeMinuteValueFormatter(v interface{}) string {
	return formatTime(v, defaultDateMinuteFormat)
}

// TimeDateValueFormatter is a ValueFormatter for timestamps.
func TimeDateValueFormatter(v interface{}) string {
	return formatTime(v, "2006-01-02")
}

// TimeValueFormatterWithFormat returns a time formatter with a given format.
func TimeValueFormatterWithFormat(format string) ValueFormatter {
	return func(v interface{}) string {
		return formatTime(v, format)
	}
}

// TimeValueFormatterWithFormat is a ValueFormatter for timestamps with a given format.
func formatTime(v interface{}, dateFormat string) string {
	if typed, isTyped := v.(time.Time); isTyped {
		return typed.Format(dateFormat)
	}
	if typed, isTyped := v.(int64); isTyped {
		return time.Unix(0, typed).Format(dateFormat)
	}
	if typed, isTyped := v.(float64); isTyped {
		return time.Unix(0, int64(typed)).Format(dateFormat)
	}
	return ""
}

// IntValueFormatter is a ValueFormatter for float64.
func IntValueFormatter(v interface{}) string {
	switch v.(type) {
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case float32:
		return strconv.FormatInt(int64(v.(float32)), 10)
	case float64:
		return strconv.FormatInt(int64(v.(float64)), 10)
	default:
		return ""
	}
}

// FloatValueFormatter is a ValueFormatter for float64.
func FloatValueFormatter(v interface{}) string {
	return FloatValueFormatterWithFormat(v, defaultFloatFormat)
}

// PercentValueFormatter is a formatter for percent values.
// NOTE: it normalizes the values, i.e. multiplies by 100.0.
func PercentValueFormatter(v interface{}) string {
	if typed, isTyped := v.(float64); isTyped {
		return FloatValueFormatterWithFormat(typed*100.0, defaultPercentValueFormat)
	}
	return ""
}

// FloatValueFormatterWithFormat is a ValueFormatter for float64 with a given format.
func FloatValueFormatterWithFormat(v interface{}, floatFormat string) string {
	if typed, isTyped := v.(int); isTyped {
		return fmt.Sprintf(floatFormat, float64(typed))
	}
	if typed, isTyped := v.(int64); isTyped {
		return fmt.Sprintf(floatFormat, float64(typed))
	}
	if typed, isTyped := v.(float32); isTyped {
		return fmt.Sprintf(floatFormat, typed)
	}
	if typed, isTyped := v.(float64); isTyped {
		return fmt.Sprintf(floatFormat, typed)
	}
	return ""
}

// KValueFormatter is a formatter for K values.
func KValueFormatter(k float64, vf ValueFormatter) ValueFormatter {
	return func(v interface{}) string {
		return fmt.Sprintf("%0.0fσ %s", k, vf(v))
	}
}
