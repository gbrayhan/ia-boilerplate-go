package handlers

import (
	"github.com/gin-gonic/gin"
	"ia-boilerplate/src/db"
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
	Role        string `json:"role" binding:"required"`
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
		Role:        req.Role,
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
	Role        string `json:"role" binding:"required"`
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
		Role:        req.Role,
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
	role.Role = updatedData.Role
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
	FullName    string `json:"fullName"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	JobPosition string `json:"jobPosition"`
	RoleID      int    `json:"roleId" binding:"required"`
	Enabled     bool   `json:"enabled"`
}

type UpdateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	FullName    string `json:"fullName"`
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
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encrypting password"})
		return
	}
	newUser := db.User{
		Username:     req.Username,
		FullName:     req.FullName,
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
	user.FullName = req.FullName
	user.Email = req.Email
	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
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
