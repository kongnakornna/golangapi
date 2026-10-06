package users

import (
	"context"
	"time"

	usermodels "icmongolang/internal/modules/users/models"

	"github.com/google/uuid"
)

// ============================================================
// CACHE
// ============================================================
type CacheRepository interface {
	GetJSON(ctx context.Context, key string, dst interface{}) (bool, error)
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeleteByPattern(ctx context.Context, pattern string) error
	ScanAndDelete(ctx context.Context, pattern string, batchSize int) (int, error)
}

// ============================================================
// SHARED
// ============================================================
type SortItem struct {
	ID       int
	Sort     int
	ParentID *int
}

// ============================================================
// ACCESS MENU
// ============================================================
type AccessMenuFilter struct {
	UserTypeID *int
	RoleID     *int
	MenuID     *int
	ParentID   *int
	Limit      int
	Offset     int
}

type AccessMenuRepository interface {
	Create(ctx context.Context, m *usermodels.SdUserAccessMenu) (*usermodels.SdUserAccessMenu, error)
	Get(ctx context.Context, id int) (*usermodels.SdUserAccessMenu, error)
	Update(ctx context.Context, m *usermodels.SdUserAccessMenu, values map[string]interface{}) (*usermodels.SdUserAccessMenu, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, f AccessMenuFilter) ([]*usermodels.SdUserAccessMenu, int64, error)
	GetMenusByRole(ctx context.Context, roleID int) ([]*usermodels.SdUserAccessMenu, error)
}

type AccessMenuUseCaseI interface {
	CreateAccessMenu(ctx context.Context, m *usermodels.SdUserAccessMenu) (*usermodels.SdUserAccessMenu, error)
	GetAccessMenu(ctx context.Context, id int) (*usermodels.SdUserAccessMenu, error)
	UpdateAccessMenu(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdUserAccessMenu, error)
	DeleteAccessMenu(ctx context.Context, id int) error
	ListAccessMenus(ctx context.Context, f AccessMenuFilter) ([]*usermodels.SdUserAccessMenu, int64, error)
}

// ============================================================
// FILE
// ============================================================
type FileFilter struct {
	UID        *string
	FileTypeID *int
	Statuses   []int16
	Limit      int
	Offset     int
}

type FileRepository interface {
	Create(ctx context.Context, f *usermodels.SdUserFile) (*usermodels.SdUserFile, error)
	Get(ctx context.Context, id int) (*usermodels.SdUserFile, error)
	Update(ctx context.Context, f *usermodels.SdUserFile, values map[string]interface{}) (*usermodels.SdUserFile, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, f FileFilter) ([]*usermodels.SdUserFile, int64, error)
}

type FileUseCaseI interface {
	CreateFile(ctx context.Context, f *usermodels.SdUserFile) (*usermodels.SdUserFile, error)
	GetFile(ctx context.Context, id int) (*usermodels.SdUserFile, error)
	UpdateFile(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdUserFile, error)
	DeleteFile(ctx context.Context, id int) error
	ListFiles(ctx context.Context, f FileFilter) ([]*usermodels.SdUserFile, int64, error)
}

// ============================================================
// LOG
// ============================================================
type LogFilter struct {
	UID       *uuid.UUID
	LogTypeID *int
	IPAddress *string
	From      *time.Time
	To        *time.Time
	Limit     int
	Offset    int
}

type LogRepository interface {
	Create(ctx context.Context, l *usermodels.SdUserLog) (*usermodels.SdUserLog, error)
	Get(ctx context.Context, id int, logTypeID int) (*usermodels.SdUserLog, error)
	List(ctx context.Context, f LogFilter) ([]*usermodels.SdUserLog, int64, error)
	Delete(ctx context.Context, id int, logTypeID int) error
}

type LogUseCaseI interface {
	CreateLog(ctx context.Context, l *usermodels.SdUserLog) (*usermodels.SdUserLog, error)
	GetLog(ctx context.Context, id int, logTypeID int) (*usermodels.SdUserLog, error)
	ListLogs(ctx context.Context, f LogFilter) ([]*usermodels.SdUserLog, int64, error)
	DeleteLog(ctx context.Context, id int, logTypeID int) error
	WriteAuditLogWithMeta(ctx context.Context, uid uuid.UUID, logTypeID int, name, detail, ip, device, browser string)
}

// ============================================================
// LOG TYPE
// ============================================================
type LogTypeRepository interface {
	Create(ctx context.Context, lt *usermodels.SdUserLogType) (*usermodels.SdUserLogType, error)
	Get(ctx context.Context, id int) (*usermodels.SdUserLogType, error)
	Update(ctx context.Context, lt *usermodels.SdUserLogType, values map[string]interface{}) (*usermodels.SdUserLogType, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*usermodels.SdUserLogType, int64, error)
}

type LogTypeUseCaseI interface {
	CreateLogType(ctx context.Context, lt *usermodels.SdUserLogType) (*usermodels.SdUserLogType, error)
	GetLogType(ctx context.Context, id int) (*usermodels.SdUserLogType, error)
	UpdateLogType(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdUserLogType, error)
	DeleteLogType(ctx context.Context, id int) error
	ListLogTypes(ctx context.Context, limit, offset int) ([]*usermodels.SdUserLogType, int64, error)
}

// ============================================================
// MODULE
// ============================================================
type ModuleFilter struct {
	ModuleParentID *int
	Statuses       []int16
	Keyword        string
	Limit          int
	Offset         int
	SortBy         string
	SortOrder      string
}

type ModuleRepository interface {
	Create(ctx context.Context, m *usermodels.SdModules) (*usermodels.SdModules, error)
	Get(ctx context.Context, id int) (*usermodels.SdModules, error)
	Update(ctx context.Context, m *usermodels.SdModules, values map[string]interface{}) (*usermodels.SdModules, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, f ModuleFilter) ([]*usermodels.SdModules, int64, error)
	BulkUpdateSort(ctx context.Context, items []SortItem) error
}

type ModuleUseCaseI interface {
	CreateModule(ctx context.Context, m *usermodels.SdModules) (*usermodels.SdModules, error)
	GetModule(ctx context.Context, id int) (*usermodels.SdModules, error)
	UpdateModule(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdModules, error)
	DeleteModule(ctx context.Context, id int) error
	ListModules(ctx context.Context, f ModuleFilter) ([]*usermodels.SdModules, int64, error)
	BulkUpdateModuleSort(ctx context.Context, items []SortItem) error
}

// ============================================================
// MENU
// ============================================================
type MenuFilter struct {
	ModulesID    *int
	MenuParentID *int
	Statuses     []int16
	Keyword      string
	Limit        int
	Offset       int
	SortBy       string
	SortOrder    string
}

type MenuRepository interface {
	Create(ctx context.Context, m *usermodels.SdMenu) (*usermodels.SdMenu, error)
	Get(ctx context.Context, id int) (*usermodels.SdMenu, error)
	Update(ctx context.Context, m *usermodels.SdMenu, values map[string]interface{}) (*usermodels.SdMenu, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, f MenuFilter) ([]*usermodels.SdMenu, int64, error)
	BulkUpdateSort(ctx context.Context, items []SortItem) error
}

type MenuUseCaseI interface {
	CreateMenu(ctx context.Context, m *usermodels.SdMenu) (*usermodels.SdMenu, error)
	GetMenu(ctx context.Context, id int) (*usermodels.SdMenu, error)
	UpdateMenu(ctx context.Context, id int, values map[string]interface{}) (*usermodels.SdMenu, error)
	DeleteMenu(ctx context.Context, id int) error
	ListMenus(ctx context.Context, f MenuFilter) ([]*usermodels.SdMenu, int64, error)
	BulkUpdateMenuSort(ctx context.Context, items []SortItem) error
}