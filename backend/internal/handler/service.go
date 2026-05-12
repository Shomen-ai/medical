// Файл: internal/handler/service.go
// Назначение: публичные HTTP-обработчики каталога медицинских услуг — список услуг (с фильтром по специальности) и карточка одной услуги.
package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

// ServiceHandler — обработчик публичных запросов каталога услуг.
type ServiceHandler struct{ svc *service.Services }

// NewServiceHandler создаёт новый ServiceHandler с подключённым сервисным слоем.
func NewServiceHandler(svc *service.Services) *ServiceHandler { return &ServiceHandler{svc: svc} }

// List возвращает список услуг; при указании specialty_id фильтрует по специальности.
// GET /api/services?specialty_id=uuid
func (h *ServiceHandler) List(c *gin.Context) {
	ss, err := h.svc.Repos.Services.List(c.Query("specialty_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ss)
}

// Get возвращает карточку услуги по идентификатору.
// GET /api/services/:id
func (h *ServiceHandler) Get(c *gin.Context) {
	s, err := h.svc.Repos.Services.FindByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}
