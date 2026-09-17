package domain

import (
	"errors"
	"testing"
)

func TestNewTimeZone(t *testing.T) {
	testCases := []struct {
		name    string
		s       string
		want    TimeZone
		wantErr error
	}{
		{"IANA zone", "Europe/Belgrade", "Europe/Belgrade", nil},
		{"IANA zone with non-hour offset", "Asia/Kolkata", "Asia/Kolkata", nil},
		{"IANA zone with nested name", "America/Argentina/Buenos_Aires", "America/Argentina/Buenos_Aires", nil},
		{"UTC", "UTC", "UTC", nil},

		{"empty is rejected, although LoadLocation treats it as UTC", "", "", ErrInvalidTimeZone},
		{"Local is rejected, it depends on the host", "Local", "", ErrInvalidTimeZone},

		{"unknown zone", "Europe/Nowhere", "", ErrInvalidTimeZone},
		{"offset is not a zone name", "+01:00", "", ErrInvalidTimeZone},
		{"whitespace is not trimmed", " Europe/Belgrade", "", ErrInvalidTimeZone},
		{"absolute path", "/Europe/Belgrade", "", ErrInvalidTimeZone},
		{"path traversal", "../Europe/Belgrade", "", ErrInvalidTimeZone},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewTimeZone(tc.s)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("NewTimeZone(%q) error = %v, want %v", tc.s, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("NewTimeZone(%q) = %q, want %q", tc.s, got, tc.want)
			}
		})
	}
}
