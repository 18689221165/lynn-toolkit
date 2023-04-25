package types

import (
	"database/sql/driver"
	"encoding/json"
)

type ExtraMap map[string]string

func (emap ExtraMap) Value() (driver.Value, error) {
	if emap == nil {
		return "{}", nil
	}
	return json.Marshal(emap)
}

func (emap *ExtraMap) Scan(v interface{}) error {
	switch v.(type) {
	case []byte:
		err := json.Unmarshal(v.([]byte), &emap)
		return err
	}
	return nil
}
