package helpers

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var defaultTimezone = "Asia/Bangkok" // fallback if not set

// SetDefaultTimezone sets the global default timezone.
// Call this once during app startup using cfg.Server.Timezone.
func SetDefaultTimezone(tz string) {
	if tz != "" {
		defaultTimezone = tz
	}
}

// GetTimeLocation returns the time.Location for the given timezone.
// If tzOverride is provided, it uses that; otherwise it uses the default
// timezone (which can be set via SetDefaultTimezone, or "Asia/Bangkok" as fallback).
// If loading from IANA fails, it falls back to a fixed UTC+7 zone.
func GetTimeLocation(tzOverride ...string) *time.Location {
	tz := defaultTimezone // ✅ use the variable, not the function
	if len(tzOverride) > 0 && tzOverride[0] != "" {
		tz = tzOverride[0]
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		// fallback to fixed UTC+7
		return time.FixedZone("Bangkok", 7*60*60)
	}
	return loc
}

func GetCurrentDatenow() string {
	return time.Now().Format("2006-01-02")
}

func GetCurrentTimenow() string {
	return time.Now().Format("15:04")
}

func GetCurrentTime() string {
	return GetCurrentTimenow()
}

// GetCurrentFullDatenow returns the current date and time in Asia/Bangkok timezone (UTC+7)
// formatted as "2006-01-02 15:04:05"
func GetCurrentFullDatenow() string {
	return TimeConvertermas(time.Now())
}

func NowDateTime() string {
	return GetCurrentFullDatenow()
}
func ConvertTZ(t time.Time, tzString string) time.Time {
	return t.In(GetTimeLocation(tzString))
}

func TimeConvertermas(t time.Time) string {
	return t.In(GetTimeLocation()).Format("2006-01-02 15:04:05")
}

// LocalTime formats a time.Time as "2006-01-02 15:04:05" in the configured
// timezone when marshaled to JSON (same output as TimeConvertermas). Swap
// response presenter fields to this type to get locale-consistent timestamps.
// ใช้สำหรับ field CreatedAt/UpdatedAt ใน presenter ให้คืนเวลา local ตาม timezone
type LocalTime struct {
	time.Time
}

// NewLocalTime wraps a time.Time into a LocalTime, normalizing the instant
// into the configured local timezone (universal time → Thai time).
// รับเวลาสากล (UTC) แล้วแปลงเป็นเวลาไทยทันที
func NewLocalTime(t time.Time) LocalTime {
	return LocalTime{Time: t.In(GetTimeLocation())}
} // GetTimeLocation

// NewLocalTimePtr wraps a *time.Time into a *LocalTime (nil-safe).
func NewLocalTimePtr(t *time.Time) *LocalTime {
	if t == nil {
		return nil
	}
	lt := NewLocalTime(*t)
	return &lt
}

// MarshalJSON implements json.Marshaler.
func (t LocalTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + TimeConvertermas(t.Time) + `"`), nil
}

// UnmarshalJSON accepts "2006-01-02 15:04:05" (local) or RFC3339.
func (t *LocalTime) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(string(b))
	s = strings.Trim(s, `"`)
	if s == "" || s == "null" {
		return nil
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", s, GetTimeLocation())
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	*t = LocalTime{Time: parsed}
	return nil
}

// String returns the formatted local time ("2006-01-02 15:04:05"), empty when zero.
func (t LocalTime) String() string {
	if t.IsZero() {
		return ""
	}
	return TimeConvertermas(t.Time)
}

func DiffMinutes(current, log time.Time) int {
	diff := current.Sub(log)
	return int(diff.Minutes())
}

func GetRandomString(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	return string(b)
}

func GetRandomInt(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		b[i] = digits[n.Int64()]
	}
	return string(b)
}

func ToThaiDate(t time.Time) string {
	monthNames := []string{"ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	year := t.Year() + 543
	month := monthNames[t.Month()-1]
	day := t.Day()
	hour := t.Format("15")
	minute := t.Format("04")
	second := t.Format("05")
	return StringConcat(day, " ", month, " ", year, " ", hour, ":", minute, ":", second, " น.")
}

func ToEnDate(t time.Time) string {
	monthNamesLong := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	year := t.Year()
	month := monthNamesLong[t.Month()-1]
	day := t.Day()
	hour := t.Format("15")
	minute := t.Format("04")
	second := t.Format("05")
	return StringConcat(day, " ", month, " ", year, " ", hour, ":", minute, ":", second)
}

func GetDayname() string {
	return strings.ToLower(time.Now().Weekday().String())
}

func GetDaynameall() string {
	return time.Now().Format("Monday, January 2, 2006")
}

func CheckEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func PasswordValidator(password string) bool {
	re := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*]).{8,}$`)
	return re.MatchString(password)
}

func GeneratePassword(length int) string {
	if length < 8 {
		length = 8
	}
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digit := "0123456789"
	special := "!@#$%^&*"
	all := lower + upper + digit + special

	var pass []byte
	pass = append(pass, lower[randInt(len(lower))])
	pass = append(pass, upper[randInt(len(upper))])
	pass = append(pass, digit[randInt(len(digit))])
	pass = append(pass, special[randInt(len(special))])

	for len(pass) < length {
		pass = append(pass, all[randInt(len(all))])
	}
	// shuffle
	for i := range pass {
		j := randInt(len(pass))
		pass[i], pass[j] = pass[j], pass[i]
	}
	return string(pass)
}

func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func StringConcat(parts ...interface{}) string {
	var sb strings.Builder
	for _, p := range parts {
		sb.WriteString(toString(p))
	}
	return sb.String()
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return ""
	}
}
