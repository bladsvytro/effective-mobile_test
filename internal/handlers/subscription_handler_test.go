package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"subscription-service/internal/models"
)

// mockSubscriptionRepository реализует SubscriptionRepository для тестов
type mockSubscriptionRepository struct {
	subscriptions map[uuid.UUID]models.Subscription
	totalPrice    int
	totalPriceErr error
}

func (m *mockSubscriptionRepository) Create(subscription *models.Subscription) error {
	if m.subscriptions == nil {
		m.subscriptions = make(map[uuid.UUID]models.Subscription)
	}
	m.subscriptions[subscription.ID] = *subscription
	return nil
}

func (m *mockSubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	sub, ok := m.subscriptions[id]
	if !ok {
		return nil, nil
	}
	return &sub, nil
}

func (m *mockSubscriptionRepository) Update(subscription *models.Subscription) error {
	m.subscriptions[subscription.ID] = *subscription
	return nil
}

func (m *mockSubscriptionRepository) Delete(id uuid.UUID) error {
	delete(m.subscriptions, id)
	return nil
}

func (m *mockSubscriptionRepository) List(offset, limit int) ([]models.Subscription, error) {
	// Простая реализация для тестов
	subs := make([]models.Subscription, 0, len(m.subscriptions))
	for _, sub := range m.subscriptions {
		subs = append(subs, sub)
	}
	if offset >= len(subs) {
		return []models.Subscription{}, nil
	}
	end := offset + limit
	if end > len(subs) {
		end = len(subs)
	}
	return subs[offset:end], nil
}

func (m *mockSubscriptionRepository) GetByUserID(userID uuid.UUID) ([]models.Subscription, error) {
	var result []models.Subscription
	for _, sub := range m.subscriptions {
		if sub.UserID == userID {
			result = append(result, sub)
		}
	}
	return result, nil
}

func (m *mockSubscriptionRepository) GetByServiceName(serviceName string) ([]models.Subscription, error) {
	var result []models.Subscription
	for _, sub := range m.subscriptions {
		if sub.ServiceName == serviceName {
			result = append(result, sub)
		}
	}
	return result, nil
}

func (m *mockSubscriptionRepository) GetTotalPrice(startDate, endDate time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	return m.totalPrice, m.totalPriceErr
}

func TestSubscriptionHandler_GetTotalPrice(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := &mockSubscriptionRepository{
		totalPrice: 1500,
	}
	handler := NewSubscriptionHandler(repo, logger)

	// Настраиваем Gin в тестовом режиме
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/analytics/total", handler.GetTotalPrice)

	// Создаем запрос с параметрами
	req, _ := http.NewRequest("GET", "/analytics/total?start_date=01-2025&end_date=12-2025", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response TotalPriceResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Total != 1500 {
		t.Errorf("Expected total 1500, got %d", response.Total)
	}
}

func TestSubscriptionHandler_Create(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := &mockSubscriptionRepository{
		subscriptions: make(map[uuid.UUID]models.Subscription),
	}
	handler := NewSubscriptionHandler(repo, logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/subscriptions", handler.Create)

	userID := uuid.New()
	reqBody := CreateSubscriptionRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      userID,
		StartDate:   "01-2025",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/subscriptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp SubscriptionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.ServiceName != "Yandex Plus" {
		t.Errorf("Expected service name Yandex Plus, got %s", resp.ServiceName)
	}
	if resp.Price != 400 {
		t.Errorf("Expected price 400, got %d", resp.Price)
	}
	if resp.UserID != userID {
		t.Errorf("UserID mismatch")
	}
}

func TestSubscriptionHandler_GetByID(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	repo := &mockSubscriptionRepository{
		subscriptions: make(map[uuid.UUID]models.Subscription),
	}
	handler := NewSubscriptionHandler(repo, logger)

	// Создаем подписку через репозиторий
	sub := models.Subscription{
		ID:          uuid.New(),
		ServiceName: "Test Service",
		Price:       300,
		UserID:      uuid.New(),
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	repo.Create(&sub)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/subscriptions/:id", handler.GetByID)

	req, _ := http.NewRequest("GET", "/subscriptions/"+sub.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp SubscriptionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.ID != sub.ID {
		t.Errorf("ID mismatch")
	}
}