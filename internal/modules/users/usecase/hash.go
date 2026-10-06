package usecase

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

func hashFilter(f interface{}) string {
	b, _ := json.Marshal(f)
	sum := sha1.Sum(b)
	return hex.EncodeToString(sum[:])
}

func strPtr(s string) *string {
	if s == "" { return nil }
	return &s
}