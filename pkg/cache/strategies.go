package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// โหมดแคช แบบ Cache-Aside (โหมดแคชที่ใช้บ่อยที่สุด)
type CacheAside struct {
	cache  Cache
	loader DataLoader
	ttl    time.Duration
}

// DataLoader ตัวโหลดข้อมูล
type DataLoader func(ctx context.Context, key string) (interface{}, error)

// NewCacheAside สร้างแคชในโหมด Cache-Aside
func NewCacheAside(cache Cache, loader DataLoader, ttl time.Duration) *CacheAside {
	return &CacheAside{
		cache:  cache,
		loader: loader,
		ttl:    ttl,
	}
}

// ดึงข้อมูล (ตรวจสอบแคชก่อน หากไม่พบในแคช ให้โหลดจากแหล่งข้อมูล)
func (ca *CacheAside) Get(ctx context.Context, key string, dest interface{}) error {
	// ดึงข้อมูลจากแคชก่อน
	err := ca.cache.GetObject(ctx, key, dest)
	if err == nil {
		return nil
	}

	// แคชไม่พบข้อมูล โหลดข้อมูลจากแหล่งข้อมูล
	data, err := ca.loader(ctx, key)
	if err != nil {
		return err
	}

	// เขียนข้อมูลลงแคช (แบบอะซิงโครนัส เพื่อหลีกเลี่ยงการบล็อก)
	go ca.cache.SetObject(context.Background(), key, data, ca.ttl)

	// คัดลอกข้อมูลไปยังเป้าหมาย
	return copyValue(data, dest)
}

// Invalidate ล้างแคช
func (ca *CacheAside) Invalidate(ctx context.Context, key string) error {
	return ca.cache.Delete(ctx, key)
}

// SingleFlight ป้องกันการโอเวอร์โหลดแคช (อนุญาตให้โหลดข้อมูลได้ครั้งละหนึ่งคำขอเท่านั้น)
type SingleFlight struct {
	cache   Cache
	loader  DataLoader
	ttl     time.Duration
	flights map[string]*flightGroup
	mu      sync.Mutex
}

type flightGroup struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// NewSingleFlight จะสร้างแคช SingleFlight ขึ้นมา
func NewSingleFlight(cache Cache, loader DataLoader, ttl time.Duration) *SingleFlight {
	return &SingleFlight{
		cache:   cache,
		loader:  loader,
		ttl:     ttl,
		flights: make(map[string]*flightGroup),
	}
}

// Get ดึงข้อมูล (เพื่อป้องกันแคชเสียหาย)
func (sf *SingleFlight) Get(ctx context.Context, key string, dest interface{}) error {
	// ดึงข้อมูลจากแคชก่อน
	err := sf.cache.GetObject(ctx, key, dest)
	if err == nil {
		return nil
	}

	// ตรวจสอบว่ามีกระบวนการโหลดใด ๆ กำลังดำเนินการอยู่หรือไม่
	sf.mu.Lock()
	if fg, ok := sf.flights[key]; ok {
		sf.mu.Unlock()
		// กำลังรอให้การโหลดเสร็จสมบูรณ์
		fg.wg.Wait()
		if fg.err != nil {
			return fg.err
		}
		return copyValue(fg.val, dest)
	}

	// สร้าง  flight group
	fg := &flightGroup{}
	fg.wg.Add(1)
	sf.flights[key] = fg
	sf.mu.Unlock()

	// กำลังโหลดข้อมูล
	fg.val, fg.err = sf.loader(ctx, key)
	if fg.err == nil {
		// เขียนลงแคช
		sf.cache.SetObject(ctx, key, fg.val, sf.ttl)
		copyValue(fg.val, dest)
	}

	// เสร็จสมบูรณ์
	fg.wg.Done()

	// การทำความสะอาด flight group
	sf.mu.Lock()
	delete(sf.flights, key)
	sf.mu.Unlock()

	return fg.err
}

// copyValue การคัดลอกค่า (การแปลงข้อมูลเป็น JSON และการแปลงกลับจาก JSON)
func copyValue(src, dest interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
