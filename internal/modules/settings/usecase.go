package settings

import (
	"context"

	"icmongolang/internal/modules/settings/presenter"
)

// SettingsUseCaseI exposes the business operations of the settings module.
// One method per endpoint family; thin implementations delegate to the
// generic list/CRUD helpers plus the special flows (MqttData, email).
type SettingsUseCaseI interface {
	ListPaginate(ctx context.Context, spec ListSpec, q ListQuery) (*presenter.ListResult, error)
	RowsAll(ctx context.Context, spec ListSpec, q ListQuery) ([]map[string]interface{}, error)
	CreateRow(ctx context.Context, table string, m map[string]interface{}) error
	GetWhere(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExistsWhere(ctx context.Context, table string, query string, args ...interface{}) (bool, error)
	CountWhere(ctx context.Context, table string, query string, args ...interface{}) (int64, error)
	UpdateFields(ctx context.Context, table string, idCol string, idVal interface{}, fields map[string]interface{}) (int64, error)
	// SetStatusExclusive ports the NestJS update_*_status rule: when activating
	// (status == 1) every other row of the table is first reset to status 0
	// (with updateddate), then the target row is set; resetAll=false keeps the
	// target-only behaviour used by update_mqtt_status.
	SetStatusExclusive(ctx context.Context, table string, idCol string, id interface{}, status int, resetAll bool) (int64, error)
	DeleteRow(ctx context.Context, table string, query string, args ...interface{}) (int64, error)

	// Special flows
	MqttData(ctx context.Context, topic string) (*presenter.MqttDataResult, error)
	SendEmailStub(to, subject, content string) *presenter.SendEmailResult
	TestGmailConnection(ctx context.Context) *presenter.SendEmailResult
	SensorTypes(lang string) []map[string]string
}
