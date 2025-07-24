package internal_service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"database/sql"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sqlc-dev/pqtype"
	"golang.org/x/crypto/bcrypt"

	"yuyu-test/internal/auth"
	"yuyu-test/internal/store/database"
)

// Service 内部服务管理服务
type Service struct {
	store           Store
	signer          auth.JWTSigner
	logger          *slog.Logger
	tokenExpiration time.Duration
}

// Store 数据存储接口
type Store interface {
	CreateInternalClient(ctx context.Context, arg database.CreateInternalClientParams) (database.InternalClient, error)
	GetInternalClient(ctx context.Context, clientID string) (database.InternalClient, error)
	ListInternalClients(ctx context.Context) ([]database.InternalClient, error)
	UpdateInternalClient(ctx context.Context, arg database.UpdateInternalClientParams) (database.InternalClient, error)
	DeactivateInternalClient(ctx context.Context, clientID string) error

	ActivateInternalClient(ctx context.Context, clientID string) error
	DeleteInternalClient(ctx context.Context, clientID string) error
	GetClientScopes(ctx context.Context, clientID string) ([]database.GetClientScopesRow, error)
	GrantScopeToClient(ctx context.Context, arg database.GrantScopeToClientParams) error
	RevokeScopeFromClient(ctx context.Context, arg database.RevokeScopeFromClientParams) error
	CheckClientHasScope(ctx context.Context, arg database.CheckClientHasScopeParams) (bool, error)
	ListAllScopes(ctx context.Context) ([]database.Scope, error)
	GetScopeByName(ctx context.Context, scopeName string) (database.Scope, error)
	CreateScope(ctx context.Context, arg database.CreateScopeParams) (database.Scope, error)
	UpdateScope(ctx context.Context, arg database.UpdateScopeParams) (database.Scope, error)
	DeactivateScope(ctx context.Context, scopeName string) error
	LogServiceAccess(ctx context.Context, arg database.LogServiceAccessParams) error
	GetServiceAccessLogs(ctx context.Context, arg database.GetServiceAccessLogsParams) ([]database.ServiceAccessLog, error)
	StoreServiceToken(ctx context.Context, arg database.StoreServiceTokenParams) error
	GetServiceToken(ctx context.Context, tokenHash string) (database.ServiceToken, error)
	RevokeServiceToken(ctx context.Context, tokenHash string) error
	CleanupExpiredTokens(ctx context.Context) error
	GetClientStatistics(ctx context.Context, arg database.GetClientStatisticsParams) (database.GetClientStatisticsRow, error)
}

// NewService 创建内部服务管理服务实例
func NewService(store Store, signer auth.JWTSigner, logger *slog.Logger, tokenExpiration time.Duration) *Service {
	return &Service{
		store:           store,
		signer:          signer,
		logger:          logger,
		tokenExpiration: tokenExpiration,
	}
}

// RegisterServiceRequest 服务注册请求
type RegisterServiceRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Description string `json:"description"`
}

// RegisterServiceResponse 服务注册响应
type RegisterServiceResponse struct {
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	ServiceName  string    `json:"service_name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	Message      string    `json:"message"`
}

// RegisterService 注册新的内部服务
func (s *Service) RegisterService(ctx context.Context, req RegisterServiceRequest) (*RegisterServiceResponse, error) {
	// 生成 client_id 和 client_secret
	clientID := generateRandomID()
	clientSecret, err := s.GenerateClientSecret()
	if err != nil {
		s.logger.Error("failed to generate client secret", "error", err)
		return nil, fmt.Errorf("failed to generate client secret: %w", err)
	}
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash client secret", "error", err)
		return nil, fmt.Errorf("failed to hash client secret: %w", err)
	}
	// 创建内部客户端
	client, err := s.store.CreateInternalClient(ctx, database.CreateInternalClientParams{
		ClientID:         clientID,
		ClientSecretHash: string(hashedSecret),
		ServiceName:      req.ServiceName,
		Description:      sql.NullString{String: req.Description, Valid: req.Description != ""},
	})
	if err != nil {
		s.logger.Error("failed to create internal client", "error", err, "client_id", clientID)
		return nil, fmt.Errorf("failed to create internal client: %w", err)
	}
	s.logger.Info("internal service registered", "client_id", clientID, "service_name", req.ServiceName)
	// 返回响应（包含明文 client_secret，仅此一次）
	return &RegisterServiceResponse{
		ClientID:     client.ClientID,
		ClientSecret: clientSecret, // 新增返回
		ServiceName:  client.ServiceName,
		Description:  client.Description.String,
		CreatedAt:    client.CreatedAt,
		Message:      "Service registered successfully",
	}, nil
}

// AuthenticateServiceRequest 服务认证请求
type AuthenticateServiceRequest struct {
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
}

// AuthenticateServiceResponse 服务认证响应
type AuthenticateServiceResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int64    `json:"expires_in"`
	Scopes      []string `json:"scopes"`
}

// AuthenticateService 认证内部服务并颁发JWT令牌
func (s *Service) AuthenticateService(ctx context.Context, req AuthenticateServiceRequest) (*AuthenticateServiceResponse, error) {
	// 获取内部客户端
	client, err := s.store.GetInternalClient(ctx, req.ClientID)
	if err != nil {
		s.logger.Error("failed to get internal client", "error", err, "client_id", req.ClientID)
		return nil, fmt.Errorf("invalid client credentials")
	}

	// 验证客户端密钥
	err = bcrypt.CompareHashAndPassword([]byte(client.ClientSecretHash), []byte(req.ClientSecret))
	if err != nil {
		s.logger.Error("invalid client secret", "error", err, "client_id", req.ClientID)
		return nil, fmt.Errorf("invalid client credentials")
	}

	// 获取客户端权限
	scopes, err := s.store.GetClientScopes(ctx, req.ClientID)
	if err != nil {
		s.logger.Error("failed to get client scopes", "error", err, "client_id", req.ClientID)
		return nil, fmt.Errorf("failed to get client scopes: %w", err)
	}

	// 构建权限列表
	scopeNames := make([]string, len(scopes))
	for i, scope := range scopes {
		scopeNames[i] = scope.ScopeName
	}

	// 生成JWT令牌
	claims := jwt.MapClaims{
		"sub":    client.ClientID,
		"iss":    "idaas-internal",
		"aud":    "internal-services",
		"scopes": scopeNames,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(s.tokenExpiration).Unix(),
	}
	tokenString, err := s.signer.Sign(claims)
	if err != nil {
		s.logger.Error("failed to sign JWT token", "error", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 存储令牌信息
	tokenHash := s.hashToken(tokenString)
	err = s.store.StoreServiceToken(ctx, database.StoreServiceTokenParams{
		ClientID:  client.ClientID,
		TokenHash: tokenHash,
		Scopes:    scopeNames,
		ExpiresAt: time.Now().Add(s.tokenExpiration),
	})
	if err != nil {
		s.logger.Error("failed to store service token", "error", err)
		// 不返回错误，因为令牌已经生成
	}

	s.logger.Info("service authenticated", "client_id", req.ClientID, "scopes_count", len(scopeNames))

	return &AuthenticateServiceResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int64(s.tokenExpiration.Seconds()),
		Scopes:      scopeNames,
	}, nil
}

// ValidateTokenRequest 令牌验证请求
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// ValidateTokenResponse 令牌验证响应
type ValidateTokenResponse struct {
	Valid    bool     `json:"valid"`
	ClientID string   `json:"client_id,omitempty"`
	Scopes   []string `json:"scopes,omitempty"`
	Message  string   `json:"message,omitempty"`
}

// ValidateToken 验证JWT令牌
func (s *Service) ValidateToken(ctx context.Context, req ValidateTokenRequest) (*ValidateTokenResponse, error) {
	claims, err := parseServiceJWTToken(req.Token, s.signer.Algorithm(), s.signer.PublicKey())
	if err != nil {
		return &ValidateTokenResponse{
			Valid:   false,
			Message: err.Error(),
		}, nil
	}
	// 解析claims
	clientID, _ := (*claims)["sub"].(string)
	scopes := []string{}
	if arr, ok := (*claims)["scopes"].([]interface{}); ok {
		for _, v := range arr {
			if s, ok := v.(string); ok {
				scopes = append(scopes, s)
			}
		}
	}
	// 检查令牌是否被撤销
	tokenHash := s.hashToken(req.Token)
	_, err = s.store.GetServiceToken(ctx, tokenHash)
	if err != nil {
		return &ValidateTokenResponse{
			Valid:   false,
			Message: "Token has been revoked",
		}, nil
	}

	return &ValidateTokenResponse{
		Valid:    true,
		ClientID: clientID,
		Scopes:   scopes,
		Message:  "Token is valid",
	}, nil
}

// GrantScopeRequest 授权权限请求
type GrantScopeRequest struct {
	ClientID  string `json:"client_id" binding:"required"`
	ScopeName string `json:"scope_name" binding:"required"`
	GrantedBy string `json:"granted_by" binding:"required"`
}

// GrantScopeResponse 授权权限响应
type GrantScopeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GrantScope 为内部服务授权权限
func (s *Service) GrantScope(ctx context.Context, req GrantScopeRequest) (*GrantScopeResponse, error) {
	// 获取权限
	scope, err := s.store.GetScopeByName(ctx, req.ScopeName)
	if err != nil {
		return nil, fmt.Errorf("scope not found: %s", req.ScopeName)
	}

	// 授权权限
	err = s.store.GrantScopeToClient(ctx, database.GrantScopeToClientParams{
		ClientID:  req.ClientID,
		ScopeID:   scope.ID,
		GrantedBy: sql.NullString{String: req.GrantedBy, Valid: req.GrantedBy != ""},
	})
	if err != nil {
		s.logger.Error("failed to grant scope", "error", err, "client_id", req.ClientID, "scope", req.ScopeName)
		return nil, fmt.Errorf("failed to grant scope: %w", err)
	}

	s.logger.Info("scope granted", "client_id", req.ClientID, "scope", req.ScopeName, "granted_by", req.GrantedBy)

	return &GrantScopeResponse{
		Success: true,
		Message: fmt.Sprintf("Scope '%s' granted to client '%s'", req.ScopeName, req.ClientID),
	}, nil
}

// RevokeScopeRequest 撤销权限请求
type RevokeScopeRequest struct {
	ClientID  string `json:"client_id" binding:"required"`
	ScopeName string `json:"scope_name" binding:"required"`
}

// RevokeScopeResponse 撤销权限响应
type RevokeScopeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RevokeScope 撤销内部服务的权限
func (s *Service) RevokeScope(ctx context.Context, req RevokeScopeRequest) (*RevokeScopeResponse, error) {
	// 获取权限
	scope, err := s.store.GetScopeByName(ctx, req.ScopeName)
	if err != nil {
		return nil, fmt.Errorf("scope not found: %s", req.ScopeName)
	}

	// 撤销权限
	err = s.store.RevokeScopeFromClient(ctx, database.RevokeScopeFromClientParams{
		ClientID: req.ClientID,
		ScopeID:  scope.ID,
	})
	if err != nil {
		s.logger.Error("failed to revoke scope", "error", err, "client_id", req.ClientID, "scope", req.ScopeName)
		return nil, fmt.Errorf("failed to revoke scope: %w", err)
	}

	s.logger.Info("scope revoked", "client_id", req.ClientID, "scope", req.ScopeName)

	return &RevokeScopeResponse{
		Success: true,
		Message: fmt.Sprintf("Scope '%s' revoked from client '%s'", req.ScopeName, req.ClientID),
	}, nil
}

// CheckPermissionRequest 权限检查请求
type CheckPermissionRequest struct {
	ClientID  string `json:"client_id" binding:"required"`
	ScopeName string `json:"scope_name" binding:"required"`
}

// CheckPermissionResponse 权限检查响应
type CheckPermissionResponse struct {
	HasPermission bool   `json:"has_permission"`
	Message       string `json:"message"`
}

// CheckPermission 检查内部服务是否有指定权限
func (s *Service) CheckPermission(ctx context.Context, req CheckPermissionRequest) (*CheckPermissionResponse, error) {
	hasScope, err := s.store.CheckClientHasScope(ctx, database.CheckClientHasScopeParams{
		ClientID:  req.ClientID,
		ScopeName: req.ScopeName,
	})
	if err != nil {
		s.logger.Error("failed to check client scope", "error", err, "client_id", req.ClientID, "scope", req.ScopeName)
		return nil, fmt.Errorf("failed to check permission: %w", err)
	}

	message := "Permission denied"
	if hasScope {
		message = "Permission granted"
	}

	return &CheckPermissionResponse{
		HasPermission: hasScope,
		Message:       message,
	}, nil
}

// ListServicesResponse 服务列表响应
type ListServicesResponse struct {
	Services []ServiceInfo `json:"services"`
	Total    int           `json:"total"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	ClientID    string    `json:"client_id"`
	ServiceName string    `json:"service_name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Scopes      []string  `json:"scopes"`
}

// ListServices 列出所有内部服务
func (s *Service) ListServices(ctx context.Context) (*ListServicesResponse, error) {
	clients, err := s.store.ListInternalClients(ctx)
	if err != nil {
		s.logger.Error("failed to list internal clients", "error", err)
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	services := make([]ServiceInfo, len(clients))
	for i, client := range clients {
		// 获取服务权限
		scopes, err := s.store.GetClientScopes(ctx, client.ClientID)
		if err != nil {
			s.logger.Error("failed to get client scopes", "error", err, "client_id", client.ClientID)
			continue
		}

		scopeNames := make([]string, len(scopes))
		for j, scope := range scopes {
			scopeNames[j] = scope.ScopeName
		}

		services[i] = ServiceInfo{
			ClientID:    client.ClientID,
			ServiceName: client.ServiceName,
			Description: client.Description.String,
			IsActive:    client.IsActive.Bool,
			CreatedAt:   client.CreatedAt,
			UpdatedAt:   client.UpdatedAt,
			Scopes:      scopeNames,
		}
	}

	return &ListServicesResponse{
		Services: services,
		Total:    len(services),
	}, nil
}

// GenerateClientSecret 生成客户端密钥
func (s *Service) GenerateClientSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// hashToken 哈希令牌用于存储
func (s *Service) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// LogAccess 记录服务访问日志
func (s *Service) LogAccess(ctx context.Context, clientID, endpoint, method string, statusCode int, responseTimeMs int, ipAddress, userAgent, requestBody, responseBody string) error {
	return s.store.LogServiceAccess(ctx, database.LogServiceAccessParams{
		ClientID:       clientID,
		Endpoint:       endpoint,
		Method:         method,
		StatusCode:     int32(statusCode),
		ResponseTimeMs: sql.NullInt32{Int32: int32(responseTimeMs), Valid: true},
		IpAddress:      pqtype.Inet{}, // 需要根据实际情况填充
		UserAgent:      sql.NullString{String: userAgent, Valid: userAgent != ""},
		RequestBody:    sql.NullString{String: requestBody, Valid: requestBody != ""},
		ResponseBody:   sql.NullString{String: responseBody, Valid: responseBody != ""},
	})
}

// GetAccessLogs 获取访问日志
func (s *Service) GetAccessLogs(ctx context.Context, clientID string, limit, offset int32) ([]database.ServiceAccessLog, error) {
	return s.store.GetServiceAccessLogs(ctx, database.GetServiceAccessLogsParams{
		ClientID: clientID,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetStatistics 获取服务统计信息
func (s *Service) GetStatistics(ctx context.Context, clientID string, since time.Time) (*database.GetClientStatisticsRow, error) {
	stats, err := s.store.GetClientStatistics(ctx, database.GetClientStatisticsParams{
		ClientID:  clientID,
		CreatedAt: since,
	})
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// CleanupExpiredTokens 清理过期令牌
func (s *Service) CleanupExpiredTokens(ctx context.Context) error {
	return s.store.CleanupExpiredTokens(ctx)
}

// 从Authorization头中提取客户端ID
func (s *Service) ExtractClientIDFromToken(tokenString string) (string, error) {
	// 移除Bearer前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 解析JWT令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.signer.PublicKey(), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	clientID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("missing client ID in token")
	}

	return clientID, nil
}

// 工具函数
func getSigningMethod(alg string) jwt.SigningMethod {
	if alg == "RS256" {
		return jwt.SigningMethodRS256
	}
	return jwt.SigningMethodHS256
}

func parseServiceJWTToken(tokenString string, alg string, publicKey interface{}) (*jwt.MapClaims, error) {
	var keyFunc jwt.Keyfunc
	if alg == "RS256" {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		}
	} else if alg == "ES256" {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		}
	} else {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		}
	}
	parsed, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return &claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// generateRandomID 生成唯一 client_id
func generateRandomID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
