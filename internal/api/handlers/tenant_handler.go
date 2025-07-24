package handlers

import (
	"net/http"

	"yuyu-test/internal/tenant"

	"github.com/gin-gonic/gin"
)

// TenantHandler 租户处理器
type TenantHandler struct {
	tenantService *tenant.Service
}

// NewTenantHandler 创建新的租户处理器
func NewTenantHandler(tenantService *tenant.Service) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// CreateTenant 创建租户
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req tenant.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.tenantService.CreateTenant(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetTenant 获取租户信息
func (h *TenantHandler) GetTenant(c *gin.Context) {
	tenantID := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant ID is required"})
		return
	}

	tenant, err := h.tenantService.GetTenantByID(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// GetTenants 获取所有租户（内部API，需tenant:read权限）
func (h *TenantHandler) GetTenants(c *gin.Context) {
	tenants, err := h.tenantService.ListTenants(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tenants": tenants})
}
