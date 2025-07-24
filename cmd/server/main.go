package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/sql"
	"yuyu-test/internal/api"
	"yuyu-test/internal/api/handlers"
	"yuyu-test/internal/api/middleware"
	"yuyu-test/internal/auth"
	"yuyu-test/internal/common"
	"yuyu-test/internal/config"
	"yuyu-test/internal/internal_service"
	"yuyu-test/internal/store/database"
	"yuyu-test/internal/tenant"
	"yuyu-test/internal/user"
)

func main() {
	// 初始化日志
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 处理用户JWT密钥
	var userPrivateKey, userPublicKey string
	if cfg.JWTAlgorithm == "HS256" {
		userPrivateKey = cfg.JWTUserSecret
		userPublicKey = cfg.JWTUserSecret
	} else if cfg.JWTAlgorithm == "RS256" || cfg.JWTAlgorithm == "ES256" {
		userPrivateKey, err = config.LoadPEMKey(cfg.JWTUserPrivateKey)
		if err != nil {
			slog.Error("Failed to load user private key", "error", err)
			os.Exit(1)
		}
		userPublicKey, err = config.LoadPEMKey(cfg.JWTUserPublicKey)
		if err != nil {
			slog.Error("Failed to load user public key", "error", err)
			os.Exit(1)
		}
	} else {
		slog.Error("Unsupported JWT_ALGORITHM", "alg", cfg.JWTAlgorithm)
		os.Exit(1)
	}

	// 处理服务JWT密钥
	var servicePrivateKey, servicePublicKey string
	if cfg.JWTAlgorithm == "HS256" {
		servicePrivateKey = cfg.JWTServiceSecret
		servicePublicKey = cfg.JWTServiceSecret
	} else if cfg.JWTAlgorithm == "RS256" || cfg.JWTAlgorithm == "ES256" {
		servicePrivateKey, err = config.LoadPEMKey(cfg.JWTServicePrivateKey)
		if err != nil {
			slog.Error("Failed to load service private key", "error", err)
			os.Exit(1)
		}
		servicePublicKey, err = config.LoadPEMKey(cfg.JWTServicePublicKey)
		if err != nil {
			slog.Error("Failed to load service public key", "error", err)
			os.Exit(1)
		}
	} else {
		slog.Error("Unsupported JWT_ALGORITHM", "alg", cfg.JWTAlgorithm)
		os.Exit(1)
	}

	// 实例化用户JWTSigner
	var userSigner auth.JWTSigner
	if cfg.JWTAlgorithm == "HS256" {
		userSigner = auth.NewHS256Signer(cfg.JWTUserSecret)
	} else if cfg.JWTAlgorithm == "RS256" {
		priv, err := common.ParseRSAPrivateKeyFromPEM([]byte(userPrivateKey))
		if err != nil {
			slog.Error("Failed to parse user RSA private key", "error", err)
			os.Exit(1)
		}
		pub, err := common.ParseRSAPublicKeyFromPEM([]byte(userPublicKey))
		if err != nil {
			slog.Error("Failed to parse user RSA public key", "error", err)
			os.Exit(1)
		}
		userSigner = auth.NewRS256Signer(priv, pub)
	} else if cfg.JWTAlgorithm == "ES256" {
		priv, err := common.ParseECPrivateKeyFromPEM([]byte(userPrivateKey))
		if err != nil {
			slog.Error("Failed to parse user EC private key", "error", err)
			os.Exit(1)
		}
		pub, err := common.ParseECPublicKeyFromPEM([]byte(userPublicKey))
		if err != nil {
			slog.Error("Failed to parse user EC public key", "error", err)
			os.Exit(1)
		}
		userSigner = auth.NewES256Signer(priv, pub)
	} else {
		slog.Error("Unsupported JWT_ALGORITHM", "alg", cfg.JWTAlgorithm)
		os.Exit(1)
	}

	// 实例化服务JWTSigner
	var internalServiceSigner auth.JWTSigner
	if cfg.JWTAlgorithm == "HS256" {
		internalServiceSigner = auth.NewHS256Signer(cfg.JWTServiceSecret)
	} else if cfg.JWTAlgorithm == "RS256" {
		priv, err := common.ParseRSAPrivateKeyFromPEM([]byte(servicePrivateKey))
		if err != nil {
			slog.Error("Failed to parse service RSA private key", "error", err)
			os.Exit(1)
		}
		pub, err := common.ParseRSAPublicKeyFromPEM([]byte(servicePublicKey))
		if err != nil {
			slog.Error("Failed to parse service RSA public key", "error", err)
			os.Exit(1)
		}
		internalServiceSigner = auth.NewRS256Signer(priv, pub)
	} else if cfg.JWTAlgorithm == "ES256" {
		priv, err := common.ParseECPrivateKeyFromPEM([]byte(servicePrivateKey))
		if err != nil {
			slog.Error("Failed to parse service EC private key", "error", err)
			os.Exit(1)
		}
		pub, err := common.ParseECPublicKeyFromPEM([]byte(servicePublicKey))
		if err != nil {
			slog.Error("Failed to parse service EC public key", "error", err)
			os.Exit(1)
		}
		internalServiceSigner = auth.NewES256Signer(priv, pub)
	} else {
		slog.Error("Unsupported JWT_ALGORITHM", "alg", cfg.JWTAlgorithm)
		os.Exit(1)
	}

	// 初始化数据库连接
	sqlDB, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()
	queries := database.New(sqlDB)

	// 初始化服务
	tenantService := tenant.NewService(queries)
	userService := user.NewService(queries, userSigner)

	// 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(tenantService, userSigner)

	// 初始化对内服务管理服务和相关组件
	internalService := internal_service.NewService(queries, internalServiceSigner, logger, time.Duration(cfg.ServiceTokenExpiration)*time.Second)
	internalServiceHandler := handlers.NewInternalServiceHandler(internalService, logger)
	internalAuthMiddleware := middleware.NewInternalAuthMiddleware(internalService, logger)

	// 初始化认证处理器，传递多算法参数
	authHandler := handlers.NewAuthHandler(userService, userSigner)
	internalAuthHandler := handlers.NewInternalAuthHandler(queries, internalServiceSigner)

	// 初始化路由
	router := api.NewRouter(
		tenantService,
		userService,
		authMiddleware,
		authHandler,
		internalAuthHandler,
		internalServiceHandler,
		internalAuthMiddleware,
		sqlDB,
	)
	httpServer := router.Setup()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: httpServer,
	}

	// 启动服务器
	go func() {
		slog.Info("Starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}
