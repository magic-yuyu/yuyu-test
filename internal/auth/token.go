package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构
type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌（支持HS256/RS256）
func GenerateToken(userID, tenantID, email, algorithm string, privateKey string, expiration time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	var token *jwt.Token
	var signed string
	var err error

	switch algorithm {
	case "HS256":
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err = token.SignedString([]byte(privateKey))
	case "RS256":
		priv, err2 := parseRSAPrivateKeyFromPEM([]byte(privateKey))
		if err2 != nil {
			return "", fmt.Errorf("invalid RSA private key: %w", err2)
		}
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		signed, err = token.SignedString(priv)
	default:
		return "", fmt.Errorf("unsupported JWT algorithm: %s", algorithm)
	}
	if err != nil {
		return "", err
	}
	return signed, nil
}

// ParseToken 解析JWT令牌（支持HS256/RS256）
func ParseToken(tokenString, algorithm, publicKey string) (*Claims, error) {
	var keyFunc jwt.Keyfunc

	switch algorithm {
	case "HS256":
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(publicKey), nil
		}
	case "RS256":
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			pub, err := parseRSAPublicKeyFromPEM([]byte(publicKey))
			if err != nil {
				return nil, fmt.Errorf("invalid RSA public key: %w", err)
			}
			return pub, nil
		}
	default:
		return nil, fmt.Errorf("unsupported JWT algorithm: %s", algorithm)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// parseRSAPrivateKeyFromPEM 解析PEM格式RSA私钥
func parseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// parseRSAPublicKeyFromPEM 解析PEM格式RSA公钥
func parseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || (block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY") {
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if rsaPub, ok := pub.(*rsa.PublicKey); ok {
		return rsaPub, nil
	}
	return nil, errors.New("not RSA public key")
}
