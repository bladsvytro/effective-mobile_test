package handlers

import (
	"time"

	"github.com/google/uuid"
)

// CreateSubscriptionRequest запрос на создание подписки
type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int       `json:"price" binding:"required,min=1" example:"400"`
	UserID      uuid.UUID `json:"user_id" binding:"required" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" binding:"required" example:"07-2025"` // формат MM-YYYY
	EndDate     *string   `json:"end_date,omitempty" example:"12-2025"`
}

// UpdateSubscriptionRequest запрос на обновление подписки
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty" example:"Yandex Music"`
	Price       *int    `json:"price,omitempty" binding:"omitempty,min=1" example:"500"`
	StartDate   *string `json:"start_date,omitempty" example:"08-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2025"`
}

// SubscriptionResponse ответ с данными подписки
type SubscriptionResponse struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TotalPriceRequest запрос для подсчета суммарной стоимости
type TotalPriceRequest struct {
	StartDate   string     `form:"start_date" binding:"required" example:"01-2025"`
	EndDate     string     `form:"end_date" binding:"required" example:"12-2025"`
	UserID      *uuid.UUID `form:"user_id,omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName *string    `form:"service_name,omitempty" example:"Yandex Plus"`
}

// TotalPriceResponse ответ с суммарной стоимостью
type TotalPriceResponse struct {
	Total int `json:"total"`
}

// ErrorResponse стандартный ответ об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}
