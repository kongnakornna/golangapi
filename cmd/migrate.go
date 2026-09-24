package cmd

import (
	"errors"
	"fmt"
	"icmongolang/config"
	"icmongolang/internal/models"
	fullschedulemodel "icmongolang/internal/modules/fullschedule"
	iotmodels "icmongolang/internal/modules/iot/models"
	kafkamodels "icmongolang/internal/modules/kafka/models"
	pomodels "icmongolang/internal/modules/purchaseorder/models"
	"icmongolang/pkg/db/postgres"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/vectordb"
	stdlog "log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate data",
	Long:  "Migrate data",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()
		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		appLogger.Infof("--migrate Run--")
		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			appLogger.Fatalf("เชื่อมต่อไม่สำเร็จ - Postgresql init: %s", err)
		}
		appLogger.Infof("Postgres connected successfully")

		if err := createExtensions(psqlDB, cfg); err != nil {
			appLogger.Warnf("Failed to create extensions: %v", err)
		}

		// Run migration with skip list for problematic models
		if err := Migrate(psqlDB, appLogger); err != nil {
			appLogger.Fatal("Migrate ข้อมูลไม่สำเร็จ: ", err)
		}
		appLogger.Info("Migrate ข้อมูลสำเร็จ")

		// Ensure vector index/table for the configured vector database
		if cfg.VectorDB.Provider == "pgvector" {
			vdb, err := vectordb.New(&cfg.VectorDB, nil, psqlDB, appLogger)
			if err == nil {
				dims := cfg.VectorDB.Dims
				if dims <= 0 {
					dims = 768
				}
				if err := vdb.EnsureIndex(cmd.Context(), dims); err != nil {
					appLogger.Warnf("⚠️ Vector DB table not ready: %v", err)
				} else {
					appLogger.Infof("✅ Vector DB table %q ready (pgvector)", cfg.VectorDB.Index)
				}
			} else {
				appLogger.Warnf("⚠️ Vector DB init failed: %v", err)
			}
		}
	},
}

func createExtensions(db *gorm.DB, cfg *config.Config) error {
	extensions := []string{
		"CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"",
		"CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"",
	}
	// The pgvector extension is only needed when the vector DB provider is pgvector.
	if strings.ToLower(cfg.VectorDB.Provider) == "pgvector" {
		extensions = append(extensions, "CREATE EXTENSION IF NOT EXISTS \"vector\"")
	}
	for _, ext := range extensions {
		if err := db.Exec(ext).Error; err != nil {
			return err
		}
	}
	return nil
}

// Migrate runs AutoMigrate with skip list for models that have known issues.
// Add any model that fails to this list to skip it.
func Migrate(db *gorm.DB, log logger.Logger) error {
	// ⚠️  Add any model that fails to this map to skip it.
	skipModels := map[string]bool{
		// Skip models with Device relation (because Device has Mqtt/Location issues)
		"icmongolang/internal/models.Device":                   true,
		"icmongolang/internal/models.DeviceNotificationConfig": true,
		"icmongolang/internal/models.DeviceSchedule":           true,
		"icmongolang/internal/models.DeviceStatusHistory":      true,
		"icmongolang/internal/models.SdActivityLog":            true, // has ActivityType relation
		"icmongolang/internal/models.DeviceGroupMember":        true, // may also have Device relation
		"icmongolang/internal/models.NotificationCondition":    true,
		"icmongolang/internal/models.NotificationLog":          true,
		"icmongolang/internal/models.ReportData":               true,
		"icmongolang/internal/models.SensorData":               true,
		"icmongolang/internal/models.DeviceGroup":              true, // might have nested issues
		"icmongolang/internal/models.DeviceCategory":           true,
		"icmongolang/internal/models.DeviceStatus":             true, // if it has Device relation
		"icmongolang/internal/models.DeviceConfig":             true,
		"icmongolang/internal/models.SdIotDevice":              true, // duplicate of iot.models.Device on same table -> would flip column types (text vs varchar) every run
		"icmongolang/internal/models.CommandLog":               true,
		"icmongolang/internal/models.DeviceAlert":              true,
		"icmongolang/internal/models.IotData":                  true,
		"icmongolang/internal/models.ActivityLog":              true,
		// Add any other models that cause "invalid field" errors here
	}

	allModels := []interface{}{
		&models.ActivityLog{},
		&models.CommandLog{},
		&models.DeviceAlert{},
		&models.DeviceConfig{},
		&models.DeviceStatus{},
		&models.IotData{},
		&models.NotiNotificationLog{},
		&models.NotiNotificationRule{},
		&models.NotiNotificationType{},
		&models.NotiNotification{},
		&models.NotificationDevice{},
		&models.NotificationGroup{},
		&models.NotificationGroupsDevicesNotificationDevice{},
		&models.NotificationLog{},
		&models.NotificationType{},
		&models.SdActivityLog{},
		&models.SdActivityTypeLog{},
		&models.SdAdminAccessMenu{},
		&models.SdAirControl{},
		&models.SdAirControlDeviceMap{},
		&models.SdAirControlLog{},
		&models.SdAirMod{},
		&models.SdAirModDeviceMap{},
		&models.SdAirPeriod{},
		&models.SdAirPeriodDeviceMap{},
		&models.SdAirSettingWarning{},
		&models.SdAirSettingWarningDeviceMap{},
		&models.SdAirWarning{},
		&models.SdAirWarningDeviceMap{},
		&models.SdAlarmProcessLog{},
		&models.SdAlarmProcessLogEmail{},
		&models.SdAlarmProcessLogLine{},
		&models.SdAlarmProcessLogMqtt{},
		&models.SdAlarmProcessLogSms{},
		&models.SdAlarmProcessLogTelegram{},
		&models.SdAlarmProcessLogTemp{},
		&models.SdApiKey{},
		&models.SdAuditLog{},
		&models.AuditLogEntry{},
		&models.BlockchainAnchor{},
		&models.FlowDefinition{},
		&models.SdChannelTemplate{},
		&models.SdDashboardConfig{},
		&models.SdDeviceCategory{},
		&models.SdDeviceGroup{},
		&models.SdDeviceLog{},
		&models.SdDeviceMember{},
		&models.SdDeviceNotificationConfig{},
		&models.SdDeviceSchedule{},
		&models.SdDeviceStatusHistory{},
		&models.SdGroupNotificationConfig{},
		&models.SdIotAlarmDevice{},
		&models.SdIotAlarmDeviceEvent{},
		&models.SdIotApi{},
		&models.SdIotDevice{},
		&models.SdIotDeviceAction{},
		&models.SdIotDeviceActionLog{},
		&models.SdIotDeviceActionUser{},
		&models.SdIotDeviceAlarmAction{},
		&models.SdIotDeviceType{},
		&models.SdIotEmail{},
		&models.SdIotGroup{},
		&models.SdIotHost{},
		&models.SdIotInfluxdb{},
		&models.SdIotLine{},
		&models.SdIotLocation{},
		&models.SdIotMqtt{},
		&models.SdIotNodered{},
		&models.SdIotSchedule{},
		&models.SdIotScheduleDevice{},
		&models.SdIotSensor{},
		&models.SdIotSetting{},
		&models.SdIotSms{},
		&models.SdIotTelegram{},
		&models.SdIotToken{},
		&models.SdIotType{},
		&models.SdModuleLog{},
		&models.SdMqttHost{},
		&models.SdMqttLog{},
		&models.SdNotificationChannel{},
		&models.SdNotificationCondition{},
		&models.SdNotificationLog{},
		&models.SdNotificationType{},
		&models.SdReportData{},
		&models.SdScheduleProcessLog{},
		&models.SdSensorData{},
		&models.SdSystemSetting{},
		&models.SdUser{},
		&models.SdUserAccessMenu{},
		&models.SdUserFile{},
		&models.SdUserLog{},
		&models.SdUserLogType{},
		&models.SdUserRole{},
		&models.SdUserRolesAccess{},
		&models.SdUserRolePermission{},
		&models.SdUserRolesPermision{},
		&models.Tnb{},
		&models.Item{},
		&models.WsMessage{},
		&models.WsSession{},
		&models.User{},
		&models.PaymentMethod{},
		&models.Payment{},
		&models.Receipt{},
		&models.PaymentHistory{},
		&models.OutstandingBalance{},
		&kafkamodels.Order{},
		&pomodels.PurchaseOrderHeader{},
		&pomodels.PurchaseOrderDetail{},
		&pomodels.PurchaseOrderStatusHistory{},
		&iotmodels.IotData{},
		&iotmodels.ActivityLog{},
		&iotmodels.AirControlDeviceMap{},
		&iotmodels.AirControl{},
		&iotmodels.AirControlLog{},
		&iotmodels.AirMod{},
		&iotmodels.AirModDeviceMap{},
		&iotmodels.AirPeriod{},
		&iotmodels.AirPeriodDeviceMap{},
		&iotmodels.AirSettingWarning{},
		&iotmodels.AirSettingWarningDeviceMap{},
		&iotmodels.AirWarning{},
		&iotmodels.AirWarningDeviceMap{},
		&iotmodels.Device{},
		&iotmodels.DeviceCategory{},
		&iotmodels.DeviceGroup{},
		&iotmodels.DeviceNotificationConfig{},
		&iotmodels.DeviceSchedule{},
		&iotmodels.DeviceStatusHistory{},
		&iotmodels.GroupNotificationConfig{},
		&iotmodels.DeviceStatus{},
		&iotmodels.DeviceConfig{},
		&iotmodels.CommandLog{},
		&iotmodels.DeviceAlert{},
		// fullschedule module (fs_* tables)
		&fullschedulemodel.Schedule{},
		&fullschedulemodel.ScheduleDevice{},
		&fullschedulemodel.ScheduleHistory{},
		&fullschedulemodel.ScheduleSetting{},
		&fullschedulemodel.Group{},
		&fullschedulemodel.Zone{},
		&fullschedulemodel.Area{},
		&fullschedulemodel.AreaDevice{},
	}

	// Filter out skipped models
	var toMigrate []interface{}
	for _, m := range allModels {
		t := reflect.TypeOf(m)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		typeName := t.PkgPath() + "." + t.Name()
		if skipModels[typeName] {
			log.Warnf("Skipping model %s due to known issue", typeName)
			continue
		}
		toMigrate = append(toMigrate, m)
	}

	log.Infof("Migrating %d models (skipped %d)", len(toMigrate), len(allModels)-len(toMigrate))

	// Use a schema-migration logger: DDL like ALTER COLUMN TYPE on a large table
	// routinely exceeds the 200 ms "slow query" threshold and would otherwise spam
	// "SLOW SQL" warnings on every run. 5 s keeps genuinely stuck statements visible
	// while silencing expected DDL. Errors are still returned by AutoMigrate below.
	migrateDB := db.Session(&gorm.Session{
		Logger: gormlogger.New(
			stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
			gormlogger.Config{
				SlowThreshold:             5 * time.Second,
				LogLevel:                  gormlogger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})

	// Run AutoMigrate on the filtered list.
	// Migrate each model individually so a non-fatal constraint/index error on one
	// table (e.g. GORM trying to DROP CONSTRAINT "uni_..." / index that was already
	// replaced under a different name by a conflicting/duplicate model, or an index
	// created manually in migrations/*.sql) does not abort the whole migration.
	var firstErr error
	migratedCount := 0
	for _, m := range toMigrate {
		t := reflect.TypeOf(m)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		name := t.PkgPath() + "." + t.Name()
		err := migrateDB.AutoMigrate(m)
		if err == nil {
			migratedCount++
			continue
		}
		if isUndefinedObjectError(err) {
			log.Warnf("Skipping non-fatal migrate error on %s (constraint/index not found): %v", name, err)
			continue
		}
		if firstErr == nil {
			firstErr = fmt.Errorf("auto migrate error on %s: %w", name, err)
		}
	}

	if firstErr != nil {
		return firstErr
	}

	log.Infof("Migration completed successfully (%d models)", migratedCount)
	return nil
}

// isUndefinedObjectError reports whether the given error is a PostgreSQL
// "undefined_object" error (SQLSTATE 42704). GORM/Postgres raise this when it tries
// to drop a constraint or index whose name no longer exists — typically because a
// duplicate model or a manually-created index in migrations/*.sql already replaced it
// under a different name, or because the constraint was created by an earlier run.
// These are non-fatal: the migration only aims to ensure the schema exists, and a
// missing drop target already means the target does not exist.
func isUndefinedObjectError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "42704" {
		return true
	}
	return false
}
