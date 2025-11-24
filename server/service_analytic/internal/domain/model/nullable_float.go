package model

import (
	"fmt"
	"reflect"

	"github.com/gocql/gocql"
)

// NullableFloat64 is a custom type to handle unmarshaling from ScyllaDB's float.
type NullableFloat64 struct {
	Float *float64
}

// UnmarshalCQL implements the gocql.Unmarshaler interface.
func (n *NullableFloat64) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	if data == nil {
		n.Float = nil
		return nil
	}

	var tempFloat32 float32
	if err := gocql.Unmarshal(info, data, &tempFloat32); err == nil {
		converted := float64(tempFloat32)
		n.Float = &converted
		return nil
	}

	var tempFloat64 float64
	if err := gocql.Unmarshal(info, data, &tempFloat64); err == nil {
		n.Float = &tempFloat64
		return nil
	}

	return fmt.Errorf("failed to unmarshal float: data is not a float32 or float64, but %s", reflect.TypeOf(data).String())
}

// MarshalCQL implements the gocql.Marshaler interface.
func (n NullableFloat64) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	if n.Float == nil {
		return nil, nil
	}
	// Marshal as float32 to match ScyllaDB float type
	return gocql.Marshal(info, float32(*n.Float))
}

// NullableFloat32 is a custom type for non-pointer float32 values from ScyllaDB.
type NullableFloat32 struct {
	Float float64
}

// UnmarshalCQL implements the gocql.Unmarshaler interface.
func (n *NullableFloat32) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	if data == nil {
		n.Float = 0.0
		return nil
	}

	var tempFloat32 float32
	if err := gocql.Unmarshal(info, data, &tempFloat32); err == nil {
		n.Float = float64(tempFloat32)
		return nil
	}

	var tempFloat64 float64
	if err := gocql.Unmarshal(info, data, &tempFloat64); err == nil {
		n.Float = tempFloat64
		return nil
	}

	return fmt.Errorf("failed to unmarshal float: data is not a float32, but %s", reflect.TypeOf(data).String())
}

// MarshalCQL implements the gocql.Marshaler interface.
func (n NullableFloat32) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return gocql.Marshal(info, float32(n.Float))
}
