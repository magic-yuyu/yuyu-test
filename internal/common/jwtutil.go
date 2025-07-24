package common

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// getSigningMethod 根据算法返回签名方法
func GetSigningMethod(alg string) jwt.SigningMethod {
	switch alg {
	case "RS256":
		return jwt.SigningMethodRS256
	case "ES256":
		return jwt.SigningMethodES256
	default:
		return jwt.SigningMethodHS256
	}
}

// signJWTToken 用指定算法和私钥签名
func SignJWTToken(token *jwt.Token, alg, privateKey string) (string, error) {
	switch alg {
	case "RS256":
		priv, err := ParseRSAPrivateKeyFromPEM([]byte(privateKey))
		if err != nil {
			return "", err
		}
		return token.SignedString(priv)
	case "ES256":
		priv, err := ParseECPrivateKeyFromPEM([]byte(privateKey))
		if err != nil {
			return "", err
		}
		return token.SignedString(priv)
	default:
		return token.SignedString([]byte(privateKey))
	}
}

// ParseRSAPrivateKeyFromPEM 解析PEM格式RSA私钥
func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// ParseRSAPublicKeyFromPEM 解析PEM格式RSA公钥
func ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
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

// ParseECPrivateKeyFromPEM 解析PEM格式EC私钥
func ParseECPrivateKeyFromPEM(key []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing EC private key")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}

// ParseECPublicKeyFromPEM 解析PEM格式EC公钥
func ParseECPublicKeyFromPEM(key []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing EC public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if ecPub, ok := pub.(*ecdsa.PublicKey); ok {
		return ecPub, nil
	}
	return nil, errors.New("not EC public key")
}
