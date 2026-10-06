package http

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// parsePageParams รองรับ page / per_page / page_size / limit
func parsePageParams(q url.Values) (page, perPage int) {
	get := func(k string) string {
		if v, ok := q[k]; ok && len(v) > 0 {
			return v[0]
		}
		return ""
	}
	page, _ = strconv.Atoi(get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ = strconv.Atoi(get("per_page"))
	if perPage < 1 {
		perPage, _ = strconv.Atoi(get("page_size"))
	}
	if perPage < 1 {
		perPage, _ = strconv.Atoi(get("limit"))
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return
}

// parseSortParam แยก "field-ASC" → (field, order)
func parseSortParam(v string) (field, order string) {
	if v == "" {
		return "", ""
	}
	if idx := strings.LastIndex(v, "-"); idx > 0 {
		return v[:idx], v[idx+1:]
	}
	return v, ""
}

// extractClientMeta ดึง IP / Device / Browser
func extractClientMeta(r *http.Request) (ip, device, browser string) {
	ip = r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if i := strings.IndexByte(ip, ','); i > 0 {
			ip = strings.TrimSpace(ip[:i])
		}
	}
	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip = host
		} else {
			ip = r.RemoteAddr
		}
	}
	browser = r.UserAgent()
	device = r.Header.Get("X-Device-Access")
	if device == "" {
		device = detectDevice(browser)
	}
	return
}

func detectDevice(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "iphone"), strings.Contains(ua, "ipad"):
		return "iOS"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "mac os"), strings.Contains(ua, "macintosh"):
		return "macOS"
	case strings.Contains(ua, "linux"):
		return "Linux"
	case strings.Contains(ua, "curl"), strings.Contains(ua, "postman"):
		return "API-Client"
	default:
		return "Unknown"
	}
}

// ============================================================================
// Small numeric aliases — ใช้ในไฟล์อื่นได้โดยไม่ต้อง import strconv
// ============================================================================
func strconvAtoi(s string) (int, error) {
	return strconv.Atoi(s)
}

func strconvAtoi16(s string) (int16, error) {
	n, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(n), nil
}

// ============================================================================
// Pointer helpers — ใช้ใน handlers ที่ต้อง deref pointer ก่อนส่งให้ presenter
// ============================================================================
func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// ============================================================================
// Pointer helpers — แปลง pointer ↔ primitive สำหรับ map model/presenter
// ============================================================================
func derefInt16(p *int16) int16 {
	if p == nil {
		return 0
	}
	return *p
}

func ptrStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func ptrInt(n int) *int {
	return &n
}

func ptrIntFromInt16(n int16) *int {
	v := int(n)
	return &v
}

func derefInt16FromIntPtr(p *int) int16 {
	if p == nil {
		return 0
	}
	return int16(*p)
}