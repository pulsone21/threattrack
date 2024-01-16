package entities

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Entity interface {
	Incident | IncidentType | Task | Worklog | User
}

type Scanable interface {
	ScanTo(ScanFunc) error
}

type NullString struct {
	sql.NullString
}

func ToNullString(s string) NullString {
	return NullString{
		sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) ToString() (string, error) {
	if !ns.Valid {
		return "", fmt.Errorf("NullString is null")
	}
	return ns.String, nil
}

type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func ToNullInt(i int64) NullInt64 {
	return NullInt64{
		sql.NullInt64{
			Int64: i,
			Valid: true,
		},
	}
}

func (ni *NullInt64) ToInt64() (int64, error) {
	if !ni.Valid {
		return -1, fmt.Errorf("NullInt64 is null")
	}
	return ni.Int64, nil
}

type NullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

func ToNullBool(b bool) NullBool {
	return NullBool{
		sql.NullBool{
			Bool:  b,
			Valid: true,
		},
	}
}

func (nb *NullBool) ToBool() (bool, error) {
	if !nb.Valid {
		return false, fmt.Errorf("NullBool is null")
	}
	return nb.Bool, nil
}

type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

func ToNullFloat(f float64) NullFloat64 {
	return NullFloat64{
		sql.NullFloat64{
			Float64: f,
			Valid:   true,
		},
	}
}

func (nf *NullFloat64) ToFloat64() (float64, error) {
	if !nf.Valid {
		return -1, fmt.Errorf("NullFloat64 is null")
	}
	return nf.Float64, nil
}
