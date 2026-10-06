package command

import "encoding/json"

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
