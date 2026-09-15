package domain

import "testing"

func TestNewTimeZone(t *testing.T) {
	testCases := []struct {
		name    string
		s       string
		want    TimeZone
		wantErr bool
	}{
		{"IANA zone", "Europe/Belgrade", "Europe/Belgrade", false},
		{"IANA zone with non-hour offset", "Asia/Kolkata", "Asia/Kolkata", false},
		{"IANA zone with nested name", "America/Argentina/Buenos_Aires", "America/Argentina/Buenos_Aires", false},
		{"UTC", "UTC", "UTC", false},

		{"empty is rejected, although LoadLocation treats it as UTC", "", "", true},
		{"Local is rejected, it depends on the host", "Local", "", true},

		{"unknown zone", "Europe/Nowhere", "", true},
		{"offset is not a zone name", "+01:00", "", true},
		{"whitespace is not trimmed", " Europe/Belgrade", "", true},
		{"absolute path", "/Europe/Belgrade", "", true},
		{"path traversal", "../Europe/Belgrade", "", true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewTimeZone(tc.s)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewTimeZone(%q) error = %v, want error %t", tc.s, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("NewTimeZone(%q) = %q, want %q", tc.s, got, tc.want)
			}
		})
	}
}
