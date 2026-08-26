package ginUtil

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lily0749labs/gopkgs/jwtUtil"
)

func TestContextRejectsNonGinContext(t *testing.T) {
	if got := Ctx.ClientIP(context.Background()); got != "" {
		t.Fatalf("ClientIP() = %q", got)
	}
	if got := Ctx.GetMbrId(context.Background()); got != 0 {
		t.Fatalf("GetMbrId() = %d", got)
	}
}

func TestSetTokenCookieDomain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "http://example.com/test", nil)
	ctx.Request.Host = "example.com:8080"
	Ctx.SetToken(ctx, "token-value", 60)
	cookie := recorder.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, "x-token=token-value") || !strings.Contains(cookie, "Domain=example.com") {
		t.Fatalf("Set-Cookie = %q", cookie)
	}
}

func TestClaimsFromHeader(t *testing.T) {
	jwtUtil.SigningKey = "test-signing-key"
	jwtUtil.ExpiresTime = "1h"
	jwtUtil.BufferTime = "10m"
	jwtUtil.Issuer = "test"
	token, _, err := (jwtUtil.Token{}).LoginToken(&jwtUtil.Login{ID: 99, Account: "tester"})
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest("GET", "http://127.0.0.1/test", nil)
	ctx.Request.Header.Set("x-token", token)

	if got := Ctx.GetMbrId(ctx); got != 99 {
		t.Fatalf("GetMbrId() = %d", got)
	}
	if got := Ctx.GetMbrName(ctx); got != "tester" {
		t.Fatalf("GetMbrName() = %q", got)
	}
}

func TestClaimsFromGinContext(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	want := &jwtUtil.CustomClaims{BaseClaims: jwtUtil.BaseClaims{ID: 100, Account: "context-user"}}
	ctx.Set("claims", want)

	if got := Ctx.GetMbrInfo(ctx); got != want {
		t.Fatalf("GetMbrInfo() = %p, want %p", got, want)
	}
	if got := Ctx.GetMbrId(ctx); got != 100 {
		t.Fatalf("GetMbrId() = %d", got)
	}
	if got := Ctx.GetMbrName(ctx); got != "context-user" {
		t.Fatalf("GetMbrName() = %q", got)
	}
}
