package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidTimeZone = errors.New("invalid time zone")
)

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
func (tz TimeZone) Validate() error {
	switch tz {
	case "":
		// time.LoadLocation resolves this as UTC, ie as valid.
		return fmt.Errorf("%w: must not be empty", ErrInvalidTimeZone)
	case "Local":
		// Resolves to the host's zone, which is not a durable identity.
		return fmt.Errorf("%w: must not be %q", ErrInvalidTimeZone, tz)
	}

	if _, err := time.LoadLocation(string(tz)); err != nil {
		return fmt.Errorf("%w: %q", ErrInvalidTimeZone, tz)
	}

	return nil
}
