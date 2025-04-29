package infrastructure

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Logger *Logger
}

func NewAuth(logger *Logger) *Auth {
	return &Auth{
		Logger: logger,
	}
}

// generateToken creates a JWT token with the given user ID, issuer, secret key, and TTL
func (a *Auth) generateToken(userID int, issuer, secretKey string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"iss":     issuer,
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		a.Logger.Error("Failed to generate token", zap.Error(err))
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return signed, nil
}

// GenerateAccessToken issues a JWT access token for the given user ID
func (a *Auth) GenerateAccessToken(userID int) (string, error) {
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		a.Logger.Error("Environment variable not set", zap.String("var", "JWT_ISSUER"))
		return "", errors.New("JWT_ISSUER environment variable is not set")

	}

	secret := os.Getenv("ACCESS_SECRET_KEY")
	if secret == "" {
		a.Logger.Error("Environment variable not set", zap.String("var", "ACCESS_SECRET_KEY"))
		return "", fmt.Errorf("ACCESS_SECRET_KEY environment variable is not set")
	}

	ttl, err := a.getEnvAsDuration("ACCESS_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		a.Logger.Error("Failed to parse ACCESS_TOKEN_TTL", zap.Error(err))
		return "", fmt.Errorf("failed to parse ACCESS_TOKEN_TTL: %w", err)
	}

	tok, err := a.generateToken(userID, issuer, secret, ttl)
	if err != nil {
		a.Logger.Error("Failed to generate access token", zap.Error(err))
		return "", err
	}
	return tok, nil
}

// GenerateRefreshToken issues a JWT refresh token for the given user ID
func (a *Auth) GenerateRefreshToken(userID int) (string, error) {
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		a.Logger.Error("Environment variable not set", zap.String("var", "JWT_ISSUER"))
		return "", fmt.Errorf("JWT_ISSUER environment variable is not set")
	}

	secret := os.Getenv("REFRESH_SECRET_KEY")
	if secret == "" {
		a.Logger.Error("Environment variable not set", zap.String("var", "REFRESH_SECRET_KEY"))
		return "", fmt.Errorf("REFRESH_SECRET_KEY environment variable is not set")
	}

	ttl, err := a.getEnvAsDuration("REFRESH_TOKEN_TTL", 7*24*time.Hour)
	if err != nil {
		a.Logger.Error("Failed to parse REFRESH_TOKEN_TTL", zap.Error(err))
		return "", fmt.Errorf("failed to parse REFRESH_TOKEN_TTL: %w", err)
	}

	tok, err := a.generateToken(userID, issuer, secret, ttl)
	if err != nil {
		return "", err
	}
	return tok, nil
}

// CheckAccessToken validates the access token string
func (a *Auth) CheckAccessToken(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("ACCESS_SECRET_KEY")
	if secret == "" {
		a.Logger.Error("Environment variable not set", zap.String("var", "ACCESS_SECRET_KEY"))
		return nil, fmt.Errorf("ACCESS_SECRET_KEY environment variable is not set")
	}
	return a.checkToken(tokenString, secret)
}

// CheckRefreshToken validates the refresh token string
func (a *Auth) CheckRefreshToken(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("REFRESH_SECRET_KEY")
	if secret == "" {
		a.Logger.Error("Environment variable not set", zap.String("var", "REFRESH_SECRET_KEY"))
		return nil, fmt.Errorf("REFRESH_SECRET_KEY environment variable is not set")
	}
	return a.checkToken(tokenString, secret)
}

// GetClaims extracts JWT claims as a MapClaims
func (a *Auth) GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		a.Logger.Error("Invalid token claims type")
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// checkToken parses and verifies a JWT token string with the given secret
func (a *Auth) checkToken(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			a.Logger.Error("Unexpected signing method", zap.String("alg", token.Header["alg"].(string)))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		a.Logger.Error("Token validation failed", zap.Error(err))
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		a.Logger.Warn("Invalid token")
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

// getEnvAsDuration reads an env var as integer minutes, with fallback default
func (a *Auth) getEnvAsDuration(key string, defaultVal time.Duration) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}
	minutes, err := strconv.Atoi(val)
	if err != nil {
		a.Logger.Error("Invalid duration value", zap.String("var", key), zap.Error(err))
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}
	return time.Duration(minutes) * time.Minute, nil
}

// HashPassword encrypts a plaintext password using bcrypt
func (a *Auth) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.Logger.Error("Error hashing password", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}

func (a *Auth) ComparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		a.Logger.Warn("Password comparison failed", zap.Error(err))
	}
	return err
}
