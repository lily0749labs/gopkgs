package ginUtil

import (
	"net/http"
	"strings"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	gincors "github.com/rs/cors/wrapper/gin"
)

// EngineCallback 在公共中间件安装完成后注册业务路由或执行引擎定制。
type EngineCallback func(*gin.Engine)

// EngineConfig 定义 Gin Engine 的通用初始化行为。
type EngineConfig struct {
	Mode                   string
	HandleMethodNotAllowed bool
	MaxMultipartMemory     int64
	Middlewares            []gin.HandlerFunc
	CORS                   *CORSConfig
	EnablePprof            bool
}

// CORSConfig 隔离底层 CORS 实现，业务项目无需直接依赖 CORS 库。
type CORSConfig struct {
	AllowedOrigins       []string
	AllowOriginFunc      func(origin string) bool
	AllowedMethods       []string
	AllowedHeaders       []string
	ExposedHeaders       []string
	MaxAge               int
	AllowCredentials     bool
	AllowPrivateNetwork  bool
	OptionsPassthrough   bool
	OptionsSuccessStatus int
	Debug                bool
}

// PermissiveCORSConfig 返回 hall、domain 服务原有的开放 CORS 配置。
func PermissiveCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
	}
}

// NewCORS 根据公共配置创建 Gin CORS 中间件。
func NewCORS(config CORSConfig) gin.HandlerFunc {
	return gincors.New(gincors.Options{
		AllowedOrigins:       config.AllowedOrigins,
		AllowOriginFunc:      config.AllowOriginFunc,
		AllowedMethods:       config.AllowedMethods,
		AllowedHeaders:       config.AllowedHeaders,
		ExposedHeaders:       config.ExposedHeaders,
		MaxAge:               config.MaxAge,
		AllowCredentials:     config.AllowCredentials,
		AllowPrivateNetwork:  config.AllowPrivateNetwork,
		OptionsPassthrough:   config.OptionsPassthrough,
		OptionsSuccessStatus: config.OptionsSuccessStatus,
		Debug:                config.Debug,
	})
}

// NewEngine 创建 Gin Engine，并通过回调把业务路由留在业务项目中。
func NewEngine(config EngineConfig, callbacks ...EngineCallback) *gin.Engine {
	setMode(config.Mode)

	engine := gin.Default()
	engine.HandleMethodNotAllowed = config.HandleMethodNotAllowed
	if config.MaxMultipartMemory > 0 {
		engine.MaxMultipartMemory = config.MaxMultipartMemory
	}
	if len(config.Middlewares) > 0 {
		engine.Use(config.Middlewares...)
	}
	if config.CORS != nil {
		engine.Use(NewCORS(*config.CORS))
	}
	for _, callback := range callbacks {
		if callback != nil {
			callback(engine)
		}
	}
	if config.EnablePprof {
		ginpprof.Wrap(engine)
	}
	return engine
}

func setMode(mode string) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "test":
		gin.SetMode(gin.TestMode)
	case "prod", "production", "release":
		gin.SetMode(gin.ReleaseMode)
	case "", "dev", "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}
