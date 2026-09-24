package repository

import (
	"context"

	"icmongolang/internal/models"
	"gorm.io/gorm"
)

// Repository provides audit + blockchain anchor persistence.
type Repository interface {
	CreateEntry(ctx context.Context, entry *models.AuditLogEntry) error
	GetEntryByID(ctx context.Context, id string) (*models.AuditLogEntry, error)
	GetLastEntry(ctx context.Context) (*models.AuditLogEntry, error)
	CreateAnchor(ctx context.Context, anchor *models.BlockchainAnchor) error
	GetAnchorByChainID(ctx context.Context, chainID string) (*models.BlockchainAnchor, error)
	CountEntries(ctx context.Context) (int64, error)
}

type pgRepository struct {
	db *gorm.DB
}

// NewPgRepository creates a DB-backed audit log repository.
func NewPgRepository(db *gorm.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) CreateEntry(ctx context.Context, entry *models.AuditLogEntry) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *pgRepository) GetEntryByID(ctx context.Context, id string) (*models.AuditLogEntry, error) {
	var e models.AuditLogEntry
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *pgRepository) GetLastEntry(ctx context.Context) (*models.AuditLogEntry, error) {
	var e models.AuditLogEntry
	if err := r.db.WithContext(ctx).Order("created_at DESC, id DESC").First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *pgRepository) CreateAnchor(ctx context.Context, anchor *models.BlockchainAnchor) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(anchor).Error
}

func (r *pgRepository) GetAnchorByChainID(ctx context.Context, chainID string) (*models.BlockchainAnchor, error) {
	var a models.BlockchainAnchor
	if err := r.db.WithContext(ctx).Where("chain_id = ?", chainID).Order("created_at DESC").First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgRepository) CountEntries(ctx context.Context) (int64, error) {
	var n int64
	if err := r.db.WithContext(ctx).Model(&models.AuditLogEntry{}).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}
