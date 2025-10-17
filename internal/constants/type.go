package constants

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// //////////////////////////////////////////////
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	return json.Marshal(j)
}

func (j *JSONMap) Scan(src interface{}) error {
	if src == nil {
		*j = nil
		return nil
	}
	var source []byte
	switch v := src.(type) {
	case string:
		source = []byte(v)
	case []byte:
		source = v
	default:
		return errors.New("incompatible type for JSONMap")
	}

	return json.Unmarshal(source, j)
}
func (j JSONMap) MustJSON() []byte {
	b, _ := json.Marshal(j)
	return b
}

// //////////////////////////////////////////