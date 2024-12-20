package auth

import (
	"github.com/maxime-joseph/Jobros/jobros-service/internal/apis"
	"github.com/maxime-joseph/Jobros/jobros-service/runtime"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestNewJWTManager(t *testing.T) {
	// Test with missing secret
	os.Unsetenv("JWT_SECRET_KEY")
	_, err := NewJWTManager()
	if err == nil {
		t.Error("Expected error when JWT_SECRET_KEY is not set")
	}

	// Test with valid secret
	os.Setenv("JWT_SECRET_KEY", "test-secret")
	manager, err := NewJWTManager()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if manager == nil {
		t.Error("Expected manager to be created")
	}
}

func TestJWTManager_GetGroupVersionKind(t *testing.T) {
	manager, _ := NewJWTManager()
	expectedGVK := runtime.GroupVersionKind{
		Group:   apis.APIGroup,
		Version: apis.APIVersion,
		Kind:    "JWTManager",
	}

	actualGVK := manager.GetGroupVersionKind()

	if !reflect.DeepEqual(actualGVK, expectedGVK) {
		t.Errorf("GetGroupVersionKind() returned %+v, expected %+v", actualGVK, expectedGVK)
	}
}

func TestJWTManager_DeepCopy(t *testing.T) {
	manager, _ := NewJWTManager()

	// Modify the secret of the original manager to ensure deep copy works
	originalSecret := []byte("test-secret-key")
	manager.secret = originalSecret

	copiedManager := manager.DeepCopy().(*JWTManager)

	// Check if the copied manager is a different instance
	if copiedManager == manager {
		t.Errorf("DeepCopy() returned the same instance")
	}

	// Check if the secret is copied correctly
	if !reflect.DeepEqual(copiedManager.secret, originalSecret) {
		t.Errorf("DeepCopy() did not copy the secret correctly, got: %v, want: %v", copiedManager.secret, originalSecret)
	}

	// Modify the copied manager's secret and check if the original is still the same
	copiedManager.secret = []byte("modified-secret")
	if reflect.DeepEqual(copiedManager.secret, manager.secret) {
		t.Errorf("DeepCopy() did not create a deep copy, modification of copy affected original")
	}
}

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret")
	manager, _ := NewJWTManager()

	token, err := manager.GenerateAccessToken("user123", "admin")
	if err != nil {
		t.Errorf("Failed to generate access token: %v", err)
	}
	if token == "" {
		t.Error("Generated token is empty")
	}

	// Validate the generated token
	claims, err := manager.GetTokenClaims(token)
	if err != nil {
		t.Errorf("Failed to parse token claims: %v", err)
	}
	if claims.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got %s", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("Expected Role 'admin', got %s", claims.Role)
	}
}

func TestJWTManager_GenerateTokenPair(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret")
	manager, _ := NewJWTManager()

	accessToken, refreshToken, err := manager.GenerateTokenPair("user123", "admin")
	if err != nil {
		t.Errorf("Failed to generate token pair: %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Error("Generated tokens should not be empty")
	}
}

func TestJWTManager_ValidateToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret")
	manager, _ := NewJWTManager()

	tests := []struct {
		name      string
		setup     func() string
		wantValid bool
	}{
		{
			name: "Valid token",
			setup: func() string {
				token, _ := manager.GenerateAccessToken("user123", "admin")
				return token
			},
			wantValid: true,
		},
		{
			name: "Expired token",
			setup: func() string {
				claims := &JWTClaims{
					UserID: "user123",
					Role:   "admin",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
				}
				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
				return token
			},
			wantValid: false,
		},
		{
			name: "Future issued token",
			setup: func() string {
				claims := &JWTClaims{
					UserID: "user123",
					Role:   "admin",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					},
				}
				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
				return token
			},
			wantValid: false,
		},
		{
			name: "Invalid signature",
			setup: func() string {
				claims := &JWTClaims{
					UserID: "user123",
					Role:   "admin",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					},
				}
				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("wrong-secret"))
				return token
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setup()
			if got := manager.ValidateToken(token); got != tt.wantValid {
				t.Errorf("ValidateToken() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}

func TestJWTManager_GetTokenClaims(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret")
	manager, _ := NewJWTManager()

	// Test valid token claims
	token, _ := manager.GenerateAccessToken("user123", "admin")
	claims, err := manager.GetTokenClaims(token)
	if err != nil {
		t.Errorf("Failed to get token claims: %v", err)
	}
	if claims.UserID != "user123" || claims.Role != "admin" {
		t.Error("Invalid claims data")
	}

	// Test invalid token
	_, err = manager.GetTokenClaims("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}
