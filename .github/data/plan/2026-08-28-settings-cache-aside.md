# Plan: Cache-aside สำหรับ slow query ใน settings module

## ภาพรวมการแก้ไข

| # | Layer | ไฟล์ | สิ่งที่แก้ |
|---|-------|------|-----------|
| 1 | pkg | `pkg/db/redis/redis_conn.go` | เพิ่ม `Delete` เข้า `Cache` interface + implement ใน `RedisCache` |
| 2 | model | `internal/modules/settings/types.go` | เพิ่ม field `CacheKey` + `CacheTTL` ใน struct `ListSpec` + helper `CacheEnabled()` |
| 3 | model | `internal/modules/settings/specs_master.go` | ใส่ `CacheKey`/`CacheTTL` ให้ spec กลุ่ม master/dropdown |
| 4 | model | `internal/modules/settings/specs_integration.go` | ใส่ `CacheKey`/`CacheTTL` ให้ spec กลุ่ม integration/dropdown |
| 5 | model | `internal/modules/settings/specs_device_schedule.go` | ใส่ `CacheKey`/`CacheTTL` ให้ `ScheduleSpec` (device/schedule ที่เป็น config) |
| 6 | usecase | `internal/modules/settings/usecase/usecase.go` | cache-aside ใน `ListPaginate`/`RowsAll` + invalidate ค่า cache ตอน write (`CreateRow`/`UpdateFields`/`SetStatusExclusive`/`DeleteRow`) |

---

## บทนำ

### ปัญหา
`ListPaginate` / `RowsAll` ใน `internal/modules/settings/usecase/usecase.go` ส่ง query ไป Postgres ทุกครั้งที่เรียก
spec ที่หนัก เช่น `DeviceSpec`, `SensorSpec`, `MqttSpec` ต้องรัน COUNT + SELECT ต่อ `LEFT JOIN` 4–5 ตาราง
กับ `CASE` expression หลายจุด → การ์ดอ่านซ้ำๆ (dropdown, master list) ช้าและกิน DB จนเกินจำเป็น

### แนวทาง
ใช้ **cache-aside** (อ่าน cache ก่อน → ไม่มีค่อยไป DB → set กลับ cache) โดย
- cache เฉพาะ spec ที่อ่านบ่อย/เขียนน้อย (master/dropdown/config)
- มี TTL สั้น (default 60 วินาที) กันข้อมูลเก่า
- ตอน write สำเร็จ → invalidate cache ของตารางนั้นทันที
- error ใดๆ จาก cache → ตกกลับไปอ่าน DB แทน ไม่ทำให้ request พัง

---

## Section 1 — เพิ่ม `Delete` ใน `Cache` interface

### Location
`pkg/db/redis/redis_conn.go:102-105` (interface) และ `redis_conn.go:115-129` (`RedisCache`)

### From → to
ตอนนี้ `Cache` interface มีแค่ `Get`/`Set` → ต้องมี `Delete` ด้วย เพื่อเอาไว้ invalidate ตอน write
`RedisCache` ยังไม่มี method `Delete` → เพิ่ม method ที่เรียก `Del` ของ go-redis

### Change detail
เพิ่ม method ใน interface:

```go
type Cache interface {
	Get(ctx context.Context, key string, dst interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
```

เพิ่ม method ใน `RedisCache`:

```go
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
```

หมายเหตุ: `iot`/`mqtt` handler ใช้ `Cache` interface ของตัวเอง (structural type) → ไม่กระทบ

---

## Section 2 — เพิ่ม field `CacheKey` + `CacheTTL` ใน `ListSpec`

### Location
`internal/modules/settings/types.go:24-33`

### From → to
ตอนนี้ `ListSpec` ยังไม่มีข้อมูลบ่งชี้ว่า spec นี้ cache ได้ไหม
→ เพิ่ม `CacheKey string` (ว่าง = ไม่ cache) และ `CacheTTL time.Duration` (default 60s เมื่อ `CacheKey` ไม่ว่าง)
พร้อม method `CacheEnabled()` และ `CacheDuration()` เพื่อให้ usecase ใช้ได้ง่าย

### Change detail
เพิ่ม field ใน struct:

```go
type ListSpec struct {
	Table        string
	Selects      []string
	Joins        []string
	Filters      []Filter
	FixedFilters []FixedFilter
	ParamConds   []ParamCond
	SortCols     map[string]string
	DefaultSort  string
	// CacheKey ใช้เปิด cache-aside; ว่าง = ไม่ cache, มีค่า = prefix ของ redis key
	CacheKey string
	// CacheTTL อายุ cache; ถ้า <= 0 ให้ใช้ defaultTTL ใน usecase
	CacheTTL time.Duration
}
```

เพิ่ม method (ใส่ไว้ข้าง struct):

```go
const defaultCacheTTL = 60 * time.Second

func (s ListSpec) CacheEnabled() bool { return s.CacheKey != "" }

func (s ListSpec) CacheDuration() time.Duration {
	if s.CacheTTL > 0 {
		return s.CacheTTL
	}
	return defaultCacheTTL
}
```

ต้องเพิ่ม import `time` ใน `types.go`

---

## Section 3 — ใส่ `CacheKey` ให้ spec กลุ่ม master

### Location
`internal/modules/settings/specs_master.go`

### From → to
spec ต่อไปนี้ยังไม่มี `CacheKey` → ใส่ `CacheKey` + `CacheTTL`

### Change detail
- `SettingSpec` → `CacheKey: "setting"`
- `LocationSpec` → `CacheKey: "location"`
- `TypeSpec` → `CacheKey: "type"`
- `DeviceTypeSpec` → `CacheKey: "devicetype"`
- `GroupSpec` → `CacheKey: "group"`
- `SensorSpec` → `CacheKey: "sensor"`

ตัวอย่าง (แทรกท้าย struct `SettingSpec`):

```go
var SettingSpec = ListSpec{
	Table:   "sd_iot_setting s",
	Selects: []string{ /* อยู่แล้ว */ },
	Joins:   []string{ /* อยู่แล้ว */ },
	Filters: []Filter{ /* อยู่แล้ว */ },
	SortCols: map[string]string{ /* อยู่แล้ว */ },
	DefaultSort: "s.createddate ASC",
	CacheKey: "setting",
}
```

---

## Section 4 — ใส่ `CacheKey` ให้ spec กลุ่ม integration

### Location
`internal/modules/settings/specs_integration.go`

### From → to
spec ต่อไปนี้ยังไม่มี `CacheKey` → ใส่ `CacheKey` (ทุกตัวเป็น config/dropdown อ่านบ่อย)

### Change detail
- `MqttSpec` → `CacheKey: "mqtt"`
- `MqttHostSpec` → `CacheKey: "mqtthost"`
- `ApiSpec` → `CacheKey: "api"`
- `EmailSpec` → `CacheKey: "email"`
- `HostSpec` → `CacheKey: "host"`
- `InfluxdbSpec` → `CacheKey: "influxdb"`
- `LineSpec` → `CacheKey: "line"`
- `NoderedSpec` → `CacheKey: "nodered"`
- `SmsSpec` → `CacheKey: "sms"`
- `TokenSpec` → `CacheKey: "token"`
- `TelegramSpec` → `CacheKey: "telegram"`
- `DashboardConfigSpec` → `CacheKey: "dashboardconfig"`

---

## Section 5 — ใส่ `CacheKey` ให้ spec กลุ่ม device/schedule

### Location
`internal/modules/settings/specs_device_schedule.go`

### From → to
`ScheduleSpec` เป็น config ที่เขียนน้อย อ่านบ่อย → cache ได้

### Change detail
- `ScheduleSpec` → `CacheKey: "schedule"`

หมายเหตุ: `DeviceSpec`/`DeviceActiveSpec`/`DeviceAllSpec`/`DeviceAllActiveSpec`
และพวก alarm/process-log (ใน `specs_alarm_logs.go`) **ไม่** ใส่ cache
เพราะเป็นข้อมูลที่เปลี่ยนบ่อย/มีข้อมูลเยอะ → cache จะได้ประโยชน์น้อยและทำingest สกปรก

---

## Section 6 — cache-aside ใน usecase + invalidate ตอน write

### Location
`internal/modules/settings/usecase/usecase.go`

### From → to
ตอนนี้ `ListPaginate`/`RowsAll` เรียก `uc.repo.*` ตรงๆ ทุกครั้ง ส่วน write (`CreateRow`/`UpdateFields`/
`SetStatusExclusive`/`DeleteRow`) ไม่ได้แตะ cache เลย
→ ใส่ cache-aside ให้ list read และ invalidate ตอน write สำเร็จ

### Change detail

เพิ่ม `time` + `sort` import และ const/helper ด้านบน:

```go
const listCachePrefix = "settings:list"

// keyForList สร้าง redis key จาก spec.CacheKey + query (sorted filters เพื่อให้ deterministic)
func (uc *settingsUseCase) keyForList(spec settings.ListSpec, q settings.ListQuery) string {
	keys := make([]string, 0, len(q.values))
	for k := range q.values {
		if settings.FilterExcluded(k) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString(listCachePrefix)
	sb.WriteString(":")
	sb.WriteString(spec.CacheKey)
	for _, k := range keys {
		sb.WriteString(":")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(q.Value(k))
	}
	sb.WriteString(":page=")
	sb.WriteString(strconv.Itoa(q.Page))
	sb.WriteString(":size=")
	sb.WriteString(strconv.Itoa(q.PageSize))
	sb.WriteString(":sort=")
	sb.WriteString(q.Sort)
	return sb.String()
}

// listKeysForTable คืนทุก prefix pattern ที่ต้อง invalidate เมื่อเขียนตารางนั้น
func (uc *settingsUseCase) listCachePrefixesFor(table string) []string {
	out := []string{}
	if uc.cache == nil {
		return out
	}
	// เก็บ registry: table -> []CacheKey (ดูด้านล่าง)
	for _, ck := range cacheKeysByTable[strings.Fields(table)[0]] {
		out = append(out, listCachePrefix+":"+ck+":")
	}
	return out
}
```

**หมายเหตุเรื่อง registry:** ใน `types.go` (หรือไฟล์ spec) เพิ่ม map กลาง:

```go
// cacheKeysByTable map ชื่อตาราง (ไม่มี alias) -> list ของ CacheKey ที่อ้างถึงตารางนั้น
var cacheKeysByTable = map[string][]string{
	"sd_iot_setting":    {"setting"},
	"sd_iot_location":   {"location"},
	"sd_iot_type":       {"type"},
	"sd_iot_device_type": {"devicetype"},
	"sd_iot_group":      {"group"},
	"sd_iot_sensor":     {"sensor"},
	"sd_iot_mqtt":       {"mqtt"},
	"sd_mqtt_host":      {"mqtthost"},
	"sd_iot_api":        {"api"},
	"sd_iot_email":      {"email"},
	"sd_iot_host":       {"host"},
	"sd_iot_influxdb":   {"influxdb"},
	"sd_iot_line":       {"line"},
	"sd_iot_nodered":    {"nodered"},
	"sd_iot_sms":        {"sms"},
	"sd_iot_token":      {"token"},
	"sd_iot_telegram":   {"telegram"},
	"sd_dashboard_config": {"dashboardconfig"},
	"sd_iot_schedule":   {"schedule"},
}
```

**ListPaginate — เปลี่ยนจาก:**

```go
func (uc *settingsUseCase) ListPaginate(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) (*presenter.ListResult, error) {
	items, total, err := uc.repo.ListPaginate(ctx, spec, q)
	...
	return presenter.NewListResult(items, total, q.Page, q.PageSize, q.FilterMap()), nil
}
```

**เป็น:**

```go
func (uc *settingsUseCase) ListPaginate(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) (*presenter.ListResult, error) {
	var res *presenter.ListResult

	// cache-aside: อ่าน cache ก่อน
	cached := false
	if uc.cache != nil && spec.CacheEnabled() {
		key := uc.keyForList(spec, q)
		var cachedRes []map[string]interface{}
		if err := uc.cache.Get(ctx, key, &cachedRes); err == nil && cachedRes != nil {
			cached = true
			// ข้อมูล cached เก็บแค่ data + total (ต้องเก็บ meta เดิมไว้ใน cache ด้วย)
		} else if err != nil {
			uc.logger.Warnf("settings ListPaginate cache get %s: %v", key, err)
		}
	}

	if !cached {
		items, total, err := uc.repo.ListPaginate(ctx, spec, q)
		if err != nil {
			uc.logger.Errorf("settings ListPaginate table=%s: %v", spec.Table, err)
			return nil, err
		}
		res = presenter.NewListResult(items, total, q.Page, q.PageSize, q.FilterMap())
		if uc.cache != nil && spec.CacheEnabled() {
			if serr := uc.cache.Set(ctx, uc.keyForList(spec, q), res, spec.CacheDuration()); serr != nil {
				uc.logger.Warnf("settings ListPaginate cache set: %v", serr)
			}
		}
	}
	return res, nil
}
```

**สำคัญ:** เพราะ `presenter.NewListResult` สร้าง meta (total, totalPages) เอง จะ cache ค่าคืน `*presenter.ListResult`
ทั้งก้อน (data + meta) เพื่อให้ผลลัพธ์ตรงกันเป๊ะทั้ง hit/miss → วิธีง่ายสุดคือ `cache.Get` ลงไปใน
`*presenter.ListResult` (ซึ่ง `redis` marshal/unmarshal ของ `ListResult` ได้ถ้าทุก field เป็น JSON-safe)

> ตรวจสอบ `presenter.ListResult` ว่า field ทั้งหมดเป็น JSON-friendly (มี json tag) — ถ้าขาด ให้ add json tag
> ให้ครบก่อน เพื่อให้ unmarshal กลับมาใช้ได้

**RowsAll — เปลี่ยนจาก:**

```go
func (uc *settingsUseCase) RowsAll(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) ([]map[string]interface{}, error) {
	items, err := uc.repo.RowsAll(ctx, spec, q)
	...
	return items, nil
}
```

**เป็น:** cache-aside เดียวกัน — `cache.Get(ctx, key, &items)` → hit คืนเลย, miss เรียก repo แล้ว `cache.Set` TTL

**Write — invalidate:** ใน `CreateRow`/`UpdateFields`/`SetStatusExclusive`/`DeleteRow`
(และ `UpdateAllMap` ที่เรียกผ่าน `SetStatusExclusive`) หลังจากสำเร็จ เรียก:

```go
func (uc *settingsUseCase) invalidateListCache(ctx context.Context, table string) {
	if uc.cache == nil {
		return
	}
	dt := table
	if i := strings.IndexByte(dt, ' '); i >= 0 {
		dt = dt[:i] // ตัด alias ออก เช่น "sd_iot_setting s" -> "sd_iot_setting"
	}
	for _, prefix := range uc.listCachePrefixesFor(dt) {
		keys, err := uc.cache.Keys(ctx, prefix+"*")
		...
		for _, key := range keys {
			uc.cache.Delete(ctx, key)
		}
	}
}
```

**ปัญหา:** interface `Cache` ตอนนี้ไม่มี `Keys` → ต้องเพิ่ม `Keys(ctx, pattern string) ([]string, error)`
ใน interface + implement ใน `RedisCache` (เรียก `client.Keys`)

**ทางเลือกที่ simplify (แนะนำ):** แทนที่จะใช้ `Keys` + prefix (มี race/ค่าใช้จ่าย scan) ให้ใช้ **ตาราง index set ใน redis**:
- ตอน `Set` list อ่าน → `SAdd` ชื่อ key เข้า `settings:idx:{CacheKey}` พร้อม `Expire` ของ set
- ตอน write → `SMembers("settings:idx:{CacheKey}")` แล้ว `Del` ทุก key แล้ว `Del` set
แต่ interface `Cache` ต้องขยายเยอะ → เก็บแบบเรียบง่าย: **ใช้ `Keys` pattern** `settings:list:{CacheKey}:*`
แทน (ยอมรับ race เล็กน้อย — ถ้า key เกิดใหม่ระหว่าง scan อาจค้างได้หนึ่ง TTL ซึ่งcache จะ expire เอง)

**ตัดสินใจ:** ใช้ `Keys` pattern ง่ายสุด + TTL เป็นตัวกันข้อมูลเก่า กรณี race
→ ต้องเพิ่ม method `Keys` ใน interface + `RedisCache`

---

## Unit Tests

**เขียน**

ครอบใน `internal/modules/settings/usecase/usecase_test.go`:
- `TestListPaginateCacheHit`: cache เก็บผล → `ListPaginate` คืนของจาก cache โดยไม่เรียก repo (`fakeRepo.listItems` ควรไม่ถูกแตะ หรือตรวจผ่าน spy)
- `TestListPaginateCacheMissPopulates`: cache miss → เรียก repo → `Set` ถูกเรียกด้วย key ที่ถูกต้อง + TTL
- `TestRowsAllCacheHit` / `TestRowsAllCacheMissPopulates`: ทำแบบเดียวกับข้างบน
- `TestInvalidateOnWrite`: `CreateRow`/`UpdateFields`/`DeleteRow` สำเร็จ → เรียก `Delete`/`Keys` สำหรับ prefix ของตารางนั้น
- อัปเดต fakes `cacheHit`/`cacheMiss` ให้ implement method ใหม่ (`Delete`, `Keys`) ที่เพิ่มใน interface
- เพิ่ม fake cache แบบมี spy (`map[string][]byte` เป็น backstore) เพื่อ assert การ get/set/delete

---

## Files touched

1. `pkg/db/redis/redis_conn.go`
2. `internal/modules/settings/types.go`
3. `internal/modules/settings/specs_master.go`
4. `internal/modules/settings/specs_integration.go`
5. `internal/modules/settings/specs_device_schedule.go`
6. `internal/modules/settings/usecase/usecase.go`
7. `internal/modules/settings/usecase/usecase_test.go`

## หมายเหตุ / ความเสี่ยง

- TTL default 60s → staleness สูงสุด ≤ 60s หลัง write ถ้า request แทรกมาระหว่าง write กับ invalidate — ยอมรับได้สำหรับข้อมูล config/master
- error จาก cache ทุกชนิด → ตกกลับไปอ่าน DB ไม่พัง request
- `filterExcludedKeys` มี `deletecache` อยู่แล้ว (ไม่ถูก echo ใน filter) — ยังไม่ wire force-refresh เว้นแต่ต้องการเพิ่มทีหลัง
- อ่านที่ cache layer ต้องทำให้ `presenter.ListResult` marshal/unmarshal ได้ครบ (มี json tag) ก่อน
