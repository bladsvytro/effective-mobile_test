package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"subscription-service/internal/models"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	GetByID(id uuid.UUID) (*models.Subscription, error)
	Update(subscription *models.Subscription) error
	Delete(id uuid.UUID) error
	List(offset, limit int) ([]models.Subscription, error)
	GetByUserID(userID uuid.UUID) ([]models.Subscription, error)
	GetByServiceName(serviceName string) ([]models.Subscription, error)
	GetTotalPrice(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (int, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepo) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.db.First(&subscription, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepo) Update(subscription *models.Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *subscriptionRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Subscription{}, "id = ?", id).Error
}

func (r *subscriptionRepo) List(offset, limit int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Offset(offset).Limit(limit).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepo) GetByUserID(userID uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Where("user_id = ?", userID).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepo) GetByServiceName(serviceName string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Where("service_name = ?", serviceName).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepo) GetTotalPrice(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	var total int
	query := r.db.Model(&models.Subscription{}).Where("start_date >= ? AND start_date <= ?", startDate, endDate)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if serviceName != nil {
		query = query.Where("service_name = ?", *serviceName)
	}

	err := query.Select("COALESCE(SUM(price), 0)").Scan(&total).Error
	return total, err
}
