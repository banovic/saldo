package domain

import (
	"fmt"
	"time"
)

// Date represents date (year, month, day) without location nor timezone.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// NewDate constructs new date.
func NewDate(y int, m time.Month, d int) (Date, error) {
	if y < 1 || y > 9999 {
		return Date{}, fmt.Errorf("year %d out of range [1, 9999]", y)
	}
	if m < time.January || m > time.December {
		return Date{}, fmt.Errorf("invalid month %d", int(m))
	}
	t := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	if t.Year() != y || t.Month() != m || t.Day() != d {
		return Date{}, fmt.Errorf("invalid date: %04d-%02d-%02d", y, int(m), d)
	}
	return Date{Year: y, Month: m, Day: d}, nil
}

// DateIn creates a new date which contains given time instant in given location.
func DateIn(t time.Time, loc *time.Location) Date {
	y, m, d := t.In(loc).Date()
	return Date{Year: y, Month: m, Day: d}
}

// Before checks if a date was before given date.
func (d Date) Before(o Date) bool {
	if d.Year != o.Year {
		return d.Year < o.Year
	}
	if d.Month != o.Month {
		return d.Month < o.Month
	}
	return d.Day < o.Day
}

// After checks if a date was after given date.
func (d Date) After(o Date) bool {
	return o.Before(d)
}

// String returns date as ISO 8601 (YYYY-MM-DD) string.
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, int(d.Month), d.Day)
}
