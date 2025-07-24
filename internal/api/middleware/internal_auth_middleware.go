package middleware

import (
	"net/http"
	"strings"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"

	"yuyu-test/internal/internal_service"
)

// InternalAuthMiddleware 内部服务认证中间件
type InternalAuthMiddleware struct {
	internalService *internal_service.Service
	logger          *slog.Logger
}

// NewInternalAuthMiddleware 创建内部服务认证中间件
func NewInternalAuthMiddleware(internalService *internal_service.Service, logger *slog.Logger) *InternalAuthMiddleware {
	return &InternalAuthMiddleware{
		internalService: internalService,
		logger:          logger,
	}
}

// RequireAuth 要求认证中间件
func (m *InternalAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 从Authorization头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Error("missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		// 验证令牌格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			m.logger.Error("invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// 提取客户端ID
		clientID, err := m.internalService.ExtractClientIDFromToken(authHeader)
		if err != nil {
			m.logger.Error("failed to extract client ID from token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// 验证令牌
		validationReq := internal_service.ValidateTokenRequest{
			Token: strings.TrimPrefix(authHeader, "Bearer "),
		}

		validationResp, err := m.internalService.ValidateToken(c.Request.Context(), validationReq)
		if err != nil {
			m.logger.Error("failed to validate token", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Failed to validate token",
			})
			c.Abort()
			return
		}

		if !validationResp.Valid {
			m.logger.Error("invalid token", "client_id", clientID)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": validationResp.Message,
			})
			c.Abort()
			return
		}

		// 将客户端信息存储到上下文中
		c.Set("client_id", clientID)
		c.Set("scopes", validationResp.Scopes)

		// 记录访问日志
		go func() {
			responseTime := time.Since(start).Milliseconds()
			err := m.internalService.LogAccess(
				c.Request.Context(),
				clientID,
				c.Request.URL.Path,
				c.Request.Method,
				c.Writer.Status(),
				int(responseTime),
				c.ClientIP(),
				c.Request.UserAgent(),
				"", // 请求体（可选）
				"", // 响应体（可选）
			)
			if err != nil {
				m.logger.Error("failed to log access", "error", err)
			}
		}()

		c.Next()
	}
}

// RequireScope 要求特定权限的中间件
func (m *InternalAuthMiddleware) RequireScope(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先进行认证
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// 获取客户端ID
		clientID, exists := c.Get("client_id")
		if !exists {
			m.logger.Error("client_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Client ID not found in context",
			})
			c.Abort()
			return
		}

		// 检查权限
		permissionReq := internal_service.CheckPermissionRequest{
			ClientID:  clientID.(string),
			ScopeName: requiredScope,
		}

		permissionResp, err := m.internalService.CheckPermission(c.Request.Context(), permissionReq)
		if err != nil {
			m.logger.Error("failed to check permission", "error", err, "client_id", clientID, "scope", requiredScope)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Failed to check permission",
			})
			c.Abort()
			return
		}

		if !permissionResp.HasPermission {
			m.logger.Error("permission denied", "client_id", clientID, "scope", requiredScope)
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Forbidden",
				"message":        "Insufficient permissions",
				"required_scope": requiredScope,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyScope 要求任意一个权限的中间件
func (m *InternalAuthMiddleware) RequireAnyScope(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先进行认证
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// 获取客户端ID
		clientID, exists := c.Get("client_id")
		if !exists {
			m.logger.Error("client_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Client ID not found in context",
			})
			c.Abort()
			return
		}

		// 检查是否有任意一个权限
		hasAnyPermission := false
		for _, scope := range requiredScopes {
			permissionReq := internal_service.CheckPermissionRequest{
				ClientID:  clientID.(string),
				ScopeName: scope,
			}

			permissionResp, err := m.internalService.CheckPermission(c.Request.Context(), permissionReq)
			if err != nil {
				m.logger.Error("failed to check permission", "error", err, "client_id", clientID, "scope", scope)
				continue
			}

			if permissionResp.HasPermission {
				hasAnyPermission = true
				break
			}
		}

		if !hasAnyPermission {
			m.logger.Error("permission denied", "client_id", clientID, "required_scopes", requiredScopes)
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "Forbidden",
				"message":         "Insufficient permissions",
				"required_scopes": requiredScopes,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllScopes 要求所有权限的中间件
func (m *InternalAuthMiddleware) RequireAllScopes(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先进行认证
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// 获取客户端ID
		clientID, exists := c.Get("client_id")
		if !exists {
			m.logger.Error("client_id not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Client ID not found in context",
			})
			c.Abort()
			return
		}

		// 检查是否有所有权限
		missingScopes := []string{}
		for _, scope := range requiredScopes {
			permissionReq := internal_service.CheckPermissionRequest{
				ClientID:  clientID.(string),
				ScopeName: scope,
			}

			permissionResp, err := m.internalService.CheckPermission(c.Request.Context(), permissionReq)
			if err != nil {
				m.logger.Error("failed to check permission", "error", err, "client_id", clientID, "scope", scope)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "Failed to check permission",
				})
				c.Abort()
				return
			}

			if !permissionResp.HasPermission {
				missingScopes = append(missingScopes, scope)
			}
		}

		if len(missingScopes) > 0 {
			m.logger.Error("permission denied", "client_id", clientID, "missing_scopes", missingScopes)
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Forbidden",
				"message":        "Insufficient permissions",
				"missing_scopes": missingScopes,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（不强制要求认证）
func (m *InternalAuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// 没有认证信息，继续处理
			c.Next()
			return
		}

		// 尝试提取客户端ID
		clientID, err := m.internalService.ExtractClientIDFromToken(authHeader)
		if err != nil {
			// 令牌无效，但不阻止请求继续
			m.logger.Warn("invalid token in optional auth", "error", err)
			c.Next()
			return
		}

		// 验证令牌
		validationReq := internal_service.ValidateTokenRequest{
			Token: strings.TrimPrefix(authHeader, "Bearer "),
		}

		validationResp, err := m.internalService.ValidateToken(c.Request.Context(), validationReq)
		if err != nil || !validationResp.Valid {
			// 令牌无效，但不阻止请求继续
			m.logger.Warn("invalid token in optional auth", "error", err)
			c.Next()
			return
		}

		// 将客户端信息存储到上下文中
		c.Set("client_id", clientID)
		c.Set("scopes", validationResp.Scopes)
		c.Set("authenticated", true)

		c.Next()
	}
}

// GetClientID 从上下文中获取客户端ID
func GetClientID(c *gin.Context) (string, bool) {
	clientID, exists := c.Get("client_id")
	if !exists {
		return "", false
	}
	return clientID.(string), true
}

// GetScopes 从上下文中获取权限列表
func GetScopes(c *gin.Context) ([]string, bool) {
	scopes, exists := c.Get("scopes")
	if !exists {
		return nil, false
	}
	return scopes.([]string), true
}

// IsAuthenticated 检查是否已认证
func IsAuthenticated(c *gin.Context) bool {
	authenticated, exists := c.Get("authenticated")
	if !exists {
		return false
	}
	return authenticated.(bool)
}
