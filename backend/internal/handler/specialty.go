// Файл: internal/handler/specialty.go
// Назначение: публичный HTTP-обработчик справочника медицинских специальностей.
package handler

import (
	"net/http"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

// SpecialtyHandler — обработчик публичных запросов справочника специальностей.
type SpecialtyHandler struct{ svc *service.Services }

// NewSpecialtyHandler создаёт новый SpecialtyHandler с подключённым сервисным слоем.
func NewSpecialtyHandler(svc *service.Services) *SpecialtyHandler { return &SpecialtyHandler{svc: svc} }

// List возвращает список всех медицинских специальностей.
// GET /api/specialties
func (h *SpecialtyHandler) List(c *gin.Context) {
	ss, err := h.svc.Repos.Specialties.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ss)
}
