package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidDate = errors.New("invalid date")
)

// Date represents date (year, month, day) without location nor timezone.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// NewDate constructs new date.
func NewDate(y int, m time.Month, d int) (Date, error) {
	date := Date{Year: y, Month: m, Day: d}
	if err := date.Validate(); err != nil {
		return Date{}, err
	}
	return date, nil
}

// Validate returns error if date is not valid.
// Date is valid if:
//   - year is between 1 and 9999 inclusive
//   - month is 1..12 inclusive
//   - day is valid for given year and month
func (d Date) Validate() error {
	if d.Year < 1 || d.Year > 9999 {
		return fmt.Errorf("%w: year %d out of range [1, 9999]", ErrInvalidDate, d.Year)
	}
	if d.Month < time.January || d.Month > time.December {
		return fmt.Errorf("%w: invalid month %d", ErrInvalidDate, int(d.Month))
	}
	t := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
	if t.Year() != d.Year || t.Month() != d.Month || t.Day() != d.Day {
		return fmt.Errorf("%w: invalid date: %04d-%02d-%02d", ErrInvalidDate, d.Year, int(d.Month), d.Day)
	}
	return nil
}

// DateIn creates a new date which contains given time instant in given location.
func DateIn(t time.Time, loc *time.Location) (Date, error) {
	y, m, d := t.In(loc).Date()
	date, err := NewDate(y, m, d)
	if err != nil {
		return Date{}, err
	}
	return date, nil
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
