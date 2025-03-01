package nullish

import "database/sql"

// StringOrEmpty returns the string value of a sql.NullString if it is valid,
// otherwise it returns an empty string.
func StringOrEmpty(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

// ParseNullString returns a pointer to a string if the sql.NullString is valid,
// otherwise it returns nil.
func ParseNullString(v sql.NullString) *string {
	if v.Valid {
		return &v.String
	}
	return nil
}

// CreateNullString returns a sql.NullString from a pointer to a string.
func CreateNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}
