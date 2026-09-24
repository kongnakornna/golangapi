package transaction

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- fake database/sql driver: ทำให้ gorm Begin/Commit/Rollback ทำงานได้ไม่ต้องมี DB จริง ---

type noopConnector struct{}

func (noopConnector) Connect(context.Context) (driver.Conn, error) { return noopConn{}, nil }
func (noopConnector) Driver() driver.Driver                        { return noopDriver{} }

type noopDriver struct{}

func (noopDriver) Open(string) (driver.Conn, error) { return noopConn{}, nil }

type noopConn struct{}

func (noopConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("not implemented") }
func (noopConn) Close() error                        { return nil }
func (noopConn) Begin() (driver.Tx, error)           { return noopTx{}, nil }

type noopTx struct{}

func (noopTx) Commit() error   { return nil }
func (noopTx) Rollback() error { return nil }

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sql.OpenDB(noopConnector{}),
	}), &gorm.Config{})
	assert.NoError(t, err)
	return gdb
}

var errBiz = errors.New("biz failure")

func TestExecuteWithOptions_FnPanics_ReturnsError(t *testing.T) {
	m := NewGormTransactionManager(newTestDB(t))

	err := m.ExecuteWithOptions(context.Background(), nil, func(ctx context.Context, tx *gorm.DB) error {
		panic("boom")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic in transaction")
	assert.Contains(t, err.Error(), "boom")
}

func TestExecuteWithOptions_Success_Commits(t *testing.T) {
	m := NewGormTransactionManager(newTestDB(t))

	err := m.ExecuteWithOptions(context.Background(), nil, func(ctx context.Context, tx *gorm.DB) error {
		return nil
	})

	assert.NoError(t, err)
}

func TestExecuteWithOptions_FnError_RollsBack(t *testing.T) {
	m := NewGormTransactionManager(newTestDB(t))

	err := m.ExecuteWithOptions(context.Background(), nil, func(ctx context.Context, tx *gorm.DB) error {
		return errBiz
	})

	assert.True(t, errors.Is(err, errBiz))
}

func TestExecute_Delegates_PanicBecomesError(t *testing.T) {
	m := NewGormTransactionManager(newTestDB(t))

	err := m.Execute(context.Background(), func(ctx context.Context, tx *gorm.DB) error {
		panic(errors.New("wrapped boom"))
	})

	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "panic in transaction"))
}

func TestRunInTransaction_FnPanics_ReturnsError(t *testing.T) {
	db := newTestDB(t)

	err := RunInTransaction(db, func(tx *gorm.DB) error {
		panic("boom")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic in transaction")
	assert.Contains(t, err.Error(), "boom")
}

func TestRunInTransaction_FnError_RollsBack(t *testing.T) {
	db := newTestDB(t)

	err := RunInTransaction(db, func(tx *gorm.DB) error {
		return errBiz
	})

	assert.True(t, errors.Is(err, errBiz))
}

func TestRunInTransaction_Success(t *testing.T) {
	db := newTestDB(t)

	err := RunInTransaction(db, func(tx *gorm.DB) error {
		return nil
	})

	assert.NoError(t, err)
}

func TestWithTransaction_Panic_ReturnsError(t *testing.T) {
	db := newTestDB(t)

	err := WithTransaction(context.Background(), db, func(ctx context.Context, tx *gorm.DB) error {
		panic("boom")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "panic in transaction")
}

func TestWithTransaction_Success(t *testing.T) {
	db := newTestDB(t)

	err := WithTransaction(context.Background(), db, func(ctx context.Context, tx *gorm.DB) error {
		return nil
	})

	assert.NoError(t, err)
}
