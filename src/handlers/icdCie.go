package handlers

import (
	"github.com/gin-gonic/gin"
	"ia-boilerplate/src/db"
	"net/http"
	"strconv"
)

func GetICDCies(c *gin.Context) {
	var records []db.ICDCie
	if result := db.DB.Find(&records); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve ICDCie records"})
		return
	}
	c.JSON(http.StatusOK, records)
}

func GetICDCie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var record db.ICDCie
	if result := db.DB.First(&record, id); result.Error != nil {
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

func CreateICDCie(c *gin.Context) {
	var req CreateICDCieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record := db.ICDCie{
		CieVersion:   req.CieVersion,
		Code:         req.Code,
		Description:  req.Description,
		ChapterNo:    req.ChapterNo,
		ChapterTitle: req.ChapterTitle,
	}
	if result := db.DB.Create(&record); result.Error != nil {
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

func UpdateICDCie(c *gin.Context) {
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
	var record db.ICDCie
	if res := db.DB.First(&record, id); res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ICDCie record not found"})
		return
	}
	record.CieVersion = req.CieVersion
	record.Code = req.Code
	record.Description = req.Description
	record.ChapterNo = req.ChapterNo
	record.ChapterTitle = req.ChapterTitle
	if save := db.DB.Save(&record); save.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update ICDCie record"})
		return
	}
	c.JSON(http.StatusOK, record)
}

func DeleteICDCie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if res := db.DB.Delete(&db.ICDCie{}, id); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete ICDCie record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ICDCie record deleted successfully"})
}

func SearchICDCiePaginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	cieLike := c.Query("cie_version_like")
	codeLike := c.Query("code_like")
	descLike := c.Query("description_like")
	chapNoLike := c.Query("chapter_no_like")
	chapTitleLike := c.Query("chapter_title_like")
	matches := map[string][]string{
		"cie_version":   c.QueryArray("cie_version_match"),
		"code":          c.QueryArray("code_match"),
		"description":   c.QueryArray("description_match"),
		"chapter_no":    c.QueryArray("chapter_no_match"),
		"chapter_title": c.QueryArray("chapter_title_match"),
	}
	var total int64
	query := db.DB.Model(&db.ICDCie{})
	if cieLike != "" {
		query = query.Where("cie_version ILIKE ?", "%"+cieLike+"%")
	}
	if codeLike != "" {
		query = query.Where("code ILIKE ?", "%"+codeLike+"%")
	}
	if descLike != "" {
		query = query.Where("description ILIKE ?", "%"+descLike+"%")
	}
	if chapNoLike != "" {
		query = query.Where("chapter_no ILIKE ?", "%"+chapNoLike+"%")
	}
	if chapTitleLike != "" {
		query = query.Where("chapter_title ILIKE ?", "%"+chapTitleLike+"%")
	}
	for col, vals := range matches {
		if len(vals) > 0 {
			query = query.Where(""+col+" IN (?)", vals)
		}
	}
	query.Count(&total)
	offset := (page - 1) * limit
	var records []db.ICDCie
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

func SearchIcdCoincidencesByProperty(c *gin.Context) {
	property := c.Query("property")
	searchText := c.Query("search_text")
	allowed := map[string]bool{"cie_version": true, "code": true, "description": true, "chapter_no": true, "chapter_title": true}
	if !allowed[property] || searchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property or search_text"})
		return
	}
	var results []string
	if res := db.DB.Model(&db.ICDCie{}).
		Distinct(property).
		Where(property+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(property, &results); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	c.JSON(http.StatusOK, results)
}
