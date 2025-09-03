package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type DifficultyAttributes map[uint32]map[string]float64

func (d *DifficultyAttributes) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan DifficultyAttributes: %v", value)
	}
	return json.Unmarshal(bytes, d)
}

func (d DifficultyAttributes) Value() (driver.Value, error) {
	return json.Marshal(d)
}
