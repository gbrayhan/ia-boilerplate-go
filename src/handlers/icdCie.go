package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"ia-boilerplate/src/repository"
	"net/http"
	"strconv"
)

func (h *Handler) GetICDCies(c *gin.Context) {
	var records []repository.ICDCie
	if result := h.Repository.DB.Find(&records); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve ICDCie records"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *Handler) GetICDCie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var record repository.ICDCie
	if result := h.Repository.DB.First(&record, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ICDCie record not found"})
		return
	}
	c.JSON(http.StatusOK, record)
}

type CreateICDCieRequest struct {
	CieVersion   string `json:"cieVersion" binding:"required"`
	Code         string `json:"code" binding:"required"`
	Description  string `json:"description"`
	ChapterNo    string `json:"chapterNo"`
	ChapterTitle string `json:"chapterTitle"`
}

func (h *Handler) CreateICDCie(c *gin.Context) {
	var req CreateICDCieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existing repository.ICDCie
	if err := h.Repository.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Could not create ICDCie record: duplicate code"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create ICDCie record"})
		return
	}
	record := repository.ICDCie{
		CieVersion:   req.CieVersion,
		Code:         req.Code,
		Description:  req.Description,
		ChapterNo:    req.ChapterNo,
		ChapterTitle: req.ChapterTitle,
	}
	if result := h.Repository.DB.Create(&record); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create ICDCie record"})
		return
	}
	c.JSON(http.StatusCreated, record)
}

type UpdateICDCieRequest struct {
	CieVersion   string `json:"cieVersion" binding:"required"`
	Code         string `json:"code" binding:"required"`
	Description  string `json:"description"`
	ChapterNo    string `json:"chapterNo"`
	ChapterTitle string `json:"chapterTitle"`
}

func (h *Handler) UpdateICDCie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req UpdateICDCieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var record repository.ICDCie
	if res := h.Repository.DB.First(&record, id); res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ICDCie record not found"})
		return
	}
	if req.Code != record.Code {
		var existing repository.ICDCie
		if err := h.Repository.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Could not update ICDCie record: duplicate code"})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update ICDCie record"})
			return
		}
	}
	record.CieVersion = req.CieVersion
	record.Code = req.Code
	record.Description = req.Description
	record.ChapterNo = req.ChapterNo
	record.ChapterTitle = req.ChapterTitle
	if save := h.Repository.DB.Save(&record); save.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update ICDCie record"})
		return
	}
	c.JSON(http.StatusOK, record)
}

func (h *Handler) DeleteICDCie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if res := h.Repository.DB.Delete(&repository.ICDCie{}, id); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete ICDCie record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ICDCie record deleted successfully"})
}

func (h *Handler) SearchICDCiePaginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	likeFilters := map[string]string{
		"cie_version":   c.Query("cie_version_like"),
		"code":          c.Query("code_like"),
		"description":   c.Query("description_like"),
		"chapter_no":    c.Query("chapter_no_like"),
		"chapter_title": c.Query("chapter_title_like"),
	}

	matches := map[string][]string{
		"cie_version":   c.QueryArray("cie_version_match"),
		"code":          c.QueryArray("code_match"),
		"description":   c.QueryArray("description_match"),
		"chapter_no":    c.QueryArray("chapter_no_match"),
		"chapter_title": c.QueryArray("chapter_title_match"),
	}

	var total int64
	query := h.Repository.DB.Model(&repository.ICDCie{})

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
	var records []repository.ICDCie
	if res := query.Offset(offset).Limit(limit).Find(&records); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	c.JSON(http.StatusOK, gin.H{
		"current_page":  page,
		"records":       records,
		"page_size":     limit,
		"total_pages":   totalPages,
		"total_records": total,
	})
}

func (h *Handler) SearchIcdCoincidencesByProperty(c *gin.Context) {
	property := c.Query("property")
	searchText := c.Query("search_text")
	allowed := map[string]bool{"cie_version": true, "code": true, "description": true, "chapter_no": true, "chapter_title": true}
	if !allowed[property] || searchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property or search_text"})
		return
	}
	var results []string
	if res := h.Repository.DB.Model(&repository.ICDCie{}).
		Distinct(property).
		Where(property+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(property, &results); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	c.JSON(http.StatusOK, results)
}
