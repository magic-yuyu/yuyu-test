package auth

import (
	"crypto/ecdsa"
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

type JWTSigner interface {
	Sign(claims jwt.Claims) (string, error)
	Parse(tokenString string, claims jwt.Claims) error
	Algorithm() string
	PublicKey() interface{}
	PrivateKey() interface{}
}

type HS256Signer struct {
	secret string
}

func NewHS256Signer(secret string) *HS256Signer {
	return &HS256Signer{secret: secret}
}

func (s *HS256Signer) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *HS256Signer) Parse(tokenString string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	return err
}

func (s *HS256Signer) Algorithm() string { return "HS256" }

func (s *HS256Signer) PublicKey() interface{}  { return s.secret }
func (s *HS256Signer) PrivateKey() interface{} { return s.secret }

// RS256Signer

type RS256Signer struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRS256Signer(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *RS256Signer {
	return &RS256Signer{privateKey: privateKey, publicKey: publicKey}
}

func (s *RS256Signer) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *RS256Signer) Parse(tokenString string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})
	return err
}

func (s *RS256Signer) Algorithm() string { return "RS256" }

func (s *RS256Signer) PublicKey() interface{}  { return s.publicKey }
func (s *RS256Signer) PrivateKey() interface{} { return s.privateKey }

// ES256Signer

type ES256Signer struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewES256Signer(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) *ES256Signer {
	return &ES256Signer{privateKey: privateKey, publicKey: publicKey}
}

func (s *ES256Signer) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(s.privateKey)
}

func (s *ES256Signer) Parse(tokenString string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})
	return err
}

func (s *ES256Signer) Algorithm() string { return "ES256" }

func (s *ES256Signer) PublicKey() interface{}  { return s.publicKey }
func (s *ES256Signer) PrivateKey() interface{} { return s.privateKey }
