package repository

import (
	"context"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/customer"
	"icmongolang/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Customer Repository ──────────────────────────────────────────────────

type CustomerPgRepo struct {
	repository.PgRepo[models.Customer]
	DB *gorm.DB
}

func CreateCustomerPgRepository(db *gorm.DB) customer.CustomerPgRepository {
	return &CustomerPgRepo{
		PgRepo: repository.CreatePgRepo[models.Customer](db),
		DB:     db,
	}
}

func (r *CustomerPgRepo) GetMultiByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Customer, error) {
	var objs []*models.Customer
	r.DB.WithContext(ctx).Where("user_id = ?", userID.String()).Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *CustomerPgRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Customer{}).Count(&count).Error
	return count, err
}

func (r *CustomerPgRepo) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Customer{}).Where("user_id = ?", userID.String()).Count(&count).Error
	return count, err
}

// ─── Car Repository ───────────────────────────────────────────────────────

type CarPgRepo struct {
	repository.PgRepo[models.Car]
	DB *gorm.DB
}

func CreateCarPgRepository(db *gorm.DB) customer.CarPgRepository {
	return &CarPgRepo{
		PgRepo: repository.CreatePgRepo[models.Car](db),
		DB:     db,
	}
}

func (r *CarPgRepo) GetMultiByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*models.Car, error) {
	var objs []*models.Car
	r.DB.WithContext(ctx).Where("customer_id = ?", customerID.String()).Limit(limit).Offset(offset).Find(&objs)
	return objs, nil
}

func (r *CarPgRepo) CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Car{}).Where("customer_id = ?", customerID.String()).Count(&count).Error
	return count, err
}
