package handler

import (
	"net/http"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct{ svc *service.Services }

func NewServiceHandler(svc *service.Services) *ServiceHandler { return &ServiceHandler{svc} }

// GET /api/services?specialty_id=uuid
func (h *ServiceHandler) List(c *gin.Context) {
	ss, err := h.svc.Repos.Services.List(c.Query("specialty_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ss)
}

// GET /api/services/:id
func (h *ServiceHandler) Get(c *gin.Context) {
	s, err := h.svc.Repos.Services.FindByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}
