package gopkgs_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lily0749labs/gopkgs"
)

func TestApiRspFacade(t *testing.T) {
	var response gopkgs.IApiRsp = (gopkgs.ApiRsp{}).OkWithData("ok")
	if response.HttpStatus() != http.StatusOK {
		t.Fatalf("HttpStatus() = %d, want %d", response.HttpStatus(), http.StatusOK)
	}
	if response.Code() != gopkgs.SUCCESS {
		t.Fatalf("Code() = %d, want %d", response.Code(), gopkgs.SUCCESS)
	}
	if response.Data() != "ok" {
		t.Fatalf("Data() = %v, want ok", response.Data())
	}
}

func TestGinEngineFacade(t *testing.T) {
	engine := gopkgs.NewEngine(
		gopkgs.EngineConfig{Mode: gin.TestMode},
		func(engine *gin.Engine) {
			engine.GET("/health", func(context *gin.Context) {
				context.Status(http.StatusNoContent)
			})
		},
	)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
}

func TestCORSFacade(t *testing.T) {
	config := gopkgs.PermissiveCORSConfig()
	if config == nil || len(config.AllowedOrigins) != 1 || config.AllowedOrigins[0] != "*" {
		t.Fatalf("PermissiveCORSConfig() = %#v", config)
	}
	if gopkgs.NewCORS(*config) == nil {
		t.Fatal("NewCORS() returned nil")
	}
}
