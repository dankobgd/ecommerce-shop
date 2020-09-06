package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
)

// TshirtVariant is the tshirt properties
type TshirtVariant struct {
	Color string `json:"color"`
	Size  string `json:"size"`
}

// ToJSON converts the tshirt properties to json string
func (v *TshirtVariant) ToJSON() string {
	b, _ := json.Marshal(v)
	return string(b)
}

// TshirtVariantFromJSON decodes the input and return the tshirt properties
func TshirtVariantFromJSON(data io.Reader) (*TshirtVariant, error) {
	var v *TshirtVariant
	err := json.NewDecoder(data).Decode(&v)
	return v, err
}

// Scan implements the scanner
func (v *TshirtVariant) Scan(val interface{}) error {
	b, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("Unsupported type: %T", v)
	}
	return json.Unmarshal(b, &v)
}

// Value implements the valuer
func (v *TshirtVariant) Value() (driver.Value, error) {
	return json.Marshal(v)
}
