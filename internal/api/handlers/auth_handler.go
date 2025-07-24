package handlers

import (
	"net/http"
	"time"

	"yuyu-test/internal/auth"
	"yuyu-test/internal/store/database"
	"yuyu-test/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService *user.Service
	signer      auth.JWTSigner
}

// NewAuthHandler 创建新的认证处理器
func NewAuthHandler(userService *user.Service, signer auth.JWTSigner) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		signer:      signer,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从中间件获取租户信息
	tenantInterface, exists := c.Get("tenant")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	tenant := tenantInterface.(*database.Tenant)
	response, err := h.userService.Register(c.Request.Context(), tenant.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从中间件获取租户信息
	tenantInterface, exists := c.Get("tenant")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	tenant := tenantInterface.(*database.Tenant)
	response, err := h.userService.Login(c.Request.Context(), tenant.ID, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GenerateToken 生成用户JWT（内部API，需要auth:token权限）
func (h *AuthHandler) GenerateToken(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 从中间件获取租户信息
	tenantInterface, exists := c.Get("tenant")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}
	tenant := tenantInterface.(*database.Tenant)
	// 这里只做简单演示，实际应校验用户、权限等
	claims := auth.Claims{
		UserID:   req.UserID,
		TenantID: tenant.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := h.signer.Sign(&claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ValidateToken 校验JWT（内部API，需auth:token权限）
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 这里只做简单演示，实际应校验token有效性
	claims := &auth.Claims{}
	err := h.signer.Parse(req.Token, claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"claims": claims})
}

// RefreshToken 刷新access_token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 解析refresh_token，查库校验
	userID, err := h.userService.VerifyAndRefreshToken(c.Request.Context(), req.RefreshToken, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// 生成新access_token和refresh_token
	resp, err := h.userService.IssueNewTokens(c.Request.Context(), userID, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
