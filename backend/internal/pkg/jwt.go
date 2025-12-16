package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    string    `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret           string
	expiresIn        time.Duration
	refreshExpiresIn time.Duration
}

func NewJWTService(secret string, expiresIn, refreshExpiresIn int) *JWTService {
	return &JWTService{
		secret:           secret,
		expiresIn:        time.Duration(expiresIn) * time.Second,
		refreshExpiresIn: time.Duration(refreshExpiresIn) * time.Second,
	}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

func (s *JWTService) GenerateTokenPair(userID, tenantID, email, role string) (*TokenPair, error) {
	accessToken, err := s.generateToken(userID, tenantID, email, role, AccessToken, s.expiresIn)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(userID, tenantID, email, role, RefreshToken, s.refreshExpiresIn)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.expiresIn.Seconds()),
	}, nil
}

func (s *JWTService) generateToken(userID, tenantID, email, role string, tokenType TokenType, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		TenantID:  tenantID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
