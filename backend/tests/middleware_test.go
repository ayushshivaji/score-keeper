package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ayush-sr/score-keeper/backend/internal/config"
	"github.com/ayush-sr/score-keeper/backend/internal/middleware"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthRouter() (*gin.Engine, *service.AuthService) {
	cfg := &config.Config{JWTSecret: "test-secret-for-middleware-tests!!"}
	authSvc := service.NewAuthService(nil, cfg)

	r := gin.New()
	r.Use(middleware.AuthRequired(authSvc))
	r.GET("/protected", func(c *gin.Context) {
		userID := c.MustGet("user_id").(uuid.UUID)
		c.JSON(200, gin.H{"user_id": userID.String()})
	})
	return r, authSvc
}

func TestAuthMiddleware_ValidTokenInCookie(t *testing.T) {
	r, authSvc := setupAuthRouter()
	token, _ := authSvc.GenerateAccessToken(uuid.New())

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidTokenInAuthHeader(t *testing.T) {
	r, authSvc := setupAuthRouter()
	token, _ := authSvc.GenerateAccessToken(uuid.New())

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	r, _ := setupAuthRouter()
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r, _ := setupAuthRouter()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer garbage")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	r, _ := setupAuthRouter()
	wrongSvc := service.NewAuthService(nil, &config.Config{JWTSecret: "different-secret-key-entirely!!!!!"})
	token, _ := wrongSvc.GenerateAccessToken(uuid.New())

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_CookieTakesPrecedence(t *testing.T) {
	r, authSvc := setupAuthRouter()
	token, _ := authSvc.GenerateAccessToken(uuid.New())

	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	req.Header.Set("Authorization", "Bearer bad-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200 (cookie precedence), got %d", w.Code)
	}
}

func TestAuthMiddleware_BearerPrefixRequired(t *testing.T) {
	r, authSvc := setupAuthRouter()
	token, _ := authSvc.GenerateAccessToken(uuid.New())

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", token) // no "Bearer "
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_SetsUserIDInContext(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test-secret-for-context-check!!!!!"}
	authSvc := service.NewAuthService(nil, cfg)
	userID := uuid.New()
	token, _ := authSvc.GenerateAccessToken(userID)

	r := gin.New()
	r.Use(middleware.AuthRequired(authSvc))
	var capturedID uuid.UUID
	r.GET("/check", func(c *gin.Context) {
		capturedID = c.MustGet("user_id").(uuid.UUID)
		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/check", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if capturedID != userID {
		t.Errorf("expected %s, got %s", userID, capturedID)
	}
}

// ---------------------------------------------------------------------------
// CORS
// ---------------------------------------------------------------------------

func TestCORS_SetsHeaders(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS("http://localhost:3000"))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Error("missing CORS origin")
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("missing credentials header")
	}
}

func TestCORS_OptionsReturns204(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS("http://localhost:3000"))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestCORS_CustomOrigin(t *testing.T) {
	r := gin.New()
	r.Use(middleware.CORS("https://myapp.com"))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Header().Get("Access-Control-Allow-Origin") != "https://myapp.com" {
		t.Error("expected custom origin")
	}
}
