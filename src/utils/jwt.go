package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserUUID string `json:"user_uuid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// SetJWTSecret 设置JWT密钥
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateJWT 生成JWT令牌
func GenerateJWT(userUUID, username string, isAdmin bool, expirationTime time.Duration) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret not configured")
	}

	claims := &JWTClaims{
		UserUUID: userUUID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT 验证JWT令牌
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret not configured")
	}

	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshJWT 刷新JWT令牌
func RefreshJWT(tokenString string, newExpirationTime time.Duration) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新的令牌
	return GenerateJWT(claims.UserUUID, claims.Username, claims.IsAdmin, newExpirationTime)
}