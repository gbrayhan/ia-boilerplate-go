package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ia-boilerplate/src/repository"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	validMedicineTypes    = map[string]bool{"tableta": true, "capsula": true}
	validTemperatureCtrls = map[string]bool{"seco": true, "temperatura_controlada": true}
	validUnitTypes        = map[string]bool{"tabletas": true, "capsulas": true}
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
	EANCode          string  `json:"eanCode" binding:"required"`
	Description      string  `json:"description" binding:"required"`
	Type             string  `json:"type" binding:"required"`
	Laboratory       string  `json:"laboratory"`
	IVA              string  `json:"iva"`
	SatKey           string  `json:"satKey"`
	ActiveIngredient string  `json:"activeIngredient"`
	TemperatureCtrl  string  `json:"temperatureControl"`
	IsControlled     bool    `json:"isControlled"`
	UnitQuantity     float64 `json:"unitQuantity"`
	UnitType         string  `json:"unitType"`
}

func (h *Handler) CreateMedicine(c *gin.Context) {
	var req createMedicineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validMedicineTypes[req.Type] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine type"})
		return
	}
	if req.TemperatureCtrl != "" && !validTemperatureCtrls[req.TemperatureCtrl] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid temperature control"})
		return
	}
	if req.UnitType != "" && !validUnitTypes[req.UnitType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit type"})
		return
	}

	m := repository.Medicine{
		EANCode:          req.EANCode,
		Description:      req.Description,
		Type:             req.Type,
		Laboratory:       req.Laboratory,
		IVA:              req.IVA,
		SatKey:           req.SatKey,
		TemperatureCtrl:  req.TemperatureCtrl,
		ActiveIngredient: req.ActiveIngredient,
		ColdChain:        false,
		IsControlled:     req.IsControlled,
		UnitQuantity:     req.UnitQuantity,
		UnitType:         req.UnitType,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
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

func (h *Handler) UpdateMedicine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(payload) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
		return
	}

	var m repository.Medicine
	if res := h.Repository.DB.Where("id = ? AND is_deleted = ?", id, false).First(&m); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	updates := map[string]interface{}{}
	for k, v := range payload {
		switch k {
		case "type":
			s, ok := v.(string)
			if !ok || !validMedicineTypes[s] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine type"})
				return
			}
			updates["type"] = s
		case "temperatureControl":
			s, ok := v.(string)
			if !ok || !validTemperatureCtrls[s] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid temperature control"})
				return
			}
			updates["temperature_ctrl"] = s
		case "unitType":
			s, ok := v.(string)
			if !ok || !validUnitTypes[s] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit type"})
				return
			}
			updates["unit_type"] = s
		case "eanCode", "description", "laboratory", "iva", "satKey", "activeIngredient", "isControlled", "unitQuantity", "coldChain", "isDeleted":
			updates[snakeCase(k)] = v
		}
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
		return
	}

	updates["updated_at"] = time.Now()

	if err := h.Repository.DB.Model(&m).Updates(updates).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusConflict, gin.H{"error": "Could not update medicine: " + err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update medicine"})
		}
		return
	}

	h.Repository.DB.First(&m, id)
	c.JSON(http.StatusOK, m)
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
