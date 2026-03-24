package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义JWT相关错误
var (
	ErrInvalidToken = errors.New("无效的令牌")
	ErrExpiredToken = errors.New("令牌已过期")
)

const (
	defaultJWTSecret = "your-secret-key"
	jwtIssuer        = "seckill_go"
	jwtTTL           = 24 * time.Hour
)

// 定义JWT声明结构
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 生成JWT令牌
func GenerateToken(userID uint) (string, error) {
	now := time.Now()
	expirationTime := now.Add(jwtTTL)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jwtIssuer,
			Subject:   "user_authentication",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(getJWTSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 验证JWT令牌
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(getJWTSecret()), nil
	}, jwt.WithIssuer(jwtIssuer), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return defaultJWTSecret
	}
	return secret
}
