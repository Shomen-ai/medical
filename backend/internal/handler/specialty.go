package handler

import (
	"net/http"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

type SpecialtyHandler struct{ svc *service.Services }

func NewSpecialtyHandler(svc *service.Services) *SpecialtyHandler { return &SpecialtyHandler{svc} }

// GET /api/specialties
func (h *SpecialtyHandler) List(c *gin.Context) {
	ss, err := h.svc.Repos.Specialties.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ss)
}
