// Package gopkgs 为常用且稳定的子包能力提供统一入口。
//
// 具体实现仍保留在各子包中；根包只做类型别名、常量转发和薄包装，
// 方便业务项目按需使用门面，同时保留直接导入子包的能力。
package gopkgs

import (
	"github.com/gin-gonic/gin"
	"github.com/lily0749labs/gopkgs/apiRsp"
	"github.com/lily0749labs/gopkgs/ginUtil"
)

// API 响应类型。
type (
	ApiRsp  = apiRsp.ApiRsp
	IApiRsp = apiRsp.IApiRsp
)

// API 响应状态码。
const (
	InternalServerError = apiRsp.InternalServerError
	SUCCESS             = apiRsp.SUCCESS
	ERROR               = apiRsp.ERROR
	PARAMS              = apiRsp.PARAMS
)

// API 响应消息。
const (
	Msg_Unknown             = apiRsp.Msg_Unknown
	Msg_OperationSuccessful = apiRsp.Msg_OperationSuccessful
	Msg_OperationFailed     = apiRsp.Msg_OperationFailed
	Msg_Success             = apiRsp.Msg_Success
	Msg_Params              = apiRsp.Msg_Params
)

// Gin Engine 相关类型。
type (
	EngineCallback = ginUtil.EngineCallback
	EngineConfig   = ginUtil.EngineConfig
	CORSConfig     = ginUtil.CORSConfig
)

// NewEngine 使用公共配置创建 Gin Engine。
func NewEngine(config EngineConfig, callbacks ...EngineCallback) *gin.Engine {
	return ginUtil.NewEngine(config, callbacks...)
}

// NewCORS 根据公共配置创建 Gin CORS 中间件。
func NewCORS(config CORSConfig) gin.HandlerFunc {
	return ginUtil.NewCORS(config)
}

// PermissiveCORSConfig 返回开放式 CORS 配置。
func PermissiveCORSConfig() *CORSConfig {
	return ginUtil.PermissiveCORSConfig()
}
