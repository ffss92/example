package data

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUniqueError(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		name     string
		table    string
		fields   []string
		err      error
		expected bool
	}{
		{
			name:     "should return true for a unique error (single field)",
			table:    "users",
			fields:   []string{"username"},
			err:      errors.New("UNIQUE constraint failed: users.username"),
			expected: true,
		},
		{
			name:     "should return false for other errors (single field)",
			table:    "users",
			fields:   []string{"username"},
			err:      errors.New("some random error"),
			expected: false,
		},
		{
			name:     "should return true for a unique error (multi field)",
			table:    "post_likes",
			fields:   []string{"user_id", "post_id"},
			err:      errors.New("UNIQUE constraint failed: post_likes.user_id, post_likes.post_id"),
			expected: true,
		},
		{
			name:     "should return false for other errors (multi field)",
			table:    "post_likes",
			fields:   []string{"user_id", "post_id"},
			err:      errors.New("some random error"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isUniqueError(tc.err, tc.table, tc.fields...)
			assert.Equal(tc.expected, result)
		})
	}
}
