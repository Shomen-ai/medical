// Файл: internal/handler/review.go
// Назначение: HTTP-обработчики отзывов — публичный список с фильтрами, список визитов для формы, создание отзыва пациентом и модерация админом.
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

// ReviewHandler — обработчик запросов, связанных с отзывами.
type ReviewHandler struct{ svc *service.Services }

// NewReviewHandler создаёт ReviewHandler с подключённым сервисным слоем.
func NewReviewHandler(svc *service.Services) *ReviewHandler { return &ReviewHandler{svc: svc} }

// reviewLimit парсит limit/offset из query с дефолтом и максимумом 50.
func reviewLimit(c *gin.Context, def int) (limit, offset int) {
	limit = def
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
		limit = v
	}
	if limit > 50 {
		limit = 50
	}
	if v, err := strconv.Atoi(c.Query("offset")); err == nil && v > 0 {
		offset = v
	}
	return
}

// List возвращает видимые отзывы с опциональными фильтрами по врачу и услуге.
// GET /api/reviews?doctor_id=&service_id=&limit=&offset=
func (h *ReviewHandler) List(c *gin.Context) {
	limit, offset := reviewLimit(c, 10)
	rs, err := h.svc.Reviews.ListPublic(c.Query("doctor_id"), c.Query("service_id"), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rs)
}

// Reviewable возвращает завершённые визиты текущего пациента (для формы отзыва).
// GET /api/cabinet/reviewable
func (h *ReviewHandler) Reviewable(c *gin.Context) {
	as, err := h.svc.Reviews.Reviewable(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, as)
}

// Create создаёт отзыв от лица любого аутентифицированного пациента.
// POST /api/reviews  body: {"rating":5,"text":"...","doctor_id":"","service_id":""}
// doctor_id/service_id — опциональны (для фильтрации на странице отзывов).
func (h *ReviewHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		Rating    int    `json:"rating"     binding:"required"`
		Text      string `json:"text"       binding:"required"`
		DoctorID  string `json:"doctor_id"`
		ServiceID string `json:"service_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rv, err := h.svc.Reviews.Create(userID, req.DoctorID, req.ServiceID, req.Rating, req.Text)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrReviewRating), errors.Is(err, service.ErrReviewText):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, rv)
}

// AdminList возвращает все отзывы (включая скрытые) для модерации.
// GET /api/admin/reviews?limit=&offset=
func (h *ReviewHandler) AdminList(c *gin.Context) {
	limit, offset := reviewLimit(c, 20)
	rs, err := h.svc.Reviews.ListAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rs)
}

// AdminSetHidden скрывает или возвращает отзыв.
// PATCH /api/admin/reviews/:id  body: {"hidden": true}
func (h *ReviewHandler) AdminSetHidden(c *gin.Context) {
	var req struct {
		Hidden bool `json:"hidden"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Reviews.SetHidden(c.Param("id"), req.Hidden); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
