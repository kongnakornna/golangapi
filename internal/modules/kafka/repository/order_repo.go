package repository

import (
	"icmongolang/internal/modules/kafka/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	Update(order *models.Order) error
	FindByID(id string) (*models.Order, error)
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepo) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepo) FindByID(id string) (*models.Order, error) {
	var order models.Order
	if err := r.db.Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
