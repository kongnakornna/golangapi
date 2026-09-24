package report

import (
	"fmt"
	"math"
	"strings"
)

var thaiNum = []string{"", "หนึ่ง", "สอง", "สาม", "สี่", "ห้า", "หก", "เจ็ด", "แปด", "เก้า"}
var thaiPos = []string{"", "สิบ", "ร้อย", "พัน", "หมื่น", "แสน", "ล้าน"}

func ThaiDigit(n int) string {
	if n == 0 {
		return ""
	}
	parts := []string{}
	unit := 0
	for n > 0 {
		digit := n % 10
		if digit != 0 {
			pos := thaiPos[unit]
			if unit == 0 && digit == 1 {
				parts = append([]string{"เอ็ด"}, parts...)
			} else if unit == 1 && digit == 2 {
				parts = append([]string{"ยี่สิบ"}, parts...)
			} else if unit == 1 && digit == 1 {
				parts = append([]string{"สิบ"}, parts...)
			} else {
				part := thaiNum[digit] + pos
				parts = append([]string{part}, parts...)
			}
		}
		n /= 10
		unit++
	}
	return strings.Join(parts, "")
}

func BahtThai(amount float64) string {
	if amount == 0 {
		return "ศูนย์บาทถ้วน"
	}

	intPart := int64(math.Floor(amount))
	satang := int64(math.Round((amount - float64(intPart)) * 100))

	var words string

	if intPart > 0 {
		million := intPart / 1000000
		remainder := intPart % 1000000

		if million > 0 {
			words += ThaiDigit(int(million)) + "ล้าน"
		}
		words += ThaiDigit(int(remainder))
		words += "บาท"
	}

	if satang > 0 {
		words += ThaiDigit(int(satang)) + "สตางค์"
	} else {
		words += "ถ้วน"
	}

	return words
}

var _ = fmt.Sprintf
