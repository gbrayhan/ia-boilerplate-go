package infrastructure

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strconv"
	"time"
)

func generateToken(userID int, issuer string, accessSecretKey string, accessTokenTTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"iss":     issuer,
		"exp":     time.Now().Add(accessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(accessSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return signedToken, nil
}

func GenerateAccessToken(userId int) (string, error) {
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		return "", fmt.Errorf("JWT_ISSUER environment variable is not set")
	}

	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	if accessSecretKey == "" {
		return "", fmt.Errorf("ACCESS_SECRET_KEY environment variable is not set")
	}

	accessTokenTTL, err := getEnvAsDuration("ACCESS_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to parse ACCESS_TOKEN_TTL: %w", err)
	}
	token, err := generateToken(userId, issuer, accessSecretKey, accessTokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateRefreshToken(userId int) (string, error) {
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		return "", fmt.Errorf("JWT_ISSUER environment variable is not set")
	}
	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	if refreshSecretKey == "" {
		return "", fmt.Errorf("REFRESH_SECRET_KEY environment variable is not set")
	}

	refreshTokenTTL, err := getEnvAsDuration("REFRESH_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to parse REFRESH_TOKEN_TTL: %w", err)
	}
	token, err := generateToken(userId, issuer, refreshSecretKey, refreshTokenTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func CheckAccessToken(tokenString string) (*jwt.Token, error) {
	accessSecretKey := os.Getenv("ACCESS_SECRET_KEY")
	if accessSecretKey == "" {
		return nil, fmt.Errorf("ACCESS_SECRET_KEY environment variable is not set")
	}

	return checkToken(tokenString, accessSecretKey)
}

func CheckRefreshToken(tokenString string) (*jwt.Token, error) {
	refreshSecretKey := os.Getenv("REFRESH_SECRET_KEY")
	if refreshSecretKey == "" {
		return nil, fmt.Errorf("REFRESH_SECRET_KEY environment variable is not set")
	}

	return checkToken(tokenString, refreshSecretKey)
}

func GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func checkToken(tokenString, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func getEnvAsDuration(key string, defaultVal time.Duration) (time.Duration, error) {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal, nil
	}

	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}

	return time.Duration(valInt) * time.Minute, nil
}
