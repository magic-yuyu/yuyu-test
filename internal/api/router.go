package api

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"yuyu-test/internal/api/handlers"
	"yuyu-test/internal/api/middleware"
	"yuyu-test/internal/config"
	"yuyu-test/internal/tenant"
	"yuyu-test/internal/user"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// 仅允许本地和内网访问的中间件 目前没有使用这个函数
func k8sProbeOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		fmt.Println("访问ip:", ip)
		if isAllowedHealthCheckIP(ip) {
			c.Next()
			return
		}
		c.JSON(403, gin.H{"error": "forbidden"})
		c.Abort()
	}
}
func isAllowedHealthCheckIP(ip string) bool {
	if ip == "127.0.0.1" ||
		strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "10.") ||
		(strings.HasPrefix(ip, "172.") && is172Private(ip)) {
		return true
	}
	return false
}

func is172Private(ip string) bool {
	// 172.16.0.0 - 172.31.255.255
	parts := strings.Split(ip, ".")
	if len(parts) < 2 {
		return false
	}
	if parts[0] == "172" {
		second, _ := strconv.Atoi(parts[1])
		return second >= 16 && second <= 31
	}
	return false
}

// Router API路由
type Router struct {
	tenantHandler          *handlers.TenantHandler
	authHandler            *handlers.AuthHandler
	userHandler            *handlers.UserHandler
	authMiddleware         *middleware.AuthMiddleware
	internalAuthHandler    *handlers.InternalAuthHandler
	internalServiceHandler *handlers.InternalServiceHandler
	internalAuthMiddleware *middleware.InternalAuthMiddleware
	sqlDB                  *sql.DB // 新增字段用于数据库健康检查
}

// NewRouter 创建新的路由
func NewRouter(
	tenantService *tenant.Service,
	userService *user.Service,
	authMiddleware *middleware.AuthMiddleware,
	authHandler *handlers.AuthHandler,
	internalAuthHandler *handlers.InternalAuthHandler,
	internalServiceHandler *handlers.InternalServiceHandler,
	internalAuthMiddleware *middleware.InternalAuthMiddleware,
	sqlDB *sql.DB, // 新增参数
) *Router {
	return &Router{
		tenantHandler:          handlers.NewTenantHandler(tenantService),
		authHandler:            authHandler,
		userHandler:            handlers.NewUserHandler(userService),
		authMiddleware:         authMiddleware,
		internalAuthHandler:    internalAuthHandler,
		internalServiceHandler: internalServiceHandler,
		internalAuthMiddleware: internalAuthMiddleware,
		sqlDB:                  sqlDB,
	}
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// 健康检查（允许所有来源访问）
	router.GET("/health", func(c *gin.Context) {
		cfg, _ := config.Load()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		host, _ := os.Hostname()

		// 数据库健康检查
		dbStatus := "ok"
		dbErr := ""
		if r.sqlDB != nil {
			err := r.sqlDB.Ping()
			if err != nil {
				dbStatus = "error"
				dbErr = err.Error()
			}
		} else {
			dbStatus = "not_initialized"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":        "ok",
			"time":          time.Now().Format(time.RFC3339),
			"uptime":        time.Since(startTime).String(),
			"env":           cfg.Environment,
			"port":          cfg.Port,
			"version":       os.Getenv("APP_VERSION"),
			"database":      cfg.DatabaseURL,
			"go_version":    runtime.Version(),
			"pid":           os.Getpid(),
			"hostname":      host,
			"mem_alloc":     m.Alloc,
			"mem_sys":       m.Sys,
			"mem_heap":      m.HeapAlloc,
			"num_goroutine": runtime.NumGoroutine(),
			"db_status":     dbStatus,
			"db_error":      dbErr,
		})
	})

	// 对内服务认证API
	router.POST("/oauth/token", r.internalAuthHandler.Token)

	// API版本控制
	v1 := router.Group("/v1")
	{
		// 租户管理（无需认证）
		tenants := v1.Group("/tenants")
		{
			tenants.POST("", r.tenantHandler.CreateTenant)
			tenants.GET("/:id", r.tenantHandler.GetTenant)
		}

		// 认证相关（需要API密钥认证）
		auth := v1.Group("/auth")
		auth.Use(r.authMiddleware.APIKeyAuth())
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken) // 新增refresh token接口
		}

		// 用户管理（需要JWT认证）
		users := v1.Group("/users")
		users.Use(r.authMiddleware.JWTAuth())
		{
			users.GET("/me", r.userHandler.GetMe)
		}

		// 用户管理（需要API密钥认证）
		adminUsers := v1.Group("/users")
		adminUsers.Use(r.authMiddleware.APIKeyAuth())
		{
			adminUsers.GET("", r.userHandler.GetUsers)
			adminUsers.GET("/:id", r.userHandler.GetUser)
		}

		// 内部服务管理API
		internal := v1.Group("/internal")
		{
			// 服务注册（无需认证）
			services := internal.Group("/services")
			{
				services.POST("/register", r.internalServiceHandler.RegisterService)
				services.POST("/authenticate", r.internalServiceHandler.AuthenticateService)
				services.POST("/validate-token", r.internalServiceHandler.ValidateToken)
			}

			// 需要内部服务认证的API
			authenticated := internal.Group("/services")
			authenticated.Use(r.internalAuthMiddleware.RequireAuth())
			{
				// 服务列表
				authenticated.GET("", r.internalServiceHandler.ListServices)

				// 权限管理
				authenticated.POST("/grant-scope", r.internalServiceHandler.GrantScope)
				authenticated.POST("/revoke-scope", r.internalServiceHandler.RevokeScope)
				authenticated.POST("/check-permission", r.internalServiceHandler.CheckPermission)

				// 访问日志和统计
				authenticated.GET("/:client_id/logs", r.internalServiceHandler.GetServiceAccessLogs)
				authenticated.GET("/:client_id/statistics", r.internalServiceHandler.GetServiceStatistics)

				// 系统维护
				authenticated.POST("/cleanup-tokens", r.internalServiceHandler.CleanupExpiredTokens)
			}
		}
	}

	// 内部服务API示例（需要特定权限）
	internalAPI := router.Group("/api/internal")
	{
		// 用户管理API（需要user:read权限）
		internalUsers := internalAPI.Group("/users")
		internalUsers.Use(r.internalAuthMiddleware.RequireScope("user:read"))
		{
			internalUsers.GET("", r.userHandler.GetUsers)
			internalUsers.GET("/:id", r.userHandler.GetUser)
		}

		// 用户写入API（需要user:write权限）
		internalUserWrite := internalAPI.Group("/users")
		internalUserWrite.Use(r.internalAuthMiddleware.RequireScope("user:write"))
		{
			internalUserWrite.POST("", r.userHandler.CreateUser)
			internalUserWrite.PUT("/:id", r.userHandler.UpdateUser)
		}

		// 租户管理API（需要tenant:read权限）
		internalTenants := internalAPI.Group("/tenants")
		internalTenants.Use(r.internalAuthMiddleware.RequireScope("tenant:read"))
		{
			internalTenants.GET("", r.tenantHandler.GetTenants)
			internalTenants.GET("/:id", r.tenantHandler.GetTenant)
		}

		// 认证API（需要auth:token权限）
		internalAuth := internalAPI.Group("/auth")
		internalAuth.Use(r.internalAuthMiddleware.RequireScope("auth:token"))
		{
			internalAuth.POST("/token", r.authHandler.GenerateToken)
			internalAuth.POST("/validate", r.authHandler.ValidateToken)
		}

		// 管理API（需要internal:admin权限）
		internalAdmin := internalAPI.Group("/admin")
		internalAdmin.Use(r.internalAuthMiddleware.RequireScope("internal:admin"))
		{
			internalAdmin.GET("/services", r.internalServiceHandler.ListServices)
			internalAdmin.POST("/services/grant-scope", r.internalServiceHandler.GrantScope)
			internalAdmin.POST("/services/revoke-scope", r.internalServiceHandler.RevokeScope)
		}

		// 复合权限API示例
		internalComposite := internalAPI.Group("/composite")
		{
			// 需要任意一个权限
			internalComposite.GET("/any",
				r.internalAuthMiddleware.RequireAnyScope("user:read", "tenant:read"),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Access granted with any scope"})
				},
			)

			// 需要所有权限
			internalComposite.GET("/all",
				r.internalAuthMiddleware.RequireAllScopes("user:read", "user:write"),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "Access granted with all scopes"})
				},
			)
		}

		// 可选认证API示例
		internalOptional := internalAPI.Group("/optional")
		internalOptional.Use(r.internalAuthMiddleware.OptionalAuth())
		{
			internalOptional.GET("/public", func(c *gin.Context) {
				if middleware.IsAuthenticated(c) {
					clientID, _ := middleware.GetClientID(c)
					scopes, _ := middleware.GetScopes(c)
					c.JSON(http.StatusOK, gin.H{
						"message":   "Public endpoint with optional auth",
						"client_id": clientID,
						"scopes":    scopes,
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"message": "Public endpoint without auth",
					})
				}
			})
		}
	}

	return router
}
