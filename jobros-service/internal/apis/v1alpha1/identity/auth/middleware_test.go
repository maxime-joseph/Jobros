package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*gin.Engine, *JWTManager) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Set JWT_SECRET_KEY for testing
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	jwtManager, err := NewJWTManager()
	if err != nil {
		t.Fatalf("Failed to create JWTManager: %v", err)
	}

	router.Use(AuthMiddleware(jwtManager))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return router, jwtManager
}

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	router, _ := setupTest(t)
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}

func TestAuthMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	router, _ := setupTest(t)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid authorization header format")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	router, _ := setupTest(t)
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	router, jwtManager := setupTest(t)

	// Generate an expired token
	accessClaims := &JWTClaims{
		UserID: "testuser",
		Role:   "testrole",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtManager.secret)
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", expiredToken))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestAuthMiddleware_InvalidSignature(t *testing.T) {
	router, _ := setupTest(t)

	// Generate a token with a different secret
	accessClaims := &JWTClaims{
		UserID: "testuser",
		Role:   "testrole",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	invalidSecret := []byte("invalid-secret")
	invalidToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(invalidSecret)
	if err != nil {
		t.Fatalf("Failed to generate invalid signature token: %v", err)
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidToken))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestAuthMiddleware_ValidToken_WithRealJWTManager(t *testing.T) {
	router, jwtManager := setupTest(t)

	// Generate a valid token
	accessToken, err := jwtManager.GenerateAccessToken("testuser", "testrole")
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
