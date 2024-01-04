package data

import (
	"fmt"
	"strings"
)

// Checks if the error is a sqlite unique constraint violation.
func isUniqueError(err error, table string, fields ...string) bool {
	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		parts = append(parts, fmt.Sprintf("%s.%s", table, field))
	}
	msg := strings.Join(parts, ", ")
	return strings.Contains(err.Error(), msg)
}
