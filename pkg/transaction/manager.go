package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Manager คืออินเทอร์เฟซผู้จัดการธุรกรรม
type Manager interface {
	// Execute ดำเนินการธุรกรรม
	Execute(ctx context.Context, fn TxFunc) error
	// ExecuteWithOptions ดำเนินการธุรกรรมพร้อมตัวเลือก
	ExecuteWithOptions(ctx context.Context, opts *sql.TxOptions, fn TxFunc) error
}

// TxFunc คือชนิดฟังก์ชันธุรกรรม
type TxFunc func(ctx context.Context, tx *gorm.DB) error

// GormTransactionManager คือผู้จัดการธุรกรรมสำหรับ GORM
type GormTransactionManager struct {
	db *gorm.DB
}

// NewGormTransactionManager สร้างผู้จัดการธุรกรรม GORM ใหม่
func NewGormTransactionManager(db *gorm.DB) Manager {
	return &GormTransactionManager{db: db}
}

// Execute ดำเนินการธุรกรรม
func (m *GormTransactionManager) Execute(ctx context.Context, fn TxFunc) error {
	return m.ExecuteWithOptions(ctx, nil, fn)
}

// ExecuteWithOptions ดำเนินการธุรกรรมพร้อมตัวเลือก
func (m *GormTransactionManager) ExecuteWithOptions(ctx context.Context, opts *sql.TxOptions, fn TxFunc) error {
	// เริ่มต้นธุรกรรม
	tx := m.db.WithContext(ctx)
	if opts != nil {
		tx = tx.Begin(opts)
	} else {
		tx = tx.Begin()
	}

	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// ใช้ defer เพื่อให้แน่ใจว่าธุรกรรมจะสิ้นสุดเสมอ
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // โยน panic ออกไปอีกครั้ง
		}
	}()

	// ดำเนินการฟังก์ชันธุรกรรม
	if err := fn(ctx, tx); err != nil {
		// ย้อนกลับธุรกรรม
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// ยืนยันธุรกรรม
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RunInTransaction ดำเนินการฟังก์ชันภายในธุรกรรม (แบบง่าย)
func RunInTransaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Transactional คือตัวตกแต่งธุรกรรม
type Transactional struct {
	db *gorm.DB
}

// NewTransactional สร้างตัวตกแต่งธุรกรรมใหม่
func NewTransactional(db *gorm.DB) *Transactional {
	return &Transactional{db: db}
}

// Wrap ห่อหุ้มฟังก์ชันให้ดำเนินการภายในธุรกรรม
func (t *Transactional) Wrap(fn func(*gorm.DB) error) error {
	return RunInTransaction(t.db, fn)
}

// NestedTransaction รองรับธุรกรรมแบบซ้อน
type NestedTransaction struct {
	db           *gorm.DB
	savepoints   []string
	currentLevel int
}

// NewNestedTransaction สร้างธุรกรรมแบบซ้อนใหม่
func NewNestedTransaction(db *gorm.DB) *NestedTransaction {
	return &NestedTransaction{
		db:           db,
		savepoints:   make([]string, 0),
		currentLevel: 0,
	}
}

// Begin เริ่มต้นธุรกรรมใหม่หรือบันทึกจุด
func (nt *NestedTransaction) Begin() error {
	if nt.currentLevel == 0 {
		// เริ่มต้นธุรกรรมใหม่
		tx := nt.db.Begin()
		if tx.Error != nil {
			return tx.Error
		}
		nt.db = tx
	} else {
		// สร้าง savepoint
		savepoint := fmt.Sprintf("sp_%d", nt.currentLevel)
		if err := nt.db.Exec("SAVEPOINT " + savepoint).Error; err != nil {
			return err
		}
		nt.savepoints = append(nt.savepoints, savepoint)
	}
	nt.currentLevel++
	return nil
}

// Commit ยืนยันธุรกรรมหรือปล่อย savepoint
func (nt *NestedTransaction) Commit() error {
	if nt.currentLevel == 0 {
		return errors.New("no transaction to commit")
	}

	if nt.currentLevel == 1 {
		// ยืนยันธุรกรรม
		if err := nt.db.Commit().Error; err != nil {
			return err
		}
	} else {
		// ปล่อย savepoint
		savepoint := nt.savepoints[len(nt.savepoints)-1]
		if err := nt.db.Exec("RELEASE SAVEPOINT " + savepoint).Error; err != nil {
			return err
		}
		nt.savepoints = nt.savepoints[:len(nt.savepoints)-1]
	}
	nt.currentLevel--
	return nil
}

// Rollback ย้อนกลับธุรกรรมหรือย้อนกลับไปยัง savepoint
func (nt *NestedTransaction) Rollback() error {
	if nt.currentLevel == 0 {
		return errors.New("no transaction to rollback")
	}

	if nt.currentLevel == 1 {
		// ย้อนกลับทั้งธุรกรรม
		if err := nt.db.Rollback().Error; err != nil {
			return err
		}
	} else {
		// ย้อนกลับไปยัง savepoint
		savepoint := nt.savepoints[len(nt.savepoints)-1]
		if err := nt.db.Exec("ROLLBACK TO SAVEPOINT " + savepoint).Error; err != nil {
			return err
		}
		nt.savepoints = nt.savepoints[:len(nt.savepoints)-1]
	}
	nt.currentLevel--
	return nil
}

// TransactionContext คือบริบทของธุรกรรม
type TransactionContext struct {
	ctx context.Context
	tx  *gorm.DB
}

// NewTransactionContext สร้างบริบทธุรกรรมใหม่
func NewTransactionContext(ctx context.Context, db *gorm.DB) (*TransactionContext, error) {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &TransactionContext{
		ctx: ctx,
		tx:  tx,
	}, nil
}

// DB รับการเชื่อมต่อฐานข้อมูลของธุรกรรม
func (tc *TransactionContext) DB() *gorm.DB {
	return tc.tx
}

// Context รับบริบท
func (tc *TransactionContext) Context() context.Context {
	return tc.ctx
}

// Commit ยืนยันธุรกรรม
func (tc *TransactionContext) Commit() error {
	return tc.tx.Commit().Error
}

// Rollback ย้อนกลับธุรกรรม
func (tc *TransactionContext) Rollback() error {
	return tc.tx.Rollback().Error
}

// Complete ตัดสินใจยืนยันหรือย้อนกลับตาม error
func (tc *TransactionContext) Complete(err error) error {
	if err != nil {
		if rbErr := tc.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
		}
		return err
	}
	return tc.Commit()
}

// WithTransaction ดำเนินการฟังก์ชันภายในธุรกรรม (พร้อมบริบท)
func WithTransaction(ctx context.Context, db *gorm.DB, fn func(context.Context, *gorm.DB) error) error {
	tc, err := NewTransactionContext(ctx, db)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tc.Rollback()
			panic(r)
		}
	}()

	err = fn(tc.Context(), tc.DB())
	return tc.Complete(err)
}
