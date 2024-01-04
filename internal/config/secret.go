package config

import (
	"encoding/json"
)

// Secret is a string that should not be printed in logs.
type Secret string

func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal("*****")
}
