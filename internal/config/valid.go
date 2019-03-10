//
// valid.go
// Copyright (C) 2019 Grigorii Sokolik <g.sokol99@g-sokol.info>
//
// Distributed under terms of the MIT license.
//

package config

import (
	"encoding/json"
)

type ValidInt struct {
	Valid bool
	Value int
}

func (v *ValidInt) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &v.Value); err != nil {
		return err
	}
	v.Valid = true
	return nil
}

type ValidString struct {
	Valid bool   `ignored:"true"`
	Value string `ignored:"true"`
}

func (v *ValidString) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &v.Value); err != nil {
		return err
	}
	v.Valid = true
	return nil
}

func (v *ValidString) UnmarshalText(b []byte) error {
	v.Value = string(b)
	v.Valid = true
	return nil
}

type ValidBool struct {
	Valid bool
	Value bool
}

func (v *ValidBool) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &v.Value); err != nil {
		return err
	}
	v.Valid = true
	return nil
}

type ValidInt64 struct {
	Valid bool
	Value int64
}

func (v *ValidInt64) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &v.Value); err != nil {
		return err
	}
	v.Valid = true
	return nil
}
