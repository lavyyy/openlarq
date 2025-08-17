package util

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(bytes)
}
