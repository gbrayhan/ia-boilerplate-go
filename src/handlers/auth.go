package handlers

import (
	"github.com/gin-gonic/gin"
	"ia-boilerplate/src/repository"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user repository.User
	result := h.Repository.DB.Preload("Role").Preload("Devices").Where("email = ?", loginRequest.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	if h.Infrastructure.ComparePasswords(user.HashPassword, loginRequest.Password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := h.Infrastructure.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	refreshToken, err := h.Infrastructure.GenerateRefreshToken(user.ID)
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

func (h *Handler) AccessTokenByRefreshToken(c *gin.Context) {
	var request AccessTokenByRefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := h.Infrastructure.CheckRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	claims, err := h.Infrastructure.GetClaims(jwt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	userID := claims["user_id"].(float64)
	userIDInt := int(userID)

	var user repository.User
	result := h.Repository.DB.Preload("Role").Preload("Devices").First(&user, userIDInt)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	accessToken, err := h.Infrastructure.GenerateAccessToken(userIDInt)
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
