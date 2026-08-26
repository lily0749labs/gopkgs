package ginUtil

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewEngineConfigAndCallback(t *testing.T) {
	oldMode := gin.Mode()
	oldWriter := gin.DefaultWriter
	oldErrorWriter := gin.DefaultErrorWriter
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	t.Cleanup(func() {
		gin.SetMode(oldMode)
		gin.DefaultWriter = oldWriter
		gin.DefaultErrorWriter = oldErrorWriter
	})

	callbackCalled := false
	engine := NewEngine(EngineConfig{
		Mode:                   "prod",
		HandleMethodNotAllowed: true,
		MaxMultipartMemory:     60 << 20,
		CORS:                   PermissiveCORSConfig(),
		EnablePprof:            true,
	}, func(engine *gin.Engine) {
		callbackCalled = true
		engine.GET("/health", func(ctx *gin.Context) {
			ctx.Status(http.StatusNoContent)
		})
		engine.GET("/panic", func(*gin.Context) {
			panic("test panic")
		})
	})

	if !callbackCalled {
		t.Fatal("engine callback was not called")
	}
	if gin.Mode() != gin.ReleaseMode {
		t.Fatalf("gin mode = %q", gin.Mode())
	}
	if !engine.HandleMethodNotAllowed {
		t.Fatal("HandleMethodNotAllowed = false")
	}
	if engine.MaxMultipartMemory != 60<<20 {
		t.Fatalf("MaxMultipartMemory = %d", engine.MaxMultipartMemory)
	}

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("Origin", "https://example.com")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("GET /health status = %d", response.Code)
	}
	if response.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatal("CORS response header was not set")
	}

	postRequest := httptest.NewRequest(http.MethodPost, "/health", nil)
	postResponse := httptest.NewRecorder()
	engine.ServeHTTP(postResponse, postRequest)
	if postResponse.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST /health status = %d", postResponse.Code)
	}

	panicRequest := httptest.NewRequest(http.MethodGet, "/panic", nil)
	panicResponse := httptest.NewRecorder()
	engine.ServeHTTP(panicResponse, panicRequest)
	if panicResponse.Code != http.StatusInternalServerError {
		t.Fatalf("GET /panic status = %d", panicResponse.Code)
	}

	if !hasRoute(engine, http.MethodGet, "/debug/pprof/") {
		t.Fatal("pprof route was not registered")
	}
}

func TestNewEngineWithoutOptionalFeatures(t *testing.T) {
	oldMode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(oldMode) })

	engine := NewEngine(EngineConfig{Mode: "test"})
	if engine.HandleMethodNotAllowed {
		t.Fatal("HandleMethodNotAllowed = true")
	}
	if hasRoute(engine, http.MethodGet, "/debug/pprof/") {
		t.Fatal("pprof route was unexpectedly registered")
	}
}

func hasRoute(engine *gin.Engine, method, path string) bool {
	for _, route := range engine.Routes() {
		if route.Method == method && route.Path == path {
			return true
		}
	}
	return false
}
