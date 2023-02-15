package patterns

import (
	"regexp"
)

const REGEX_PATTERN_NOTIFICATION = `^([a-zA-Z0-9_.,!?-]+\s?)*(@[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\s*)*$`
const REGEX_PATTERN_EMAIL = `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`

// Validates a given string based on the given regexp pattern.
func ValidatePattern(pattern string, str string) bool {
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(str)
}
