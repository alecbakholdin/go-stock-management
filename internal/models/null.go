package models

import "database/sql"

func NullStringIfZero(v string) sql.NullString {
	return sql.NullString{String: v, Valid: v != ""}
}

// null string which is null if v==z
func NullStringIfMatch(v, z string) sql.NullString {
	if v == z {
		return sql.NullString{}
	}
	return sql.NullString{String: v, Valid: true}
}

func NullInt32IfZero(v int32) sql.NullInt32 {
	return sql.NullInt32{Int32: v, Valid: v != 0}
}

func NullFloat64IfZero(v float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: v, Valid: v != 0}
}