// Файл: internal/handler/doctor.go
// Назначение: публичные HTTP-обработчики каталога врачей — список всех врачей (с фильтром по специальности) и карточка конкретного врача.
package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"beautymed/internal/service"

	"github.com/gin-gonic/gin"
)

// DoctorHandler — обработчик публичных запросов каталога врачей.
type DoctorHandler struct{ svc *service.Services }

// NewDoctorHandler создаёт новый DoctorHandler с подключённым сервисным слоем.
func NewDoctorHandler(svc *service.Services) *DoctorHandler { return &DoctorHandler{svc: svc} }

// List возвращает всех активных врачей; при указании specialty_id фильтрует по специальности.
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

// Get возвращает карточку врача по идентификатору.
// GET /api/doctors/:id
func (h *DoctorHandler) Get(c *gin.Context) {
	d, err := h.svc.Repos.Doctors.FindByID(c.Param("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}
