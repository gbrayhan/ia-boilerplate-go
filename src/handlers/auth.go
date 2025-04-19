package handlers

import (
	"github.com/gin-gonic/gin"
	"ia-boilerplate/src/db"
	"ia-boilerplate/src/infrastructure"
	"ia-boilerplate/utils"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user db.User
	result := db.DB.Preload("Role").Preload("Devices").Where("email = ?", loginRequest.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	if utils.ComparePasswords(user.HashPassword, loginRequest.Password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := infrastructure.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	refreshToken, err := infrastructure.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":        user.ID,
		"userFirstName": user.FirstName,
		"userLastName":  user.LastName,
		"userEmail":     user.Email,
		"accessToken":   accessToken,
		"refreshToken":  refreshToken,
	})
}

type AccessTokenByRefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func AccessTokenByRefreshToken(c *gin.Context) {
	var request AccessTokenByRefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := infrastructure.CheckRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	claims, err := infrastructure.GetClaims(jwt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	userID := claims["user_id"].(float64)
	userIDInt := int(userID)

	var user db.User
	result := db.DB.Preload("Role").Preload("Devices").First(&user, userIDInt)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	accessToken, err := infrastructure.GenerateAccessToken(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"userId":       user.ID,
		"userUsername": user.Username,
		"userEmail":    user.Email,
	})
}
