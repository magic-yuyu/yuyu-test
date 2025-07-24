package tenant

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"

	"yuyu-test/internal/auth"
	"yuyu-test/internal/store/database"
)

// Service 租户服务
type Service struct {
	db database.Querier
}

// NewService 创建新的租户服务
func NewService(db database.Querier) *Service {
	return &Service{db: db}
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateTenantResponse 创建租户响应
type CreateTenantResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
	CreatedAt string `json:"created_at"`
}

// CreateTenant 创建新租户
func (s *Service) CreateTenant(ctx context.Context, req CreateTenantRequest) (*CreateTenantResponse, error) {
	// 生成租户ID
	tenantID := generateID("tnt")

	// 生成API密钥
	publicKey := generateAPIKey()
	secretKey := generateAPIKey()

	// 哈希密钥
	secretKeyHash, err := auth.HashPassword(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to hash secret key: %w", err)
	}

	// 创建租户
	tenant, err := s.db.CreateTenant(ctx, database.CreateTenantParams{
		ID:               tenantID,
		Name:             req.Name,
		ApiSecretKeyHash: secretKeyHash,
		ApiPublicKey:     publicKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	slog.Info("Tenant created", "tenant_id", tenantID, "name", req.Name)

	return &CreateTenantResponse{
		ID:        tenant.ID,
		Name:      tenant.Name,
		PublicKey: publicKey,
		SecretKey: secretKey,
		CreatedAt: tenant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetTenantByID 根据ID获取租户
func (s *Service) GetTenantByID(ctx context.Context, tenantID string) (*database.Tenant, error) {
	tenant, err := s.db.GetTenantByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	return &tenant, nil
}

// ValidateAPIKey 验证API密钥
func (s *Service) ValidateAPIKey(ctx context.Context, apiKey string) (*database.Tenant, error) {
	// 先尝试通过公钥查找
	tenant, err := s.db.GetTenantByPublicKey(ctx, apiKey)
	if err == nil {
		return &tenant, nil
	}

	// 如果公钥不匹配，尝试通过密钥哈希查找
	tenant, err = s.db.GetTenantBySecretKeyHash(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid API key")
	}

	return &tenant, nil
}

// ListTenants 获取所有租户
func (s *Service) ListTenants(ctx context.Context) ([]database.Tenant, error) {
	return s.db.ListTenants(ctx)
}

// generateID 生成唯一ID
func generateID(prefix string) string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return prefix + "_" + hex.EncodeToString(bytes)
}

// generateAPIKey 生成API密钥
func generateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
