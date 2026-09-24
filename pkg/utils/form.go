package utils

import (
	"encoding/json"
	"net/http"
)

// DecodeFormToStruct populates a struct from form values.
// Works for simple fields (string, int, bool). For nested types, use JSON.
func DecodeFormToStruct(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	// Convert form values to a map and then JSON marshal/unmarshal
	formMap := make(map[string]interface{})
	for key, values := range r.Form {
		if len(values) == 1 {
			formMap[key] = values[0]
		} else {
			formMap[key] = values
		}
	}
	data, err := json.Marshal(formMap)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}
