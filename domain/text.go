package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// validateText returns error if string is not valid.
// It returns bare error, wrapping it in sentinel error is responsibility of caller.
// String is valid if:
//   - valid utf8
//   - is not empty
//   - no leading nor trailing whitespaces
//   - no control chars
//   - length (in runes, valid utf8 chars) between [minLength, maxLength] - inclusive
func validateText(s string, minLength, maxLength int) error {
	if s == "" {
		return errors.New("empty string")
	}
	if !utf8.ValidString(s) {
		return errors.New("invalid utf8")
	}
	if s != strings.TrimSpace(s) {
		return errors.New("leading or trailing whitespace")
	}
	rc := utf8.RuneCountInString(s)
	if rc < minLength || rc > maxLength {
		return fmt.Errorf("min %d, max %d, got: %d", minLength, maxLength, rc)
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return fmt.Errorf("invalid char (control): %U", r)
		}
	}
	return nil
}
