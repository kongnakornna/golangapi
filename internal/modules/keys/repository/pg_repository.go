package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	pkgapikey "icmongolang/pkg/apikey"
)

var errNotFound = errors.New("API key not found")

// Postgres is the GORM-backed implementation of APIKeyRepository.
type Postgres struct {
	db *gorm.DB
}

// NewPostgresRepository builds a Postgres repository.
func NewPostgresRepository(db *gorm.DB) APIKeyRepository {
	return &Postgres{db: db}
}

// ---------------------------------------------------------------- create

func (r *Postgres) Create(ctx context.Context, k *pkgapikey.Key) (*pkgapikey.Key, error) {
	if err := r.db.WithContext(ctx).Create(k).Error; err != nil {
		return nil, err
	}
	return k, nil
}

// ---------------------------------------------------------------- reads

func (r *Postgres) FindByID(ctx context.Context, id uuid.UUID) (*pkgapikey.Key, error) {
	var k pkgapikey.Key
	if err := r.db.WithContext(ctx).First(&k, "id = ?", id).Error; err != nil {
		return nil, mapErr(err)
	}
	return &k, nil
}

func (r *Postgres) FindByHashedKey(ctx context.Context, hash string) (*pkgapikey.Key, error) {
	var k pkgapikey.Key
	if err := r.db.WithContext(ctx).First(&k, "hashed_key = ?", hash).Error; err != nil {
		return nil, mapErr(err)
	}
	return &k, nil
}

// ---------------------------------------------------------------- list

func (r *Postgres) List(ctx context.Context, limit, offset int) ([]*pkgapikey.Key, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&pkgapikey.Key{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []*pkgapikey.Key
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Postgres) ListByOwner(ctx context.Context, owner uuid.UUID, limit, offset int) ([]*pkgapikey.Key, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&pkgapikey.Key{}).
		Where("created_by = ?", owner).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []*pkgapikey.Key
	if err := r.db.WithContext(ctx).
		Where("created_by = ?", owner).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ---------------------------------------------------------------- mutate

func (r *Postgres) Rotate(ctx context.Context, k *pkgapikey.Key) (*pkgapikey.Key, error) {
	// Save() on a struct with a populated PK issues an UPDATE and lets GORM
	// auto-bump UpdatedAt (the field has no `autoUpdateTime:false` override).
	if err := r.db.WithContext(ctx).Save(k).Error; err != nil {
		return nil, err
	}
	return k, nil
}

func (r *Postgres) Deactivate(ctx context.Context, id, updatedBy uuid.UUID) error {
	res := r.db.WithContext(ctx).Model(&pkgapikey.Key{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_active":  false,
			"updated_by": updatedBy,
			"updated_at": gorm.Expr("NOW()"),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return pkgapikey.ErrNotFound
	}
	return nil
}

func (r *Postgres) DeactivateOwned(ctx context.Context, id, owner, updatedBy uuid.UUID) error {
	res := r.db.WithContext(ctx).Model(&pkgapikey.Key{}).
		Where("id = ? AND created_by = ?", id, owner).
		Updates(map[string]any{
			"is_active":  false,
			"updated_by": updatedBy,
			"updated_at": gorm.Expr("NOW()"),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// Either the key doesn't exist or the caller isn't the owner.
		// We deliberately return NotFound here so non-admins can't
		// enumerate other users' key IDs.
		return errNotFound
	}
	return nil
}

// ---------------------------------------------------------------- helper

func mapErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errNotFound
	}
	return err
}