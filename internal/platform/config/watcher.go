package config

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/kongnakornna/golangapi/pkg/logger"
)

var log = logger.Default()

// ConfigWatcher ตัวเฝ้าติดตามการเปลี่ยนแปลงไฟล์คอนฟิก
type ConfigWatcher struct {
	mu        sync.RWMutex
	config    *AppConfig
	callbacks []func(*AppConfig)
	watcher   *fsnotify.Watcher
	stopCh    chan struct{}
}

// NewConfigWatcher สร้างตัวเฝ้าติดตามการเปลี่ยนแปลงไฟล์คอนฟิก
func NewConfigWatcher(configPath string) (*ConfigWatcher, error) {
	// โหลดคอนฟิกเริ่มต้น
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// สร้างตัวเฝ้าติดตามไฟล์
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// เพิ่มไฟล์คอนฟิกเข้าไปในรายการเฝ้าติดตาม
	if err := watcher.Add(configPath); err != nil {
		watcher.Close()
		return nil, err
	}

	cw := &ConfigWatcher{
		config:    cfg,
		callbacks: make([]func(*AppConfig), 0),
		watcher:   watcher,
		stopCh:    make(chan struct{}),
	}

	// เริ่มต้นการเฝ้าติดตาม
	go cw.watch(configPath)

	return cw, nil
}

// watch เฝ้าติดตามการเปลี่ยนแปลงของไฟล์คอนฟิก
func (cw *ConfigWatcher) watch(configPath string) {
	// ตัวจับเวลาสำหรับป้องกันการเรียกซ้ำถี่เกินไป (debounce)
	var debounceTimer *time.Timer
	debounceDuration := 100 * time.Millisecond

	for {
		select {
		case event, ok := <-cw.watcher.Events:
			if !ok {
				return
			}

			// จัดการเฉพาะเหตุการณ์เขียนและสร้างไฟล์
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				// ใช้ debounce เพื่อหลีกเลี่ยงการโหลดซ้ำบ่อยเกินไป
				if debounceTimer != nil {
					debounceTimer.Stop()
				}

				debounceTimer = time.AfterFunc(debounceDuration, func() {
					cw.reloadConfig(configPath)
				})
			}

		case err, ok := <-cw.watcher.Errors:
			if !ok {
				return
			}
			log.Error("เกิดข้อผิดพลาดในการเฝ้าติดตามไฟล์คอนฟิก", "error", err)

		case <-cw.stopCh:
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return
		}
	}
}

// reloadConfig โหลดคอนฟิกใหม่
func (cw *ConfigWatcher) reloadConfig(configPath string) {
	log.Info("ตรวจพบการเปลี่ยนแปลงไฟล์คอนฟิก กำลังโหลดคอนฟิกใหม่", "path", configPath)

	// อ่านคอนฟิกใหม่
	newCfg, err := LoadConfig(configPath)
	if err != nil {
		log.Error("โหลดคอนฟิกใหม่ล้มเหลว", "error", err)
		return
	}

	// ตรวจสอบความถูกต้องของคอนฟิกใหม่
	if err := cw.validateConfig(newCfg); err != nil {
		log.Error("ตรวจสอบคอนฟิกไม่ผ่าน", "error", err)
		return
	}

	// อัปเดตคอนฟิก
	cw.mu.Lock()
	oldCfg := cw.config
	cw.config = newCfg
	cw.mu.Unlock()

	log.Info("โหลดคอนฟิกใหม่สำเร็จ")

	// แจ้งเตือน callback ทั้งหมด
	cw.notifyCallbacks(oldCfg, newCfg)
}

// validateConfig ตรวจสอบความถูกต้องของคอนฟิก
func (cw *ConfigWatcher) validateConfig(cfg *AppConfig) error {
	// ที่นี่สามารถเพิ่มตรรกะการตรวจสอบคอนฟิกได้
	// เช่น: ตรวจสอบฟิลด์ที่จำเป็น, รูปแบบข้อมูล เป็นต้น

	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return ErrInvalidPort
	}

	if cfg.Database.Host == "" {
		return ErrMissingDatabaseHost
	}

	return nil
}

// notifyCallbacks แจ้งเตือนฟังก์ชัน callback ทั้งหมด
func (cw *ConfigWatcher) notifyCallbacks(oldCfg, newCfg *AppConfig) {
	// บันทึกการเปลี่ยนแปลงคอนฟิก
	cw.logConfigChanges(oldCfg, newCfg)

	// เรียกใช้ callback
	for _, callback := range cw.callbacks {
		go func(cb func(*AppConfig)) {
			defer func() {
				if r := recover(); r != nil {
					log.Error("การทำงาน callback สำหรับการเปลี่ยนแปลงคอนฟิกล้มเหลว", "error", r)
				}
			}()
			cb(newCfg)
		}(callback)
	}
}

// logConfigChanges บันทึกการเปลี่ยนแปลงหลักของคอนฟิก
func (cw *ConfigWatcher) logConfigChanges(oldCfg, newCfg *AppConfig) {
	// บันทึกการเปลี่ยนแปลงคอนฟิกที่สำคัญ
	if oldCfg.Server.Port != newCfg.Server.Port {
		log.Info("พอร์ตเซิร์ฟเวอร์มีการเปลี่ยนแปลง", "old", oldCfg.Server.Port, "new", newCfg.Server.Port)
	}

	if oldCfg.Log.Level != newCfg.Log.Level {
		log.Info("ระดับการบันทึก Log มีการเปลี่ยนแปลง", "old", oldCfg.Log.Level, "new", newCfg.Log.Level)
	}

	if oldCfg.Database.MaxOpenConns != newCfg.Database.MaxOpenConns {
		log.Info("จำนวนการเชื่อมต่อฐานข้อมูลสูงสุดมีการเปลี่ยนแปลง", "old", oldCfg.Database.MaxOpenConns, "new", newCfg.Database.MaxOpenConns)
	}
}

// GetConfig ดึงคอนฟิกปัจจุบัน
func (cw *ConfigWatcher) GetConfig() *AppConfig {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	return cw.config
}

// OnConfigChange ลงทะเบียน callback สำหรับการเปลี่ยนแปลงคอนฟิก
func (cw *ConfigWatcher) OnConfigChange(callback func(*AppConfig)) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.callbacks = append(cw.callbacks, callback)
}

// Stop หยุดการเฝ้าติดตาม
func (cw *ConfigWatcher) Stop() error {
	close(cw.stopCh)
	return cw.watcher.Close()
}

// WatchConfig เฝ้าติดตามการเปลี่ยนแปลงไฟล์คอนฟิก (ใช้ Viper)
func WatchConfig(onChange func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("ไฟล์คอนฟิกมีการเปลี่ยนแปลง", "file", e.Name)
		if onChange != nil {
			onChange()
		}
	})
}

// HotReloadConfig ตัวอย่างคอนฟิกรองรับการโหลดขณะรันไทม์ (Hot Reload)
type HotReloadConfig struct {
	mu     sync.RWMutex
	config *AppConfig
}

// NewHotReloadConfig สร้างคอนฟิกที่รองรับ Hot Reload
func NewHotReloadConfig(configPath string) (*HotReloadConfig, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	hrc := &HotReloadConfig{
		config: cfg,
	}

	// ตั้งค่า Viper ให้เฝ้าติดตาม
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("ไฟล์คอนฟิกเปลี่ยนแปลง กำลังโหลดใหม่", "file", e.Name)

		// แยกวิเคราะห์คอนฟิกใหม่
		newCfg := &AppConfig{}
		if err := viper.Unmarshal(newCfg); err != nil {
			log.Error("แยกวิเคราะห์คอนฟิกใหม่ล้มเหลว", "error", err)
			return
		}

		// อัปเดตคอนฟิก
		hrc.mu.Lock()
		hrc.config = newCfg
		hrc.mu.Unlock()

		log.Info("โหลดคอนฟิกใหม่ขณะรันไทม์สำเร็จ")
	})

	return hrc, nil
}

// Get ดึงคอนฟิก (ป้องกันปัญหาข้อมูลขัดแย้งจากการเข้าถึงพร้อมกัน)
func (hrc *HotReloadConfig) Get() *AppConfig {
	hrc.mu.RLock()
	defer hrc.mu.RUnlock()
	return hrc.config
}
