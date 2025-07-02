package handlers

import (
	"errors"
	"ia-boilerplate/src/repository"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func snakeCase(s string) string {
	var out []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(r))
	}
	return string(out)
}

func (h *Handler) GetMedicine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var m repository.Medicine
	res := h.Repository.DB.Where("id = ? AND is_deleted = ?", id, false).First(&m)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, m)
}

type createMedicineRequest struct {
	EANCode            string  `json:"eanCode" binding:"required"`
	Description        string  `json:"description" binding:"required"`
	Type               string  `json:"type" binding:"required"`
	Laboratory         string  `json:"laboratory"`
	IVA                string  `json:"iva"`
	SatKey             string  `json:"satKey"`
	ActiveIngredient   string  `json:"activeIngredient"`
	TemperatureControl string  `json:"temperatureControl"`
	IsControlled       bool    `json:"isControlled"`
	UnitQuantity       float64 `json:"unitQuantity"`
	UnitType           string  `json:"unitType"`
}

func (h *Handler) CreateMedicine(c *gin.Context) {
	var req createMedicineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing repository.Medicine
	if err := h.Repository.DB.Where("ean_code = ? AND is_deleted = ?", req.EANCode, false).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Could not create medicine: duplicate code"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create medicine"})
		return
	}

	medicineType := repository.MedicineType(req.Type)
	if !medicineType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine type, must be one of: " + strings.Join(repository.ValidMedicineTypes, ", ")})
		return
	}

	temperatureControl := repository.TemperatureControlType(req.TemperatureControl)

	if !temperatureControl.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid temperature control, must be one of: " + strings.Join(repository.ValidTemperatureCtrls, ", ")})
		return
	}

	unitType := repository.UnitType(req.UnitType)
	if !unitType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit type, must be one of: " + strings.Join(repository.ValidUnitTypes, ", ")})
		return
	}

	m := repository.Medicine{
		EANCode:            req.EANCode,
		Description:        req.Description,
		Type:               medicineType,
		Laboratory:         req.Laboratory,
		IVA:                req.IVA,
		SatKey:             req.SatKey,
		TemperatureControl: temperatureControl,
		ActiveIngredient:   req.ActiveIngredient,
		ColdChain:          false,
		IsControlled:       req.IsControlled,
		UnitQuantity:       req.UnitQuantity,
		UnitType:           unitType,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if res := h.Repository.DB.Create(&m); res.Error != nil {
		if strings.Contains(res.Error.Error(), "duplicate key value") {
			c.JSON(http.StatusConflict, gin.H{"error": "Could not create medicine: " + res.Error.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create medicine"})
		}
		return
	}

	c.JSON(http.StatusCreated, m)
}

func (h *Handler) DeleteMedicine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if res := h.Repository.DB.Model(&repository.Medicine{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("is_deleted", true); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete medicine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted successfully"})
}

func (h *Handler) SearchMedicinesPaginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	likeFilters := map[string]string{
		"description":       c.Query("description_like"),
		"laboratory":        c.Query("laboratory_like"),
		"ean_code":          c.Query("ean_code_like"),
		"sat_key":           c.Query("sat_key_like"),
		"active_ingredient": c.Query("active_ingredient_like"),
	}

	matches := map[string][]string{
		"description":       c.QueryArray("description_match"),
		"laboratory":        c.QueryArray("laboratory_match"),
		"ean_code":          c.QueryArray("ean_code_match"),
		"sat_key":           c.QueryArray("sat_key_match"),
		"active_ingredient": c.QueryArray("active_ingredient_match"),
	}

	var (
		medicines []repository.Medicine
		total     int64
	)

	query := h.Repository.DB.
		Model(&repository.Medicine{}).
		Where("is_deleted = ?", false)

	for col, val := range likeFilters {
		if val != "" {
			query = query.Where(col+" ILIKE ?", "%"+val+"%")
		}
	}

	for col, vals := range matches {
		if len(vals) > 0 {
			query = query.Where(col+" IN (?)", vals)
		}
	}

	query.Count(&total)

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&medicines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not perform search"})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	c.JSON(http.StatusOK, gin.H{
		"current_page":  page,
		"medicines":     medicines,
		"page_size":     limit,
		"total_pages":   totalPages,
		"total_records": total,
	})
}

func (h *Handler) SearchMedicineCoincidencesByProperty(c *gin.Context) {
	property := c.Query("property")
	searchText := c.Query("search_text")
	allowed := map[string]bool{
		"description":       true,
		"laboratory":        true,
		"ean_code":          true,
		"sat_key":           true,
		"active_ingredient": true,
	}
	if !allowed[property] || searchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property"})
		return
	}

	var results []string
	if err := h.Repository.DB.
		Model(&repository.Medicine{}).
		Distinct(property).
		Where("is_deleted = ?", false).
		Where(property+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(property, &results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not query coincidences"})
		return
	}

	c.JSON(http.StatusOK, results)
}

type updateMedicineRequest struct {
	EANCode            *string  `json:"eanCode"`
	Description        *string  `json:"description"`
	Type               *string  `json:"type"`
	Laboratory         *string  `json:"laboratory"`
	IVA                *string  `json:"iva"`
	SatKey             *string  `json:"satKey"`
	ActiveIngredient   *string  `json:"activeIngredient"`
	TemperatureControl *string  `json:"temperatureControl"`
	IsControlled       *bool    `json:"isControlled"`
	UnitQuantity       *float64 `json:"unitQuantity"`
	UnitType           *string  `json:"unitType"`
}

func (h *Handler) UpdateMedicine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Verificar que el medicamento existe
	var existingMedicine repository.Medicine
	if err := h.Repository.DB.Where("id = ? AND is_deleted = ?", id, false).First(&existingMedicine).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	var req updateMedicineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preparar los campos a actualizar
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	// Validar y agregar cada campo si está presente en el request
	if req.EANCode != nil {
		// Verificar que el EAN code no esté duplicado (excluyendo el registro actual)
		var duplicateCheck repository.Medicine
		if err := h.Repository.DB.Where("ean_code = ? AND id != ? AND is_deleted = ?", *req.EANCode, id, false).First(&duplicateCheck).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "EAN code already exists"})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking EAN code"})
			return
		}
		updates["ean_code"] = *req.EANCode
	}

	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if req.Type != nil {
		medicineType := repository.MedicineType(*req.Type)
		if !medicineType.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine type, must be one of: " + strings.Join(repository.ValidMedicineTypes, ", ")})
			return
		}
		updates["type"] = medicineType
	}

	if req.Laboratory != nil {
		updates["laboratory"] = *req.Laboratory
	}

	if req.IVA != nil {
		updates["iva"] = *req.IVA
	}

	if req.SatKey != nil {
		updates["sat_key"] = *req.SatKey
	}

	if req.ActiveIngredient != nil {
		updates["active_ingredient"] = *req.ActiveIngredient
	}

	if req.TemperatureControl != nil {
		temperatureControl := repository.TemperatureControlType(*req.TemperatureControl)
		if !temperatureControl.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid temperature control, must be one of: " + strings.Join(repository.ValidTemperatureCtrls, ", ")})
			return
		}
		updates["temperature_control"] = temperatureControl
	}

	if req.IsControlled != nil {
		updates["is_controlled"] = *req.IsControlled
	}

	if req.UnitQuantity != nil {
		updates["unit_quantity"] = *req.UnitQuantity
	}

	if req.UnitType != nil {
		unitType := repository.UnitType(*req.UnitType)
		if !unitType.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit type, must be one of: " + strings.Join(repository.ValidUnitTypes, ", ")})
			return
		}
		updates["unit_type"] = unitType
	}

	// Si no hay campos para actualizar, retornar error
	if len(updates) <= 1 { // Solo updated_at
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Realizar la actualización
	if err := h.Repository.DB.Model(&repository.Medicine{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update medicine"})
		return
	}

	// Obtener el medicamento actualizado
	var updatedMedicine repository.Medicine
	if err := h.Repository.DB.Where("id = ? AND is_deleted = ?", id, false).First(&updatedMedicine).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve updated medicine"})
		return
	}

	c.JSON(http.StatusOK, updatedMedicine)
}
