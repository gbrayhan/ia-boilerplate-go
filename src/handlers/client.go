package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"ia-boilerplate/src/db"
)

type CreateClientRequest struct {
	Alias          string `json:"alias" binding:"required"`
	LegalName      string `json:"legalName" binding:"required"`
	TIN            string `json:"tin" binding:"required"`
	ContractNumber string `json:"contractNumber"`
	FiscalAddress  string `json:"fiscalAddress" binding:"required"`
}

type UpdateClientRequest struct {
	Alias          string `json:"alias" binding:"required"`
	LegalName      string `json:"legalName" binding:"required"`
	TIN            string `json:"tin" binding:"required"`
	FiscalAddress  string `json:"fiscalAddress" binding:"required"`
	ContractNumber string `json:"contractNumber"`
}

func GetClients(c *gin.Context) {
	var clients []db.Client
	result := db.DB.
		Where("is_deleted = ?", 0).
		Preload("SubClients", "is_deleted = ?", 0).
		Preload("SubClients.Programs", "is_deleted = ?", 0).
		Find(&clients)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve clients"})
		return
	}
	c.JSON(http.StatusOK, clients)
}

func GetClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var client db.Client
	result := db.DB.
		Where("is_deleted = ?", 0).
		Preload("SubClients", "is_deleted = ?", 0).
		Preload("SubClients.Programs", "is_deleted = ?", 0).
		First(&client, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	c.JSON(http.StatusOK, client)
}

func CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client := db.Client{
		Alias:          req.Alias,
		LegalName:      req.LegalName,
		TIN:            req.TIN,
		FiscalAddress:  req.FiscalAddress,
		ContractNumber: req.ContractNumber,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if result := db.DB.Create(&client); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create client"})
		return
	}
	c.JSON(http.StatusCreated, client)
}

func UpdateClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var client db.Client
	if result := db.DB.First(&client, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	client.Alias = req.Alias
	client.LegalName = req.LegalName
	client.TIN = req.TIN
	client.FiscalAddress = req.FiscalAddress
	client.ContractNumber = req.ContractNumber
	client.UpdatedAt = time.Now()
	if save := db.DB.Save(&client); save.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update client"})
		return
	}
	c.JSON(http.StatusOK, client)
}

func DeleteClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if result := db.DB.Model(&db.Client{}).
		Where("id = ? AND is_deleted = ?", id, 0).
		Update("is_deleted", 1); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete client"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}

func SearchClientsPaginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	aliasLike := c.Query("alias_like")
	legalNameLike := c.Query("legal_name_like")
	tinLike := c.Query("tin_like")
	addressLike := c.Query("fiscal_address_like")
	contractNumberLike := c.Query("contract_number_like")

	matches := map[string][]string{
		"alias":           c.QueryArray("alias_name_match"),
		"legal_name":      c.QueryArray("legal_name_match"),
		"tin":             c.QueryArray("tin_match"),
		"fiscal_address":  c.QueryArray("fiscal_address_match"),
		"contract_number": c.QueryArray("contract_number_match"),
	}

	var total int64
	query := db.DB.Model(&db.Client{}).
		Where("is_deleted = ?", 0)

	if aliasLike != "" {
		query = query.Where("alias ILIKE ?", "%"+aliasLike+"%")
	}
	if legalNameLike != "" {
		query = query.Where("legal_name ILIKE ?", "%"+legalNameLike+"%")
	}
	if tinLike != "" {
		query = query.Where("tin ILIKE ?", "%"+tinLike+"%")
	}
	if addressLike != "" {
		query = query.Where("fiscal_address ILIKE ?", "%"+addressLike+"%")
	}
	if contractNumberLike != "" {
		query = query.Where("contract_number ILIKE ?", "%"+contractNumberLike+"%")
	}

	for col, vals := range matches {
		if len(vals) > 0 {
			query = query.Where(col+" IN (?)", vals)
		}
	}

	query.Count(&total)

	offset := (page - 1) * limit
	var clients []db.Client
	if res := query.Offset(offset).Limit(limit).Find(&clients); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	c.JSON(http.StatusOK, gin.H{
		"current_page":  page,
		"clients":       clients,
		"page_size":     limit,
		"total_pages":   totalPages,
		"total_records": total,
	})
}

func SearchClientCoincidencesByProperty(c *gin.Context) {
	property := c.Query("property")
	searchText := c.Query("search_text")

	allowed := map[string]bool{
		"alias_name":     true,
		"client_name":    true,
		"rfc":            true,
		"fiscal_address": true,
	}

	if !allowed[property] || searchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property or search_text"})
		return
	}

	var results []string
	if res := db.DB.Model(&db.Client{}).
		Distinct(property).
		Where(property+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(property, &results); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}
