// ฟังก์ชันยูทิลิตี้
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// GenerateRandomString สร้างสตริงแบบสุ่ม
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateRandomNumber สร้างตัวเลขสุ่ม
func GenerateRandomNumber(min, max int) int {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return min
	}
	return min + int(bigInt.Int64())
}
