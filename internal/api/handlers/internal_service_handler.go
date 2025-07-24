package handlers

import (
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"

	"yuyu-test/internal/internal_service"
)

// InternalServiceHandler 内部服务管理处理器
type InternalServiceHandler struct {
	service *internal_service.Service
	logger  *slog.Logger
}

// NewInternalServiceHandler 创建内部服务处理器
func NewInternalServiceHandler(service *internal_service.Service, logger *slog.Logger) *InternalServiceHandler {
	return &InternalServiceHandler{
		service: service,
		logger:  logger,
	}
}

// RegisterService 注册新的内部服务
// @Summary 注册内部服务
// @Description 注册一个新的内部服务，获取客户端ID和密钥
// @Tags 内部服务管理
// @Accept json
// @Produce json
// @Param request body internal_service.RegisterServiceRequest true "服务注册信息"
// @Success 200 {object} internal_service.RegisterServiceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/register [post]
func (h *InternalServiceHandler) RegisterService(c *gin.Context) {
	var req internal_service.RegisterServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind register service request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// 已由 service 层自动生成 client_secret，无需在 handler 生成

	response, err := h.service.RegisterService(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to register service", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	// 在响应中包含生成的密钥（仅此一次）
	responseWithSecret := map[string]interface{}{
		"client_id":     response.ClientID,
		"client_secret": response.ClientSecret, // 仅返回一次
		"service_name":  response.ServiceName,
		"description":   response.Description,
		"created_at":    response.CreatedAt,
		"message":       response.Message,
		"warning":       "Please save the client_secret securely. It will not be shown again.",
	}

	c.JSON(http.StatusCreated, responseWithSecret)
}

// AuthenticateService 内部服务认证
// @Summary 内部服务认证
// @Description 使用客户端ID和密钥进行认证，获取JWT令牌
// @Tags 内部服务管理
// @Accept json
// @Produce json
// @Param request body internal_service.AuthenticateServiceRequest true "认证信息"
// @Success 200 {object} internal_service.AuthenticateServiceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/authenticate [post]
func (h *InternalServiceHandler) AuthenticateService(c *gin.Context) {
	var req internal_service.AuthenticateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind authenticate service request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.AuthenticateService(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to authenticate service", "error", err)
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Authentication failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ValidateToken 验证JWT令牌
// @Summary 验证JWT令牌
// @Description 验证内部服务JWT令牌的有效性
// @Tags 内部服务管理
// @Accept json
// @Produce json
// @Param request body internal_service.ValidateTokenRequest true "令牌验证信息"
// @Success 200 {object} internal_service.ValidateTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/validate-token [post]
func (h *InternalServiceHandler) ValidateToken(c *gin.Context) {
	var req internal_service.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind validate token request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.ValidateToken(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to validate token", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GrantScope 授权权限
// @Summary 授权权限
// @Description 为内部服务授权指定权限
// @Tags 内部服务管理
// @Accept json
// @Produce json
// @Param request body internal_service.GrantScopeRequest true "权限授权信息"
// @Success 200 {object} internal_service.GrantScopeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/grant-scope [post]
func (h *InternalServiceHandler) GrantScope(c *gin.Context) {
	var req internal_service.GrantScopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind grant scope request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.GrantScope(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to grant scope", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RevokeScope 撤销权限
// @Summary 撤销权限
// @Description 撤销内部服务的指定权限
// @Tags 内部服务管理
// @Accept json
// @Produce json
// @Param request body internal_service.RevokeScopeRequest true "权限撤销信息"
// @Success 200 {object} internal_service.RevokeScopeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/revoke-scope [post]
func (h *InternalServiceHandler) RevokeScope(c *gin.Context) {
	var req internal_service.RevokeScopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind revoke scope request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.RevokeScope(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to revoke scope", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CheckPermission 检查权限
// @Summary 检查权限
// @Description 检查内部服务是否有指定权限
// @Tags 内部服务管理
// @Accept json
// @Produce json
// @Param request body internal_service.CheckPermissionRequest true "权限检查信息"
// @Success 200 {object} internal_service.CheckPermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/check-permission [post]
func (h *InternalServiceHandler) CheckPermission(c *gin.Context) {
	var req internal_service.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind check permission request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.CheckPermission(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to check permission", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListServices 列出所有内部服务
// @Summary 列出内部服务
// @Description 获取所有已注册的内部服务列表
// @Tags 内部服务管理
// @Produce json
// @Success 200 {object} internal_service.ListServicesResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services [get]
func (h *InternalServiceHandler) ListServices(c *gin.Context) {
	response, err := h.service.ListServices(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list services", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetServiceAccessLogs 获取服务访问日志
// @Summary 获取访问日志
// @Description 获取指定服务的访问日志
// @Tags 内部服务管理
// @Produce json
// @Param client_id path string true "客户端ID"
// @Param limit query int false "限制数量" default(50)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {array} ServiceAccessLog
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/{client_id}/logs [get]
func (h *InternalServiceHandler) GetServiceAccessLogs(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad request",
			Message: "Client ID is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad request",
			Message: "Invalid limit parameter",
		})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad request",
			Message: "Invalid offset parameter",
		})
		return
	}

	logs, err := h.service.GetAccessLogs(c.Request.Context(), clientID, int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("failed to get access logs", "error", err, "client_id", clientID)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetServiceStatistics 获取服务统计信息
// @Summary 获取服务统计
// @Description 获取指定服务的访问统计信息
// @Tags 内部服务管理
// @Produce json
// @Param client_id path string true "客户端ID"
// @Param since query string false "统计开始时间" default("24h")
// @Success 200 {object} ServiceStatistics
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/{client_id}/statistics [get]
func (h *InternalServiceHandler) GetServiceStatistics(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad request",
			Message: "Client ID is required",
		})
		return
	}

	sinceStr := c.DefaultQuery("since", "24h")
	var since time.Time

	switch sinceStr {
	case "1h":
		since = time.Now().Add(-1 * time.Hour)
	case "24h":
		since = time.Now().Add(-24 * time.Hour)
	case "7d":
		since = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		since = time.Now().Add(-30 * 24 * time.Hour)
	default:
		// 尝试解析自定义时间
		if parsed, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			since = parsed
		} else {
			since = time.Now().Add(-24 * time.Hour) // 默认24小时
		}
	}

	stats, err := h.service.GetStatistics(c.Request.Context(), clientID, since)
	if err != nil {
		h.logger.Error("failed to get service statistics", "error", err, "client_id", clientID)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"client_id":         clientID,
		"since":             since,
		"total_requests":    stats.TotalRequests,
		"avg_response_time": stats.AvgResponseTime,
		"error_count":       stats.ErrorCount,
		"success_rate":      0.0,
	}

	if stats.TotalRequests > 0 {
		response["success_rate"] = float64(stats.TotalRequests-stats.ErrorCount) / float64(stats.TotalRequests) * 100
	}

	c.JSON(http.StatusOK, response)
}

// CleanupExpiredTokens 清理过期令牌
// @Summary 清理过期令牌
// @Description 清理数据库中已过期的JWT令牌
// @Tags 内部服务管理
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /internal/services/cleanup-tokens [post]
func (h *InternalServiceHandler) CleanupExpiredTokens(c *gin.Context) {
	err := h.service.CleanupExpiredTokens(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to cleanup expired tokens", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Expired tokens cleaned up successfully",
	})
}

// ErrorResponse 通用错误响应结构体
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// ServiceAccessLog 服务访问日志结构
type ServiceAccessLog struct {
	ID             int32     `json:"id"`
	ClientID       string    `json:"client_id"`
	Endpoint       string    `json:"endpoint"`
	Method         string    `json:"method"`
	StatusCode     int32     `json:"status_code"`
	ResponseTimeMs int32     `json:"response_time_ms"`
	IpAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	RequestBody    string    `json:"request_body"`
	ResponseBody   string    `json:"response_body"`
	CreatedAt      time.Time `json:"created_at"`
}

// ServiceStatistics 服务统计信息
type ServiceStatistics struct {
	ClientID        string  `json:"client_id"`
	TotalRequests   int64   `json:"total_requests"`
	AvgResponseTime float64 `json:"avg_response_time"`
	ErrorCount      int64   `json:"error_count"`
	SuccessRate     float64 `json:"success_rate"`
}
