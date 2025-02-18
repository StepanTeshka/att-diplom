package helper

import (
	"database/sql"
	"strconv"
)

// Helper для sql.NullString
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// Helper для sql.NullInt64
func NewNullInt64(i string) sql.NullInt64 {
	if i == "" {
		return sql.NullInt64{Valid: false}
	}
	val, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: val, Valid: true}
}
