package ginUtil

import (
	"context"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lily0749labs/gopkgs/jwtUtil"
	"github.com/lily0749labs/gopkgs/log"
)

var (
	Ctx    = contextUtil{}
	CtxPtr = &Ctx
)

type contextUtil struct{}

func asGinContext(ctx context.Context) (*gin.Context, bool) {
	ginCtx, ok := ctx.(*gin.Context)
	return ginCtx, ok && ginCtx != nil
}

func (contextUtil) ClientIP(ctx context.Context) string {
	ginCtx, ok := asGinContext(ctx)
	if !ok {
		return ""
	}
	ip := ginCtx.ClientIP()
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

func (contextUtil) ClearToken(ctx context.Context) {
	ginCtx, ok := asGinContext(ctx)
	if !ok {
		return
	}
	setTokenCookie(ginCtx, "", -1)
}

func (contextUtil) SetToken(ctx context.Context, token string, maxAge int) {
	ginCtx, ok := asGinContext(ctx)
	if !ok {
		return
	}
	setTokenCookie(ginCtx, token, maxAge)
}

func setTokenCookie(ctx *gin.Context, token string, maxAge int) {
	host, _, err := net.SplitHostPort(ctx.Request.Host)
	if err != nil {
		host = ctx.Request.Host
	}
	if net.ParseIP(host) != nil {
		host = ""
	}
	ctx.SetCookie("x-token", token, maxAge, "/", host, false, false)
}

func (c contextUtil) GetToken(ctx context.Context) string {
	ginCtx, ok := asGinContext(ctx)
	if !ok {
		return ""
	}

	token, _ := ginCtx.Cookie("x-token")
	if token != "" {
		return token
	}

	token = ginCtx.Request.Header.Get("x-token")
	claims, err := jwtUtil.NewJWT().ParseToken(token)
	if err != nil {
		log.Error("重新写入cookie token失败,未能成功解析token,请检查请求头是否存在x-token且claims是否为规定结构", log.Line())
		return token
	}
	c.SetToken(ginCtx, token, int((claims.ExpiresAt.Unix()-time.Now().Unix())/60))
	return token
}

func (c contextUtil) GetClaims(ctx context.Context) (*jwtUtil.CustomClaims, error) {
	ginCtx, ok := asGinContext(ctx)
	if !ok {
		return nil, nil
	}

	claims, err := jwtUtil.NewJWT().ParseToken(c.GetToken(ginCtx))
	if err != nil {
		log.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构", log.Line())
	}
	return claims, err
}

func (c contextUtil) GetMbrId(ctx context.Context) uint64 {
	claims := c.GetMbrInfo(ctx)
	if claims == nil {
		return 0
	}
	return claims.BaseClaims.ID
}

func (c contextUtil) GetMbrInfo(ctx context.Context) *jwtUtil.CustomClaims {
	ginCtx, ok := asGinContext(ctx)
	if !ok {
		return nil
	}

	if claims, exists := ginCtx.Get("claims"); exists {
		result, _ := claims.(*jwtUtil.CustomClaims)
		return result
	}
	claims, err := c.GetClaims(ginCtx)
	if err != nil {
		return nil
	}
	return claims
}

func (c contextUtil) GetMbrName(ctx context.Context) string {
	claims := c.GetMbrInfo(ctx)
	if claims == nil {
		return ""
	}
	return claims.Account
}
