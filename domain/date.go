package domain

import (
	"fmt"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

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
