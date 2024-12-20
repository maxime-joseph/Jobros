package auth

import (
	"fmt"
	"github.com/maxime-joseph/Jobros/jobros-service/internal/apis"
	"github.com/maxime-joseph/Jobros/jobros-service/runtime"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	runtime.Object // Embed the Object interface
	secret         []byte
}

func NewJWTManager() (*JWTManager, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize JWT manager: %w", err)
	}
	return &JWTManager{secret: []byte(secret)}, nil
}

func getJWTSecret() (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY environment variable is not set")
	}
	return secret, nil
}
func (m *JWTManager) GetGroupVersionKind() runtime.GroupVersionKind {
	return runtime.GroupVersionKind{
		Group:   apis.APIGroup,
		Version: apis.APIVersion,
		Kind:    "JWTManager",
	}
}

func (m *JWTManager) DeepCopy() runtime.Object {
	// Since JWTManager only contains a byte slice, a simple copy is sufficient
	secretCopy := make([]byte, len(m.secret))
	copy(secretCopy, m.secret)

	return &JWTManager{
		secret: secretCopy,
	}
}

func (m *JWTManager) ValidateToken(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return m.secret, nil
	})

	if err != nil {
		return false
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return false
	}

	if claims.UserID == "" || claims.Role == "" {
		return false
	}

	// reject tokens that have expired
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return false
	}

	// reject tokens that claim to be issued in the future:
	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(time.Now()) {
		return false
	}

	// ensure the token isn't being used before its valid time:
	if claims.NotBefore != nil && time.Now().Before(claims.NotBefore.Time) {
		return false
	}

	return true
}

func (m *JWTManager) GetTokenClaims(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func (m *JWTManager) GenerateRefreshToken(userID string, role string) (string, error) {
	refreshClaims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(m.secret)
}

func (m *JWTManager) GenerateTokenPair(userID string, role string) (accessToken string, refreshToken string, err error) {
	accessToken, err = m.GenerateAccessToken(userID, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = m.GenerateRefreshToken(userID, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *JWTManager) GenerateAccessToken(userID string, role string) (string, error) {
	accessClaims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(m.secret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
