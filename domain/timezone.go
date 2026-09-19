package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidTimeZone = errors.New("invalid time zone")
)

// TimeZone is an IANA time zone identifier ("Europe/Belgrade", ...).
type TimeZone string

// NewTimeZone creates a new valid TimeZone.
// It returns an error if s is unknown timezone.
func NewTimeZone(s string) (TimeZone, error) {
	tz := TimeZone(s)
	if err := tz.Validate(); err != nil {
		return "", err
	}
	return tz, nil
}

// Validate returns error if timezone is not valid.
// Empty string ("") and "UTC" both resolve to "UTC" timezone.
// Timezone is valid if:
//   - not empty string; "" resolves to "UTC", but it represents zero for TimeZone, so forbidden here.
//   - not string 'Local'; 'Local' is valid timezone on current system, it is not durable.
//   - must exist in time zone database (as used by time.LoadLocation())
func (tz TimeZone) Validate() error {
	switch tz {
	case "":
		return fmt.Errorf("%w: must not be empty", ErrInvalidTimeZone)
	case "Local":
		return fmt.Errorf("%w: resolved to host time zone and is not durable %q", ErrInvalidTimeZone, tz)
	}

	if _, err := time.LoadLocation(string(tz)); err != nil {
		return fmt.Errorf("%w: %q: %w", ErrInvalidTimeZone, tz, err)
	}

	return nil
}

// Location returns location for timezone.
func (tz TimeZone) Location() (*time.Location, error) {
	if err := tz.Validate(); err != nil {
		return nil, err
	}
	return time.LoadLocation(string(tz))
}
