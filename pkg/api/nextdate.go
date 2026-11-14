package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date %q: %w", dstart, err)
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		return nextByDays(now, date, parts, repeat)
	case "y":
		return nextByYear(now, date, parts, repeat)
	case "w":
		return nextByWeekday(now, date, parts, repeat)
	case "m":
		return nextByMonthDay(now, date, parts, repeat)
	default:
		return "", fmt.Errorf("unsupported repeat rule: %q", repeat)
	}
}

func nextByDays(now, date time.Time, parts []string, repeat string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid repeat format: %q", repeat)
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil || n < 1 || n > 400 {
		return "", fmt.Errorf("invalid day interval: %q", parts[1])
	}
	for {
		date = date.AddDate(0, 0, n)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextByYear(now, date time.Time, parts []string, repeat string) (string, error) {
	if len(parts) != 1 {
		return "", fmt.Errorf("invalid repeat format: %q", repeat)
	}
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextByWeekday(now, date time.Time, parts []string, repeat string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid repeat format: %q", repeat)
	}
	weekdays, err := parseWeekdays(parts[1])
	if err != nil {
		return "", err
	}
	for {
		date = date.AddDate(0, 0, 1)
		if weekdays[weekdayISO(date)] && afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextByMonthDay(now, date time.Time, parts []string, repeat string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("invalid repeat format: %q", repeat)
	}
	days, err := parseMonthDays(parts[1])
	if err != nil {
		return "", err
	}
	var months map[int]bool
	if len(parts) == 3 {
		months, err = parseMonths(parts[2])
		if err != nil {
			return "", err
		}
	}
	for {
		date = date.AddDate(0, 0, 1)
		if !afterNow(date, now) {
			continue
		}
		if months != nil && !months[int(date.Month())] {
			continue
		}
		if matchMonthDay(date, days) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return y1 > y2 ||
		(y1 == y2 && m1 > m2) ||
		(y1 == y2 && m1 == m2 && d1 > d2)
}

func weekdayISO(t time.Time) int {
	wd := int(t.Weekday())
	if wd == 0 {
		return 7
	}
	return wd
}

func parseWeekdays(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, part := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n < 1 || n > 7 {
			return nil, fmt.Errorf("invalid weekday: %q", part)
		}
		result[n] = true
	}
	return result, nil
}

func parseMonthDays(s string) ([]int, error) {
	result := make([]int, 0)
	for _, part := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n == 0 || n > 31 || n < -2 {
			return nil, fmt.Errorf("invalid month day: %q", part)
		}
		result = append(result, n)
	}
	return result, nil
}

func parseMonths(s string) (map[int]bool, error) {
	result := make(map[int]bool)
	for _, part := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n < 1 || n > 12 {
			return nil, fmt.Errorf("invalid month: %q", part)
		}
		result[n] = true
	}
	return result, nil
}

func matchMonthDay(d time.Time, days []int) bool {
	lastDay := time.Date(d.Year(), d.Month()+1, 0, 0, 0, 0, 0, d.Location()).Day()
	for _, day := range days {
		if day > 0 && d.Day() == day {
			return true
		}
		if day < 0 && d.Day() == lastDay+day+1 {
			return true
		}
	}
	return false
}
