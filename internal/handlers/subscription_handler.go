package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"subscription-service/internal/models"
	"subscription-service/internal/repository"
)

type SubscriptionHandler struct {
	repo   repository.SubscriptionRepository
	logger *zap.Logger
}

func NewSubscriptionHandler(repo repository.SubscriptionRepository, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo, logger: logger}
}

// Create обрабатывает создание новой подписки
// @Summary Создать подписку
// @Description Создает новую запись о подписке пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Невалидный запрос", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Невалидный запрос: " + err.Error()})
		return
	}

	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат start_date, ожидается MM-YYYY"})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		ed, err := parseMonthYear(*req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат end_date, ожидается MM-YYYY"})
			return
		}
		endDate = &ed
	}

	subscription := &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.repo.Create(subscription); err != nil {
		h.logger.Error("Ошибка создания подписки", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Внутренняя ошибка сервера"})
		return
	}

	h.logger.Info("Подписка создана", zap.String("id", subscription.ID.String()))
	c.JSON(http.StatusCreated, subscriptionToResponse(subscription))
}

// GetByID возвращает подписку по ID
// @Summary Получить подписку по ID
// @Description Возвращает запись о подписке по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат ID"})
		return
	}

	subscription, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.Warn("Подписка не найдена", zap.String("id", id.String()), zap.Error(err))
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Подписка не найдена"})
		return
	}

	c.JSON(http.StatusOK, subscriptionToResponse(subscription))
}

// Update обновляет подписку
// @Summary Обновить подписку
// @Description Обновляет существующую запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Param input body UpdateSubscriptionRequest true "Обновляемые поля"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат ID"})
		return
	}

	subscription, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Подписка не найдена"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Невалидный запрос: " + err.Error()})
		return
	}

	if req.ServiceName != nil {
		subscription.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		if *req.Price < 1 {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Цена должна быть положительным числом"})
			return
		}
		subscription.Price = *req.Price
	}
	if req.StartDate != nil {
		startDate, err := parseMonthYear(*req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат start_date"})
			return
		}
		subscription.StartDate = startDate
	}
	if req.EndDate != nil {
		endDate, err := parseMonthYear(*req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат end_date"})
			return
		}
		subscription.EndDate = &endDate
	}

	if err := h.repo.Update(subscription); err != nil {
		h.logger.Error("Ошибка обновления подписки", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Внутренняя ошибка сервера"})
		return
	}

	h.logger.Info("Подписка обновлена", zap.String("id", subscription.ID.String()))
	c.JSON(http.StatusOK, subscriptionToResponse(subscription))
}

// Delete удаляет подписку
// @Summary Удалить подписку
// @Description Удаляет запись о подписке (мягкое удаление)
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат ID"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		h.logger.Warn("Подписка не найдена", zap.String("id", id.String()), zap.Error(err))
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Подписка не найдена"})
		return
	}

	h.logger.Info("Подписка удалена", zap.String("id", id.String()))
	c.Status(http.StatusNoContent)
}

// List возвращает список подписок с пагинацией
// @Summary Список подписок
// @Description Возвращает список подписок с поддержкой пагинации
// @Tags subscriptions
// @Produce json
// @Param offset query int false "Смещение (по умолчанию 0)" default(0)
// @Param limit query int false "Лимит (по умолчанию 10, максимум 100)" default(10)
// @Success 200 {array} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	h.logger.Debug("Запрос списка подписок", zap.Int("offset", offset), zap.Int("limit", limit))

	subscriptions, err := h.repo.List(offset, limit)
	if err != nil {
		h.logger.Error("Ошибка получения списка", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Внутренняя ошибка сервера"})
		return
	}

	h.logger.Debug("Получено подписок", zap.Int("count", len(subscriptions)))

	response := make([]SubscriptionResponse, len(subscriptions))
	for i, s := range subscriptions {
		response[i] = subscriptionToResponse(&s)
	}

	c.JSON(http.StatusOK, response)
}

// GetTotalPrice подсчитывает суммарную стоимость подписок за период
// @Summary Суммарная стоимость подписок
// @Description Возвращает суммарную стоимость всех подписок за выбранный период с фильтрацией по пользователю и названию сервиса. Учитывается количество месяцев активности подписки в пределах запрошенного периода (пересечение периодов). Стоимость каждой подписки умножается на количество месяцев пересечения.
// @Tags analytics
// @Produce json
// @Param start_date query string true "Начало периода (MM-YYYY)"
// @Param end_date query string true "Конец периода (MM-YYYY)"
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Success 200 {object} TotalPriceResponse
// @Failure 400 {object} ErrorResponse
// @Router /analytics/total [get]
func (h *SubscriptionHandler) GetTotalPrice(c *gin.Context) {
	var req TotalPriceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Невалидные параметры запроса: " + err.Error()})
		return
	}

	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат start_date"})
		return
	}
	endDate, err := parseMonthYear(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Неверный формат end_date"})
		return
	}

	h.logger.Debug("Запрос аналитики",
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate),
		zap.Any("user_id", req.UserID),
		zap.Any("service_name", req.ServiceName),
	)

	total, err := h.repo.GetTotalPrice(startDate, endDate, req.UserID, req.ServiceName)
	if err != nil {
		h.logger.Error("Ошибка подсчета стоимости", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Внутренняя ошибка сервера"})
		return
	}

	h.logger.Info("Рассчитана суммарная стоимость", zap.Int("total", total))
	c.JSON(http.StatusOK, TotalPriceResponse{Total: total})
}

// parseMonthYear преобразует строку формата "MM-YYYY" в time.Time (первый день месяца)
func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

// subscriptionToResponse преобразует модель в ответ API
func subscriptionToResponse(s *models.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
