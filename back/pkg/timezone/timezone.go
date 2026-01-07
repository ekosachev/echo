package timezone

import (
	"fmt"
	"time"
)

var (
	SupportedTimezones = []string{
		"UTC",
		"Europe/Moscow",
		"Europe/London",
		"America/New_York",
		"America/Los_Angeles",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Asia/Novosibirsk",
	}
)

func IsValidTimezone(tz string) bool {
	_, err := time.LoadLocation(tz)
	if err != nil {
		return false
	}
	return true
}

func ConvertToUTC(localTime time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %w", err)
	}

	return time.Date(
		localTime.Year(),
		localTime.Month(),
		localTime.Day(),
		localTime.Hour(),
		localTime.Minute(),
		localTime.Second(),
		localTime.Nanosecond(),
		loc,
	).UTC(), nil
}

func ConvertFromUTC(utcTime time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %w", err)
	}

	return utcTime.In(loc), nil
}

func GetUserTimezone(userTimezone string) string {
	if userTimezone == "" || !IsValidTimezone(userTimezone) {
		return "UTC"
	}
	return userTimezone
}
