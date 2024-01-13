package models

import (
	"database/sql/driver"
	"encoding/json"
)

type dString struct {
	Hash       string `json:"hash"`
	Identifier string `json:"identifier"`
}
type dStringArray []dString

func (sla *dStringArray) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &sla)
}
func (sla dStringArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}
