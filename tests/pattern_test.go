package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"govtech/pkg/utilities/patterns"
)

func TestPattern(t *testing.T) {
	t.Run("email regexp", EmailRegexp)
	t.Run("notification field regexp", NotificationRegexp)
}

// Test for email regexp.
func EmailRegexp(t *testing.T) {

	// Valid email format.

	v_email_1 := "test@test.com"
	v_email_2 := "test12@gmail.com"
	v_email_3 := "11test@test.com"
	v_email_4 := "T3st@test.com"
	v_email_5 := "test@test.org"
	v_email_6 := "test12@gmail.org"
	v_email_7 := "11test@test.org"
	v_email_8 := "T3st@test.org"

	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, v_email_1))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, v_email_2))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, v_email_3))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, v_email_4))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, v_email_5))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, v_email_6))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, v_email_7))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, v_email_8))

	// Invalid email format.

	i_email_1 := "!*($#*)@!(*@!.!@&@!"
	i_email_2 := "@#*&#@@#"
	i_email_3 := "test@gmailcom"
	i_email_4 := "testgmail.com"

	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, i_email_1))
	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, i_email_2))
	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, i_email_3))
	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_EMAIL, i_email_4))
}

// Test for notification field regexp.
func NotificationRegexp(t *testing.T) {

	// Valid notification field format.

	v_notification_1 := "hello world!"
	v_notification_2 := "hello world! "
	v_notification_3 := "hello world"
	v_notification_4 := "hello world! @tagged@gmail.com"
	v_notification_5 := "hello world! @tagged@gmail.com @tagged2@gmail.com"

	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, v_notification_1))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, v_notification_2))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, v_notification_3))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, v_notification_4))
	assert.Equal(t, true, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, v_notification_5))

	// Inalid notification field format.

	i_notification_1 := "hello world @"
	i_notification_2 := "@tagged@gmail.com hello world @tagged@gmail.com"
	i_notification_3 := "hello world @tagged@gmail.com extra"
	i_notification_4 := "hello world @wrongformat.com"
	i_notification_5 := "hello world @ @tagged@gmail.com"

	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, i_notification_1))
	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, i_notification_2))
	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, i_notification_3))
	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, i_notification_4))
	assert.Equal(t, false, patterns.ValidatePattern(patterns.REGEX_PATTERN_NOTIFICATION, i_notification_5))
}
