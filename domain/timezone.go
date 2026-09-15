package domain

import (
	"errors"
	"fmt"
	"time"
)

type TimeZone string

// NewTimeZone creates a new valid TimeZone.
// It returns an error if s is unknown timezone.
func NewTimeZone(s string) (TimeZone, error) {
	switch s {
	case "":
		// time.LoadLocation resolves this as UTC, ie as valid.
		return "", errors.New("time zone must not be empty")
	case "Local":
		// Resolves to the host's zone, which is not a durable identity.
		return "", errors.New(`time zone "Local" is not permitted`)
	}

	if _, err := time.LoadLocation(s); err != nil {
		return "", fmt.Errorf("invalid timezone %q: %w", s, err)
	}
	return TimeZone(s), nil
}
