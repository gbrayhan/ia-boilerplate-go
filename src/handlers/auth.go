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

	if h.Auth.ComparePasswords(user.HashPassword, loginRequest.Password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := h.Auth.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	refreshToken, err := h.Auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"email":        user.Email,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
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

	jwt, err := h.Auth.CheckRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	claims, err := h.Auth.GetClaims(jwt)
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

	accessToken, err := h.Auth.GenerateAccessToken(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not authenticate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
	})
}
