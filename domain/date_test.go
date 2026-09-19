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

func TestDateIn(t *testing.T) {
	belgrade, err := time.LoadLocation("Europe/Belgrade")
	if err != nil {
		t.Fatal(err)
	}
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name string
		t    time.Time
		loc  *time.Location
		want Date
	}{
		{"UTC", time.Date(2026, time.September, 18, 23, 30, 0, 0, time.UTC), time.UTC, Date{2026, time.September, 18}},
		{"east of UTC, next day", time.Date(2026, time.September, 18, 23, 30, 0, 0, time.UTC), belgrade, Date{2026, time.September, 19}},
		{"west of UTC, previous day", time.Date(2026, time.September, 19, 1, 0, 0, 0, time.UTC), newYork, Date{2026, time.September, 18}},
		{"east of UTC, same day", time.Date(2026, time.September, 18, 12, 0, 0, 0, time.UTC), tokyo, Date{2026, time.September, 18}},
		{"input location is ignored", time.Date(2026, time.September, 19, 1, 0, 0, 0, belgrade), time.UTC, Date{2026, time.September, 18}},
		{"year boundary", time.Date(2026, time.December, 31, 23, 0, 0, 0, time.UTC), belgrade, Date{2027, time.January, 1}},
		// Belgrade is UTC+1 in winter, UTC+2 in summer.
		{"winter offset", time.Date(2026, time.January, 15, 23, 30, 0, 0, time.UTC), belgrade, Date{2026, time.January, 16}},
		{"winter offset, before midnight", time.Date(2026, time.January, 15, 22, 30, 0, 0, time.UTC), belgrade, Date{2026, time.January, 15}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := DateIn(tc.t, tc.loc); got != tc.want {
				t.Errorf("DateIn(%v, %v) = %v, want %v", tc.t, tc.loc, got, tc.want)
			}
		})
	}
}

func TestDateBeforeAfter(t *testing.T) {
	testCases := []struct {
		name       string
		d, o       Date
		wantBefore bool
		wantAfter  bool
	}{
		{"equal", Date{2026, time.September, 18}, Date{2026, time.September, 18}, false, false},
		{"earlier day", Date{2026, time.September, 17}, Date{2026, time.September, 18}, true, false},
		{"later day", Date{2026, time.September, 19}, Date{2026, time.September, 18}, false, true},
		{"earlier month, later day", Date{2026, time.August, 31}, Date{2026, time.September, 1}, true, false},
		{"later month, earlier day", Date{2026, time.October, 1}, Date{2026, time.September, 30}, false, true},
		{"earlier year, later month and day", Date{2025, time.December, 31}, Date{2026, time.January, 1}, true, false},
		{"later year, earlier month and day", Date{2027, time.January, 1}, Date{2026, time.December, 31}, false, true},
		{"zero value is before any date", Date{}, Date{1, time.January, 1}, true, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.d.Before(tc.o); got != tc.wantBefore {
				t.Errorf("%v.Before(%v) = %t, want %t", tc.d, tc.o, got, tc.wantBefore)
			}
			if got := tc.d.After(tc.o); got != tc.wantAfter {
				t.Errorf("%v.After(%v) = %t, want %t", tc.d, tc.o, got, tc.wantAfter)
			}
		})
	}
}

func TestDateString(t *testing.T) {
	testCases := []struct {
		name string
		d    Date
		want string
	}{
		{"typical date", Date{2026, time.September, 18}, "2026-09-18"},
		{"single digit month and day are padded", Date{2026, time.January, 5}, "2026-01-05"},
		{"small year is padded", Date{1, time.January, 1}, "0001-01-01"},
		{"highest supported date", Date{9999, time.December, 31}, "9999-12-31"},
		{"zero value", Date{}, "0000-00-00"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.d.String(); got != tc.want {
				t.Errorf("%#v.String() = %q, want %q", tc.d, got, tc.want)
			}
		})
	}
}
