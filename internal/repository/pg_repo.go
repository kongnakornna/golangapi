// Package repository ให้ base generic repository
package repository

import (
	"context"

	"gorm.io/gorm"
)

type PGRepo[T any] struct {
	db *gorm.DB
}

func CreatePGRepo[T any](db *gorm.DB) PGRepo[T] {
	return PGRepo[T]{db: db}
}

func (r *PGRepo[T]) DB() *gorm.DB { return r.db }

func (r *PGRepo[T]) Create(ctx context.Context, e *T) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *PGRepo[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	var out T
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PGRepo[T]) Update(ctx context.Context, e *T) error {
	return r.db.WithContext(ctx).Save(e).Error
}

func (r *PGRepo[T]) Delete(ctx context.Context, id interface{}) error {
	var zero T
	return r.db.WithContext(ctx).Delete(&zero, "id = ?", id).Error
}

func (r *PGRepo[T]) List(ctx context.Context) ([]T, error) {
	var out []T
	if err := r.db.WithContext(ctx).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}
