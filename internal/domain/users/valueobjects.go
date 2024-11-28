package users

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	IP    net.IP `json:"ip"`
}

type AccessToken string

func NewAccessToken(email string, ip net.IP, tokenLifetime time.Duration, secretKey string) (*AccessToken, error) {
	const method = "NewAccessToken"

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		},
		Email: email,
		IP:    ip,
	})

	signedJWT, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}
	access := AccessToken(signedJWT)

	return &access, nil
}

func (t AccessToken) String() string {
	return string(t)
}

func (t AccessToken) ParseIP(secretKey string) (net.IP, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(
		t.String(),
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
		jwt.WithValidMethods([]string{"HS512"}),
	)
	if err != nil || !token.Valid {
		return nil, ErrInvalidAccessToken
	}

	return claims.IP, nil
}

func (t AccessToken) ParseEmail(secretKey string) (string, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(
		t.String(),
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
		jwt.WithValidMethods([]string{"HS512"}),
	)
	if err != nil || !token.Valid {
		return "", ErrInvalidAccessToken
	}

	return claims.Email, nil
}

type RefreshToken []byte

func NewRefreshToken() (*RefreshToken, error) {
	const method = "NewRefreshToken"

	randomBytes := make([]byte, 72)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}
	refresh := RefreshToken(randomBytes)

	return &refresh, nil
}

func ParseRefreshToken(refreshBase64 string) (*RefreshToken, error) {
	token, err := base64.StdEncoding.DecodeString(refreshBase64)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	refresh := RefreshToken(token)

	return &refresh, nil
}

func (t RefreshToken) GenerateHash() ([]byte, error) {
	return bcrypt.GenerateFromPassword(t, bcrypt.DefaultCost)
}

func (t RefreshToken) Valid(hash string) bool {
	return errors.Is(nil, bcrypt.CompareHashAndPassword([]byte(hash), t))
}

func (t RefreshToken) String() string {
	return base64.StdEncoding.EncodeToString(t)
}
