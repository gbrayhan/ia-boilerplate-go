package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ia-boilerplate/src/db"
	"net/http"
	"strconv"
	"time"
)

func GetMedicine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var m db.Medicine
	res := db.DB.Where("id = ? AND is_deleted = ?", id, false).First(&m)
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
	ColdChain        bool    `json:"coldChain"`
	IsControlled     bool    `json:"isControlled"`
	UnitQuantity     float64 `json:"unitQuantity"`
	UnitType         string  `json:"unitType"`
}

func CreateMedicine(c *gin.Context) {
	var req createMedicineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := db.Medicine{
		EANCode:          req.EANCode,
		Description:      req.Description,
		Type:             req.Type,
		Laboratory:       req.Laboratory,
		IVA:              req.IVA,
		SatKey:           req.SatKey,
		ActiveIngredient: req.ActiveIngredient,
		ColdChain:        req.ColdChain,
		IsControlled:     req.IsControlled,
		UnitQuantity:     req.UnitQuantity,
		UnitType:         req.UnitType,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if res := db.DB.Create(&m); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create medicine"})
		return
	}

	c.JSON(http.StatusCreated, m)
}

type updateMedicineRequest struct {
	EANCode          string  `json:"eanCode" binding:"required"`
	Description      string  `json:"description" binding:"required"`
	Type             string  `json:"type" binding:"required"`
	Laboratory       string  `json:"laboratory"`
	IVA              string  `json:"iva"`
	SatKey           string  `json:"satKey"`
	ActiveIngredient string  `json:"activeIngredient"`
	ColdChain        bool    `json:"coldChain"`
	IsControlled     bool    `json:"isControlled"`
	IsDeleted        bool    `json:"isDeleted"`
	UnitQuantity     float64 `json:"unitQuantity"`
	UnitType         string  `json:"unitType"`
}

func UpdateMedicine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req updateMedicineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var m db.Medicine
	if res := db.DB.Where("id = ? AND is_deleted = ?", id, false).First(&m); res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Asignar campos
	m.EANCode = req.EANCode
	m.Description = req.Description
	m.Type = req.Type
	m.Laboratory = req.Laboratory
	m.IVA = req.IVA
	m.SatKey = req.SatKey
	m.ActiveIngredient = req.ActiveIngredient
	m.ColdChain = req.ColdChain
	m.IsControlled = req.IsControlled
	m.IsDeleted = req.IsDeleted
	m.UnitQuantity = req.UnitQuantity
	m.UnitType = req.UnitType
	m.UpdatedAt = time.Now()

	if res := db.DB.Save(&m); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update medicine"})
		return
	}

	c.JSON(http.StatusOK, m)
}

func DeleteMedicine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if res := db.DB.Model(&db.Medicine{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("is_deleted", true); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete medicine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted successfully"})
}

func SearchMedicinesPaginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	descriptionLike := c.Query("description_like")
	laboratoryLike := c.Query("laboratory_like")
	eanCodeLike := c.Query("ean_code_like")
	satKeyLike := c.Query("sat_key_like")
	activeIngredientLike := c.Query("active_ingredient_like")

	descriptionMatches := c.QueryArray("description_match")
	laboratoryMatches := c.QueryArray("laboratory_match")
	eanCodeMatches := c.QueryArray("ean_code_match")
	satKeyMatches := c.QueryArray("sat_key_match")
	activeIngredientMatches := c.QueryArray("active_ingredient_match")

	isOnCostCatalog, err := strconv.ParseBool(c.DefaultQuery("is_on_cost_catalog", "false"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid is_on_cost_catalog parameter"})
		return
	}
	clientID, err := strconv.Atoi(c.DefaultQuery("client_id", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client_id parameter"})
		return
	}
	if isOnCostCatalog && clientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id is required when is_on_cost_catalog is true"})
		return
	}

	var medicines []db.Medicine
	var total int64
	today := time.Now()

	query := db.DB.Model(&db.Medicine{}).Where("is_deleted = ?", false)

	if descriptionLike != "" {
		query = query.Where("description ILIKE ?", "%"+descriptionLike+"%")
	}
	if laboratoryLike != "" {
		query = query.Where("laboratory ILIKE ?", "%"+laboratoryLike+"%")
	}
	if eanCodeLike != "" {
		query = query.Where("ean_code ILIKE ?", "%"+eanCodeLike+"%")
	}
	if satKeyLike != "" {
		query = query.Where("sat_key ILIKE ?", "%"+satKeyLike+"%")
	}
	if activeIngredientLike != "" {
		query = query.Where("active_ingredient ILIKE ?", "%"+activeIngredientLike+"%")
	}
	if len(descriptionMatches) > 0 {
		query = query.Where("description IN ?", descriptionMatches)
	}
	if len(laboratoryMatches) > 0 {
		query = query.Where("laboratory IN ?", laboratoryMatches)
	}
	if len(eanCodeMatches) > 0 {
		query = query.Where("ean_code IN ?", eanCodeMatches)
	}
	if len(satKeyMatches) > 0 {
		query = query.Where("sat_key IN ?", satKeyMatches)
	}
	if len(activeIngredientMatches) > 0 {
		query = query.Where("active_ingredient IN ?", activeIngredientMatches)
	}

	if isOnCostCatalog {
		query = query.Joins("JOIN medicine_costs mc ON mc.medicine_id = medicines.id AND mc.client_id = ? AND ? BETWEEN mc.start_effective_date AND mc.end_effective_date", clientID, today).Distinct("medicines.*")
	}

	query.Count(&total)

	if isOnCostCatalog {
		query = query.Preload("MedicineCost", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("client_id = ? AND ? BETWEEN start_effective_date AND end_effective_date", clientID, today)
		})
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&medicines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not perform search"})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	c.JSON(http.StatusOK, gin.H{
		"current_page":            page,
		"medicines":               medicines,
		"page_size":               limit,
		"description_like":        descriptionLike,
		"laboratory_like":         laboratoryLike,
		"ean_code_like":           eanCodeLike,
		"sat_key_like":            satKeyLike,
		"active_ingredient_like":  activeIngredientLike,
		"description_match":       descriptionMatches,
		"laboratory_match":        laboratoryMatches,
		"ean_code_match":          eanCodeMatches,
		"sat_key_match":           satKeyMatches,
		"active_ingredient_match": activeIngredientMatches,
		"total_pages":             totalPages,
		"total_records":           total,
	})
}

func SearchMedicineCoincidencesByProperty(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property or search_text"})
		return
	}

	var results []string
	if err := db.DB.
		Model(&db.Medicine{}).
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
