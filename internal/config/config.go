package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config 应用配置结构
type Config struct {
	DatabaseURL            string
	JWTAlgorithm           string // 新增，HS256/RS256
	JWTUserSecret          string // HS256密钥
	JWTServiceSecret       string // HS256密钥
	JWTUserPrivateKey      string // RS256私钥PEM内容
	JWTUserPublicKey       string // RS256公钥PEM内容
	JWTServicePrivateKey   string // RS256私钥PEM内容
	JWTServicePublicKey    string // RS256公钥PEM内容
	UserTokenExpiration    int    // 单位秒
	ServiceTokenExpiration int    // 单位秒
	Port                   int
	Environment            string
}

// JWTConfigValidator 定义算法校验接口
type JWTConfigValidator func(cfg *Config) error

// 算法到校验器的映射
var jwtValidators = map[string]JWTConfigValidator{
	"HS256": validateHS256,
	"RS256": validateRS256,
	"ES256": validateES256,
}

// 校验HS256配置
func validateHS256(cfg *Config) error {
	env := cfg.Environment
	if cfg.JWTUserSecret == "" {
		return fmt.Errorf("JWT_USER_SECRET_KEY or JWT_SECRET environment variable is required for user JWT")
	}
	if cfg.JWTServiceSecret == "" {
		return fmt.Errorf("JWT_SERVICE_SECRET_KEY or JWT_SECRET environment variable is required for service JWT")
	}
	if env == "production" {
		if len(cfg.JWTUserSecret) < 32 {
			return fmt.Errorf("Production JWT_USER_SECRET_KEY must be at least 32 characters")
		}
		if len(cfg.JWTServiceSecret) < 32 {
			return fmt.Errorf("Production JWT_SERVICE_SECRET_KEY must be at least 32 characters")
		}
	} else {
		if len(cfg.JWTUserSecret) < 32 {
			fmt.Println("[WARN] JWT_USER_SECRET_KEY is less than 32 characters. This is insecure and should only be used for local development.")
		}
		if len(cfg.JWTServiceSecret) < 32 {
			fmt.Println("[WARN] JWT_SERVICE_SECRET_KEY is less than 32 characters. This is insecure and should only be used for local development.")
		}
	}
	return nil
}

// 校验RS256配置
func validateRS256(cfg *Config) error {
	if cfg.JWTUserPrivateKey == "" || cfg.JWTUserPublicKey == "" {
		return fmt.Errorf("JWT_USER_PRIVATE_KEY and JWT_USER_PUBLIC_KEY are required for RS256 user JWT")
	}
	if cfg.JWTServicePrivateKey == "" || cfg.JWTServicePublicKey == "" {
		return fmt.Errorf("JWT_SERVICE_PRIVATE_KEY and JWT_SERVICE_PUBLIC_KEY are required for RS256 service JWT")
	}
	return nil
}

// 校验ES256配置
func validateES256(cfg *Config) error {
	if cfg.JWTUserPrivateKey == "" || cfg.JWTUserPublicKey == "" {
		return fmt.Errorf("JWT_USER_PRIVATE_KEY and JWT_USER_PUBLIC_KEY are required for ES256 user JWT")
	}
	if cfg.JWTServicePrivateKey == "" || cfg.JWTServicePublicKey == "" {
		return fmt.Errorf("JWT_SERVICE_PRIVATE_KEY and JWT_SERVICE_PUBLIC_KEY are required for ES256 service JWT")
	}
	return nil
}

// LoadPEMKey 支持从内容或文件路径加载PEM密钥
func LoadPEMKey(val string) (string, error) {
	if val == "" {
		return "", fmt.Errorf("empty key")
	}
	if strings.HasPrefix(val, "-----BEGIN ") {
		return val, nil // 直接是PEM内容
	}
	// 否则认为是文件路径
	data, err := os.ReadFile(val)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}
	return string(data), nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		// 去除首尾空格
		return strings.TrimSpace(value)
	}
	return defaultValue
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	userTokenExp, _ := strconv.Atoi(getEnv("USER_TOKEN_EXPIRATION", "3600"))      // 默认1小时
	serviceTokenExp, _ := strconv.Atoi(getEnv("SERVICE_TOKEN_EXPIRATION", "300")) // 默认5分钟

	algorithm := strings.ToUpper(getEnv("JWT_ALGORITHM", "HS256"))

	userSecret := getEnv("JWT_USER_SECRET_KEY", "")
	if userSecret == "" {
		return nil, fmt.Errorf("JWT_USER_SECRET_KEY environment variable is required for user JWT")
	}

	serviceSecret := getEnv("JWT_SERVICE_SECRET_KEY", "")
	if serviceSecret == "" {
		return nil, fmt.Errorf("JWT_SERVICE_SECRET_KEY environment variable is required for service JWT")
	}

	userPrivKey := getEnv("JWT_USER_PRIVATE_KEY", "")
	userPubKey := getEnv("JWT_USER_PUBLIC_KEY", "")
	servicePrivKey := getEnv("JWT_SERVICE_PRIVATE_KEY", "")
	servicePubKey := getEnv("JWT_SERVICE_PUBLIC_KEY", "")

	config := &Config{
		DatabaseURL:            getEnv("DATABASE_URL", ""),
		JWTAlgorithm:           algorithm,
		JWTUserSecret:          userSecret,
		JWTServiceSecret:       serviceSecret,
		JWTUserPrivateKey:      userPrivKey,
		JWTUserPublicKey:       userPubKey,
		JWTServicePrivateKey:   servicePrivKey,
		JWTServicePublicKey:    servicePubKey,
		UserTokenExpiration:    userTokenExp,
		ServiceTokenExpiration: serviceTokenExp,
		Port:                   port,
		Environment:            getEnv("GO_ENV", "development"),
	}

	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// 校验算法和密钥
	validator, ok := jwtValidators[config.JWTAlgorithm]
	if !ok {
		return nil, fmt.Errorf("Unsupported JWT_ALGORITHM: %s", config.JWTAlgorithm)
	}
	if err := validator(config); err != nil {
		return nil, err
	}

	return config, nil
}
