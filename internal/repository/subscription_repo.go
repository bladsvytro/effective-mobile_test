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

// monthsBetween возвращает количество полных месяцев между двумя датами.
// Ожидается, что даты нормализованы к первому дню месяца.
func monthsBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}
	years := b.Year() - a.Year()
	months := int(b.Month()) - int(a.Month())
	totalMonths := years*12 + months
	if totalMonths < 0 {
		return 0
	}
	return totalMonths
}

// monthsIntersection возвращает количество месяцев пересечения двух интервалов.
// Если end2 nil, считается бессрочным интервалом (пересечение до end1).
func monthsIntersection(start1, end1 time.Time, start2 time.Time, end2 *time.Time) int {
	// Определяем границы пересечения
	intervalStart := start1
	if start2.After(intervalStart) {
		intervalStart = start2
	}

	var intervalEnd time.Time
	if end2 == nil {
		// Бессрочная подписка: пересечение ограничено end1
		intervalEnd = end1
	} else {
		intervalEnd = end1
		if end2.Before(intervalEnd) {
			intervalEnd = *end2
		}
	}

	if intervalStart.After(intervalEnd) {
		return 0
	}

	return monthsBetween(intervalStart, intervalEnd) + 1 // +1 потому что включаем оба месяца
}

func (r *subscriptionRepo) GetTotalPrice(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	// Получаем все подписки, которые пересекаются с периодом
	query := r.db.Model(&models.Subscription{}).
		Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", endDate, startDate)

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if serviceName != nil {
		query = query.Where("service_name = ?", *serviceName)
	}

	var subscriptions []models.Subscription
	err := query.Find(&subscriptions).Error
	if err != nil {
		return 0, err
	}

	total := 0
	for _, sub := range subscriptions {
		months := monthsIntersection(startDate, endDate, sub.StartDate, sub.EndDate)
		total += sub.Price * months
	}

	return total, nil
}
