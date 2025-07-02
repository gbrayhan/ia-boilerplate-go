package handlers

import (
	"ia-boilerplate/src/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRoles(c *gin.Context) {
	var roles []repository.RoleUser
	result := h.Repository.DB.Find(&roles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *Handler) GetRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var role repository.RoleUser
	result := h.Repository.DB.First(&role, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newRole := repository.RoleUser{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	result := h.Repository.DB.Create(&newRole)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create role"})
		return
	}
	c.JSON(http.StatusCreated, newRole)
}

type UpdateRoleRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Enabled     *bool   `json:"enabled"`
}

func (h *Handler) UpdateRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID, must be an integer"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role repository.RoleUser
	result := h.Repository.DB.First(&role, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Preparar los campos a actualizar
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	// Actualizar solo los campos que están presentes en el request
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	// Si no hay campos para actualizar, retornar error
	if len(updates) <= 1 { // Solo updated_at
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Realizar la actualización
	if err := h.Repository.DB.Model(&repository.RoleUser{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update role"})
		return
	}

	// Obtener el rol actualizado
	var updatedRole repository.RoleUser
	if err := h.Repository.DB.First(&updatedRole, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve updated role"})
		return
	}

	c.JSON(http.StatusOK, updatedRole)
}

func (h *Handler) DeleteRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result := h.Repository.DB.Delete(&repository.RoleUser{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	JobPosition string `json:"jobPosition"`
	RoleID      int    `json:"roleId" binding:"required"`
	Enabled     bool   `json:"enabled"`
}

type UpdateUserRequest struct {
	Username    *string `json:"username"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Email       *string `json:"email"`
	Password    *string `json:"password"`
	JobPosition *string `json:"jobPosition"`
	RoleID      *int    `json:"roleId"`
	Enabled     *bool   `json:"enabled"`
}

func (h *Handler) GetUsers(c *gin.Context) {
	var users []repository.User
	result := h.Repository.DB.Preload("Role").Preload("Devices").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var user repository.User
	result := h.Repository.DB.Preload("Role").Preload("Devices").First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := h.Auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
		return
	}
	newUser := repository.User{
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		HashPassword: hashedPassword,
		JobPosition:  req.JobPosition,
		RoleID:       req.RoleID,
		Enabled:      req.Enabled,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	result := h.Repository.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el usuario existe
	var existingUser repository.User
	if err := h.Repository.DB.Preload("Role").Preload("Devices").First(&existingUser, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Preparar los campos a actualizar
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	// Validar y agregar cada campo si está presente en el request
	if req.Username != nil {
		updates["username"] = *req.Username
	}

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}

	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}

	if req.Email != nil {
		updates["email"] = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := h.Auth.HashPassword(*req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
			return
		}
		updates["hash_password"] = hashedPassword
	}

	if req.JobPosition != nil {
		updates["job_position"] = *req.JobPosition
	}

	if req.RoleID != nil {
		updates["role_id"] = *req.RoleID
	}

	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	// Si no hay campos para actualizar, retornar error
	if len(updates) <= 1 { // Solo updated_at
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Realizar la actualización
	if err := h.Repository.DB.Model(&repository.User{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	// Obtener el usuario actualizado
	var updatedUser repository.User
	if err := h.Repository.DB.Preload("Role").Preload("Devices").First(&updatedUser, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve updated user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result := h.Repository.DB.Delete(&repository.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *Handler) SearchUsersPaginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	likeFilters := map[string]string{
		"username":     c.Query("username_like"),
		"first_name":   c.Query("first_name_like"),
		"last_name":    c.Query("last_name_like"),
		"email":        c.Query("email_like"),
		"job_position": c.Query("job_position_like"),
	}

	matches := map[string][]string{
		"username":     c.QueryArray("username_match"),
		"first_name":   c.QueryArray("first_name_match"),
		"last_name":    c.QueryArray("last_name_match"),
		"email":        c.QueryArray("email_match"),
		"job_position": c.QueryArray("job_position_match"),
	}

	var (
		total int64
		users []repository.User
	)

	query := h.Repository.DB.Model(&repository.User{})

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
	query = query.Preload("Role").Preload("Devices")

	offset := (page - 1) * limit
	if res := query.Offset(offset).Limit(limit).Find(&users); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	c.JSON(http.StatusOK, gin.H{
		"current_page":  page,
		"users":         users,
		"page_size":     limit,
		"total_pages":   totalPages,
		"total_records": total,
	})
}

func (h *Handler) SearchUserCoincidencesByProperty(c *gin.Context) {
	property := c.Query("property")
	searchText := c.Query("search_text")

	allowed := map[string]bool{
		"username":     true,
		"first_name":   true,
		"last_name":    true,
		"email":        true,
		"job_position": true,
	}
	if !allowed[property] || searchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property or search_text"})
		return
	}

	var results []string
	if res := h.Repository.DB.
		Model(&repository.User{}).
		Distinct(property).
		Where(property+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(property, &results); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}

type CreateDeviceRequest struct {
	UserID         int    `json:"userId" binding:"required"`
	IPAddress      string `json:"ip_address" binding:"required"`
	UserAgent      string `json:"user_agent"`
	DeviceType     string `json:"device_type"`
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	OS             string `json:"os"`
	Language       string `json:"language"`
}

type UpdateDeviceRequest struct {
	IPAddress      *string `json:"ip_address"`
	UserAgent      *string `json:"user_agent"`
	DeviceType     *string `json:"device_type"`
	Browser        *string `json:"browser"`
	BrowserVersion *string `json:"browser_version"`
	OS             *string `json:"os"`
	Language       *string `json:"language"`
}

func (h *Handler) GetDevicesByUser(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var devices []repository.DeviceDetails
	result := h.Repository.DB.Where("user_id = ?", userID).Find(&devices)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve devices"})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *Handler) GetDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var device repository.DeviceDetails
	result := h.Repository.DB.First(&device, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (h *Handler) CreateDevice(c *gin.Context) {
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newDevice := repository.DeviceDetails{
		UserID:         req.UserID,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
		DeviceType:     req.DeviceType,
		Browser:        req.Browser,
		BrowserVersion: req.BrowserVersion,
		OS:             req.OS,
		Language:       req.Language,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	result := h.Repository.DB.Create(&newDevice)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create device"})
		return
	}
	c.JSON(http.StatusCreated, newDevice)
}

func (h *Handler) UpdateDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el dispositivo existe
	var existingDevice repository.DeviceDetails
	if err := h.Repository.DB.First(&existingDevice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	// Preparar los campos a actualizar
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	// Validar y agregar cada campo si está presente en el request
	if req.IPAddress != nil {
		updates["ip_address"] = *req.IPAddress
	}

	if req.UserAgent != nil {
		updates["user_agent"] = *req.UserAgent
	}

	if req.DeviceType != nil {
		updates["device_type"] = *req.DeviceType
	}

	if req.Browser != nil {
		updates["browser"] = *req.Browser
	}

	if req.BrowserVersion != nil {
		updates["browser_version"] = *req.BrowserVersion
	}

	if req.OS != nil {
		updates["os"] = *req.OS
	}

	if req.Language != nil {
		updates["language"] = *req.Language
	}

	// Si no hay campos para actualizar, retornar error
	if len(updates) <= 1 { // Solo updated_at
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Realizar la actualización
	if err := h.Repository.DB.Model(&repository.DeviceDetails{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update device"})
		return
	}

	// Obtener el dispositivo actualizado
	var updatedDevice repository.DeviceDetails
	if err := h.Repository.DB.First(&updatedDevice, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve updated device"})
		return
	}

	c.JSON(http.StatusOK, updatedDevice)
}

func (h *Handler) DeleteDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result := h.Repository.DB.Delete(&repository.DeviceDetails{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete device"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

func (h *Handler) SearchDeviceDetailsPaginated(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	likeFilters := map[string]string{
		"ip_address":      c.Query("ip_address_like"),
		"user_agent":      c.Query("user_agent_like"),
		"device_type":     c.Query("device_type_like"),
		"browser":         c.Query("browser_like"),
		"browser_version": c.Query("browser_version_like"),
		"os":              c.Query("os_like"),
		"language":        c.Query("language_like"),
	}

	matches := map[string][]string{
		"ip_address":      c.QueryArray("ip_address_match"),
		"user_agent":      c.QueryArray("user_agent_match"),
		"device_type":     c.QueryArray("device_type_match"),
		"browser":         c.QueryArray("browser_match"),
		"browser_version": c.QueryArray("browser_version_match"),
		"os":              c.QueryArray("os_match"),
		"language":        c.QueryArray("language_match"),
	}

	var (
		total   int64
		records []repository.DeviceDetails
	)

	query := h.Repository.DB.
		Model(&repository.DeviceDetails{})

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

func (h *Handler) SearchDeviceCoincidencesByProperty(c *gin.Context) {
	property := c.Query("property")
	searchText := c.Query("search_text")

	allowed := map[string]bool{
		"ip_address":      true,
		"user_agent":      true,
		"device_type":     true,
		"browser":         true,
		"browser_version": true,
		"os":              true,
		"language":        true,
	}
	if !allowed[property] || searchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property or search_text"})
		return
	}

	var results []string
	if res := h.Repository.DB.
		Model(&repository.DeviceDetails{}).
		Distinct(property).
		Where(property+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(property, &results); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}
