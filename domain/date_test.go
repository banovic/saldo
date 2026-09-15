package domain

import (
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	testCases := []struct {
		name    string
		y       int
		m       time.Month
		d       int
		want    Date
		wantErr bool
	}{
		{"typical date", 2026, time.September, 15, Date{2026, time.September, 15}, false},
		{"first day of year", 2026, time.January, 1, Date{2026, time.January, 1}, false},
		{"last day of year", 2026, time.December, 31, Date{2026, time.December, 31}, false},

		{"lowest supported date", 1, time.January, 1, Date{1, time.January, 1}, false},
		{"highest supported date", 9999, time.December, 31, Date{9999, time.December, 31}, false},
		{"year zero", 0, time.January, 1, Date{}, true},
		{"negative year", -1, time.January, 1, Date{}, true},
		{"year past 9999", 10000, time.January, 1, Date{}, true},

		{"month zero", 2026, 0, 1, Date{}, true},
		{"month 13", 2026, 13, 1, Date{}, true},
		{"negative month", 2026, -1, 1, Date{}, true},

		{"day zero", 2026, time.January, 0, Date{}, true},
		{"negative day", 2026, time.January, -1, Date{}, true},
		{"day 32", 2026, time.January, 32, Date{}, true},
		{"day 31 in 30-day month", 2026, time.April, 31, Date{}, true},
		{"day 30 in 30-day month", 2026, time.April, 30, Date{2026, time.April, 30}, false},

		{"Feb 28 in common year", 2026, time.February, 28, Date{2026, time.February, 28}, false},
		{"Feb 29 in common year", 2026, time.February, 29, Date{}, true},
		{"Feb 29 in leap year", 2024, time.February, 29, Date{2024, time.February, 29}, false},
		{"Feb 30 in leap year", 2024, time.February, 30, Date{}, true},
		{"Feb 29 in century year is not leap", 1900, time.February, 29, Date{}, true},
		{"Feb 29 in year divisible by 400 is leap", 2000, time.February, 29, Date{2000, time.February, 29}, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewDate(tc.y, tc.m, tc.d)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewDate(%d, %d, %d) error = %v, want error %t", tc.y, tc.m, tc.d, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("NewDate(%d, %d, %d) = %v, want %v", tc.y, tc.m, tc.d, got, tc.want)
			}
		})
	}
}
