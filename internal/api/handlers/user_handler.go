package handlers

import (
	"net/http"

	"yuyu-test/internal/store/database"
	"yuyu-test/internal/user"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe 获取当前用户信息
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	response, err := h.userService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUser 获取指定用户信息（需要API密钥认证）
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// 从中间件获取租户信息
	tenantInterface, exists := c.Get("tenant")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	_ = tenantInterface.(*database.Tenant)
	response, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUsers 获取租户下的所有用户
func (h *UserHandler) GetUsers(c *gin.Context) {
	// 从中间件获取租户信息
	tenantInterface, exists := c.Get("tenant")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}

	tenant := tenantInterface.(*database.Tenant)
	users, err := h.userService.GetUsersByTenant(c.Request.Context(), tenant.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// CreateUser 创建用户（内部API，需user:write权限）
func (h *UserHandler) CreateUser(c *gin.Context) {
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

// UpdateUser 更新用户（内部API，需user:write权限）
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}
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
	// 这里只做简单演示，实际应调用 UpdateUser 业务逻辑
	response, err := h.userService.Register(c.Request.Context(), tenant.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
