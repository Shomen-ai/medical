package handler

import (
	"net/http"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct{ svc *service.Services }

func NewDoctorHandler(svc *service.Services) *DoctorHandler { return &DoctorHandler{svc} }

// GET /api/doctors?specialty_id=uuid
func (h *DoctorHandler) List(c *gin.Context) {
	specialtyID := c.Query("specialty_id")
	if specialtyID != "" {
		ds, err := h.svc.Repos.Doctors.ListBySpecialty(specialtyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ds)
		return
	}
	ds, err := h.svc.Repos.Doctors.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ds)
}

// GET /api/doctors/:id
func (h *DoctorHandler) Get(c *gin.Context) {
	d, err := h.svc.Repos.Doctors.FindByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}
	c.JSON(http.StatusOK, d)
}
