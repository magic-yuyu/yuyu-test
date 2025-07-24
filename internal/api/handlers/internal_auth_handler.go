package handlers

import (
	"net/http"
	"strings"
	"time"

	"yuyu-test/internal/store/database"

	"encoding/base64"

	authpkg "yuyu-test/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// InternalAuthHandler 对内服务认证处理器
// 只实现/oauth/token端点

type InternalAuthHandler struct {
	db     database.Querier
	signer authpkg.JWTSigner
}

func NewInternalAuthHandler(db database.Querier, signer authpkg.JWTSigner) *InternalAuthHandler {
	return &InternalAuthHandler{
		db:     db,
		signer: signer,
	}
}

// POST /oauth/token
// Basic Auth: client_id/client_secret
// grant_type=client_credentials
func (h *InternalAuthHandler) Token(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Basic authorization required"})
		return
	}
	// 解码Basic Auth
	payload, err := decodeBasicAuth(auth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid basic auth"})
		return
	}
	clientID, clientSecret, ok := payload[0], payload[1], len(payload) == 2
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid basic auth format"})
		return
	}
	// 查找client
	client, err := h.db.GetInternalClientByID(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Client not found"})
		return
	}
	// 校验secret
	if bcrypt.CompareHashAndPassword([]byte(client.ClientSecretHash), []byte(clientSecret)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid client secret"})
		return
	}
	// grant_type
	grantType := c.PostForm("grant_type")
	if grantType != "client_credentials" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "grant_type must be client_credentials"})
		return
	}
	// 查询scope
	scopes, err := h.db.GetClientScopes(c.Request.Context(), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get scopes"})
		return
	}
	scopeNames := make([]string, 0, len(scopes))
	for _, s := range scopes {
		scopeNames = append(scopeNames, s.ScopeName)
	}
	// 签发JWT
	expiresIn := 300 // 5分钟
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   clientID,
		"scope": strings.Join(scopeNames, " "),
		"exp":   now.Add(time.Second * time.Duration(expiresIn)).Unix(),
		"iat":   now.Unix(),
		"iss":   "https://auth.yoursaas.com",
	}
	tokenStr, err := h.signer.Sign(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenStr,
		"token_type":   "Bearer",
		"expires_in":   expiresIn,
	})
}

// decodeBasicAuth 解析Basic Auth
func decodeBasicAuth(auth string) ([]string, error) {
	// "Basic base64(client_id:client_secret)"
	payload := strings.TrimPrefix(auth, "Basic ")
	decoded, err := decodeBase64(payload)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(decoded, ":", 2)
	return parts, nil
}

func decodeBase64(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
