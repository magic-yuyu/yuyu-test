package middleware

import (
	"net/http"
	"strings"

	"yuyu-test/internal/auth"
	"yuyu-test/internal/tenant"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	tenantService *tenant.Service
	signer        auth.JWTSigner
}

// NewAuthMiddleware 创建新的认证中间件
func NewAuthMiddleware(tenantService *tenant.Service, signer auth.JWTSigner) *AuthMiddleware {
	return &AuthMiddleware{
		tenantService: tenantService,
		signer:        signer,
	}
}

// APIKeyAuth API密钥认证中间件
func (m *AuthMiddleware) APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取API密钥
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		apiKey := parts[1]

		// 验证API密钥
		tenant, err := m.tenantService.ValidateAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		// 将租户信息存储到上下文中
		c.Set("tenant", tenant)
		c.Next()
	}
}

// JWTAuth JWT认证中间件
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取JWT令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析JWT令牌
		claims := &auth.Claims{}
		err := m.signer.Parse(tokenString, claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
