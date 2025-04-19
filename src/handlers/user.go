package handlers

import (
	"github.com/gin-gonic/gin"
	"ia-boilerplate/src/db"
	"ia-boilerplate/utils"
	"net/http"
	"strconv"
	"time"
)

func GetRoles(c *gin.Context) {
	var roles []db.RoleUser
	result := db.DB.Find(&roles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func GetRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var role db.RoleUser
	result := db.DB.First(&role, id)
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

func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newRole := db.RoleUser{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	result := db.DB.Create(&newRole)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create role"})
		return
	}
	c.JSON(http.StatusCreated, newRole)
}

type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func UpdateRole(c *gin.Context) {
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
	updatedData := db.RoleUser{
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		UpdatedAt:   time.Now(),
	}

	var role db.RoleUser
	result := db.DB.First(&role, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	role.Name = updatedData.Name
	role.Description = updatedData.Description
	role.Enabled = updatedData.Enabled
	role.UpdatedAt = time.Now()
	saveResult := db.DB.Save(&role)

	if saveResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update role"})
		return
	}
	c.JSON(http.StatusOK, role)
}

func DeleteRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result := db.DB.Delete(&db.RoleUser{}, id)
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
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password"`
	JobPosition string `json:"jobPosition"`
	RoleID      int    `json:"roleId" binding:"required"`
	Enabled     bool   `json:"enabled"`
}

func GetUsers(c *gin.Context) {
	var users []db.User
	result := db.DB.Preload("Role").Preload("Devices").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var user db.User
	result := db.DB.Preload("Role").Preload("Devices").First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
		return
	}
	newUser := db.User{
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
	result := db.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	c.JSON(http.StatusCreated, newUser)
}

func UpdateUser(c *gin.Context) {
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
	var user db.User
	result := db.DB.Preload("Role").Preload("Devices").First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Username = req.Username
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
			return
		}
		user.HashPassword = hashedPassword
	}
	user.JobPosition = req.JobPosition
	user.RoleID = req.RoleID
	user.Enabled = req.Enabled
	user.UpdatedAt = time.Now()
	saveResult := db.DB.Save(&user)
	if saveResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result := db.DB.Delete(&db.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func SearchUsersPaginated(c *gin.Context) {
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
		users []db.User
	)

	query := db.DB.Model(&db.User{})

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

func SearchUserCoincidencesByProperty(c *gin.Context) {
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
	if res := db.DB.
		Model(&db.User{}).
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
	IPAddress      string `json:"ip_address"`
	UserAgent      string `json:"user_agent"`
	DeviceType     string `json:"device_type"`
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	OS             string `json:"os"`
	Language       string `json:"language"`
}

func GetDevicesByUser(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var devices []db.DeviceDetails
	result := db.DB.Where("user_id = ?", userID).Find(&devices)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve devices"})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func GetDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var device db.DeviceDetails
	result := db.DB.First(&device, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func CreateDevice(c *gin.Context) {
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newDevice := db.DeviceDetails{
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
	result := db.DB.Create(&newDevice)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create device"})
		return
	}
	c.JSON(http.StatusCreated, newDevice)
}

func UpdateDevice(c *gin.Context) {
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
	var device db.DeviceDetails
	result := db.DB.First(&device, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	device.IPAddress = req.IPAddress
	device.UserAgent = req.UserAgent
	device.DeviceType = req.DeviceType
	device.Browser = req.Browser
	device.BrowserVersion = req.BrowserVersion
	device.OS = req.OS
	device.Language = req.Language
	device.UpdatedAt = time.Now()
	saveResult := db.DB.Save(&device)
	if saveResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update device"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func DeleteDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result := db.DB.Delete(&db.DeviceDetails{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete device"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

func SearchDeviceDetailsPaginated(c *gin.Context) {
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
		records []db.DeviceDetails
	)

	query := db.DB.
		Model(&db.DeviceDetails{})

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

func SearchDeviceCoincidencesByProperty(c *gin.Context) {
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
	if res := db.DB.
		Model(&db.DeviceDetails{}).
		Distinct(property).
		Where(property+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(property, &results); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}
