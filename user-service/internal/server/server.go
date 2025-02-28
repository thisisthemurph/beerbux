package server

import "database/sql"

func parseNullString(v sql.NullString) *string {
	if v.Valid {
		return &v.String
	}
	return nil
}

func createNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
}
