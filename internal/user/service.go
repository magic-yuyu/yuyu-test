package user

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"yuyu-test/internal/auth"
	"yuyu-test/internal/store/database"

	"crypto/sha256"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sqlc-dev/pqtype"
)

// Service 用户服务
type Service struct {
	db     database.Querier
	signer auth.JWTSigner
}

// NewService 创建新的用户服务
func NewService(db database.Querier, signer auth.JWTSigner) *Service {
	return &Service{db: db, signer: signer}
}

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Email    string                 `json:"email" binding:"required,email"`
	Password string                 `json:"password" binding:"required,min=6"`
	Profile  map[string]interface{} `json:"profile"`
}

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	Profile   map[string]interface{} `json:"profile"`
	CreatedAt string                 `json:"created_at"`
}

// Register 用户注册
func (s *Service) Register(ctx context.Context, tenantID string, req RegisterRequest) (*RegisterResponse, error) {
	// 检查用户是否已存在
	existingUser, err := s.db.GetUserByEmail(ctx, database.GetUserByEmailParams{
		TenantID: tenantID,
		Email:    req.Email,
	})
	if err == nil && existingUser.ID != "" {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// 生成用户ID
	userID := generateID("usr")

	// 哈希密码
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// profile序列化
	var profileRaw pqtype.NullRawMessage
	if req.Profile != nil {
		profileBytes, err := json.Marshal(req.Profile)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal profile: %w", err)
		}
		profileRaw = pqtype.NullRawMessage{
			RawMessage: profileBytes,
			Valid:      true,
		}
	} else {
		profileRaw = pqtype.NullRawMessage{Valid: false}
	}

	// 创建用户
	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:             userID,
		TenantID:       tenantID,
		Email:          req.Email,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
		Profile:        profileRaw,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("User registered", "user_id", userID, "email", req.Email, "tenant_id", tenantID)

	// 反序列化profile
	var profile map[string]interface{}
	if user.Profile.Valid && len(user.Profile.RawMessage) > 0 {
		_ = json.Unmarshal(user.Profile.RawMessage, &profile)
	} else {
		profile = make(map[string]interface{})
	}

	return &RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		Profile:   profile,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 用户登录响应
type LoginResponse struct {
	User         *RegisterResponse `json:"user"`
	Token        string            `json:"token"`
	RefreshToken string            `json:"refresh_token"`
}

// generateRefreshToken 生成高强度refresh token
func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, tenantID string, req LoginRequest) (*LoginResponse, error) {
	// 获取用户
	user, err := s.db.GetUserByEmail(ctx, database.GetUserByEmailParams{
		TenantID: tenantID,
		Email:    req.Email,
	})
	if err != nil || !user.HashedPassword.Valid {
		return nil, fmt.Errorf("invalid email or password")
	}

	// 验证密码
	if !auth.CheckPassword(req.Password, user.HashedPassword.String) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// 生成JWT令牌
	claims := auth.Claims{
		UserID:   user.ID,
		TenantID: tenantID,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := s.signer.Sign(&claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 生成refresh_token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshTokenHash := sha256.Sum256([]byte(refreshToken))
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	clientIP := ""                                // 可从ctx或请求中获取
	userAgent := ""                               // 可从ctx或请求中获取
	_ = s.db.DeleteAllRefreshTokens(ctx, user.ID) // 单端策略，先清理
	err = s.db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: fmt.Sprintf("%x", refreshTokenHash[:]),
		ExpiresAt: expiresAt,
		ClientIp:  sql.NullString{String: clientIP, Valid: clientIP != ""},
		UserAgent: sql.NullString{String: userAgent, Valid: userAgent != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	slog.Info("User logged in", "user_id", user.ID, "email", user.Email, "tenant_id", tenantID)

	// 反序列化profile
	var profile map[string]interface{}
	if user.Profile.Valid && len(user.Profile.RawMessage) > 0 {
		_ = json.Unmarshal(user.Profile.RawMessage, &profile)
	} else {
		profile = make(map[string]interface{})
	}

	return &LoginResponse{
		User: &RegisterResponse{
			ID:        user.ID,
			Email:     user.Email,
			Profile:   profile,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *Service) GetUserByID(ctx context.Context, userID string) (*RegisterResponse, error) {
	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	var profile map[string]interface{}
	if user.Profile.Valid && len(user.Profile.RawMessage) > 0 {
		_ = json.Unmarshal(user.Profile.RawMessage, &profile)
	} else {
		profile = make(map[string]interface{})
	}
	return &RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		Profile:   profile,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetUsersByTenant 获取租户下的所有用户
func (s *Service) GetUsersByTenant(ctx context.Context, tenantID string) ([]*RegisterResponse, error) {
	users, err := s.db.GetUsersByTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	var responses []*RegisterResponse
	for _, user := range users {
		var profile map[string]interface{}
		if user.Profile.Valid && len(user.Profile.RawMessage) > 0 {
			_ = json.Unmarshal(user.Profile.RawMessage, &profile)
		} else {
			profile = make(map[string]interface{})
		}
		responses = append(responses, &RegisterResponse{
			ID:        user.ID,
			Email:     user.Email,
			Profile:   profile,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return responses, nil
}

// generateID 生成唯一ID
func generateID(prefix string) string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return prefix + "_" + hex.EncodeToString(bytes)
}

// VerifyAndRefreshToken 校验refresh_token，返回userID
func (s *Service) VerifyAndRefreshToken(ctx context.Context, refreshToken, clientIP, userAgent string) (string, error) {
	hash := sha256.Sum256([]byte(refreshToken))
	// 需要用户ID，实际可通过前端传递或解析token内容
	// 这里假设前端传user_id，或可遍历所有用户（不推荐，建议优化）
	// 这里简化为遍历所有用户
	// 实际生产建议refresh_token中带user_id信息
	users, err := s.db.GetUsersByTenant(ctx, "") // 获取所有用户，实际应优化
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}
	var foundToken *database.UserRefreshToken
	for _, u := range users {
		token, err := s.db.GetRefreshToken(ctx, database.GetRefreshTokenParams{
			UserID:    u.ID,
			TokenHash: fmt.Sprintf("%x", hash[:]),
		})
		if err == nil && token.ExpiresAt.After(time.Now()) {
			foundToken = &token
			break
		}
	}
	if foundToken == nil {
		return "", errors.New("invalid or expired refresh token")
	}
	_ = s.db.DeleteRefreshToken(ctx, database.DeleteRefreshTokenParams{
		UserID:    foundToken.UserID,
		TokenHash: foundToken.TokenHash,
	})
	return foundToken.UserID, nil
}

// IssueNewTokens 为用户签发新access_token和refresh_token
func (s *Service) IssueNewTokens(ctx context.Context, userID, clientIP, userAgent string) (*LoginResponse, error) {
	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	claims := auth.Claims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := s.signer.Sign(&claims)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}
	refreshTokenHash := sha256.Sum256([]byte(refreshToken))
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	err = s.db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: fmt.Sprintf("%x", refreshTokenHash[:]),
		ExpiresAt: expiresAt,
		ClientIp:  sql.NullString{String: clientIP, Valid: clientIP != ""},
		UserAgent: sql.NullString{String: userAgent, Valid: userAgent != ""},
	})
	if err != nil {
		return nil, errors.New("failed to store refresh token")
	}
	var profile map[string]interface{}
	if user.Profile.Valid && len(user.Profile.RawMessage) > 0 {
		_ = json.Unmarshal(user.Profile.RawMessage, &profile)
	} else {
		profile = make(map[string]interface{})
	}
	return &LoginResponse{
		User: &RegisterResponse{
			ID:        user.ID,
			Email:     user.Email,
			Profile:   profile,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
